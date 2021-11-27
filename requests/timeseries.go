package requests

/*

This turned out a bit more complex than originally anticipated.

- We try to compact the timeseries upon storage as much as possible. For this we use the gorilla algorithm from go-tsz.
- However, this requries all timeseries to be pushed in a timely ordered fashion.
- We want to push timeseries from each go-routine
- We want to preview the timeseries as it is populated on an interval.

*/

import (
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/tsenart/go-tsz"
)

type TimeSeries struct {
	ordered        bool
	finished       bool
	orderingValues *orderingSlice
	tsz.Series
}

type orderingSlice []Serie

func (o orderingSlice) Len() int           { return len(o) }
func (o orderingSlice) Less(i, j int) bool { return o[i][0] > o[j][0] }
func (o orderingSlice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

type TimeSeriesOptions struct {
	// Set to true if the caller can guarantee that the pushed values are coming in at an timely ordered fashion.
	// This will reduce the memory-requirement, but if items are out of order, the times produced will be garbled.
	Ordered bool
}

func NewTimeSeries(startTime time.Time, options ...TimeSeriesOptions) *TimeSeries {
	var opts TimeSeriesOptions
	if len(options) > 0 {
		opts = options[0]
	} else {
		opts = TimeSeriesOptions{false}
	}
	st := uint64(startTime.UnixMilli())
	t := TimeSeries{
		ordered: opts.Ordered,
		Series:  *tsz.New(st),
	}
	if !t.ordered {
		t.orderingValues = &orderingSlice{}
	}
	return &t
}

func (ts *TimeSeries) Finish() {
	// This may require locking
	if ts.finished {
		return
	}
	if !ts.ordered {
		ts.ordered = true
		sort.Sort(*ts.orderingValues)
		for i := 0; i < len(*ts.orderingValues); i++ {
			// fmt.Println(
			// 	(*ts.orderingValues)[i].time, (*ts.orderingValues)[i].v,
			// )
			ts.Push(TimeSeriePointToTime((*ts.orderingValues)[i][0]), float64((*ts.orderingValues)[i][1]))
		}
	}
	ts.Series.Finish()
	ts.finished = true
}

// Each push must be in order!
func (ts *TimeSeries) Push(t time.Time, v float64) {
	st := uint64(t.UnixMilli())
	if !ts.ordered {
		*ts.orderingValues = append(*ts.orderingValues, Serie{st, uint64(v)})
		return
	}
	ts.Series.Push(st, v)
	return
}
func (ts *TimeSeries) Expand() *TimeSeriesExpanded {
	if ts == nil {
		return nil
	}
	exp := TimeSeriesExpanded{}
	if !ts.finished && !ts.ordered {
		if len(*ts.orderingValues) == 0 {
			return &exp
		}
		sort.Sort(ts.orderingValues)
		t0 := uint64((*ts.orderingValues)[0][0])
		exp.StartTime = TimeSeriePointToTime((*ts.orderingValues)[0][0])
		for i := 0; i < len(*ts.orderingValues); i++ {
			t := (*ts.orderingValues)[i][0]
			exp.Series = append(exp.Series, Serie{t - t0, uint64((*ts.orderingValues)[i][1])})
		}

		return &exp
	}
	iter := ts.Iter()
	if ok := iter.Next(); !ok {
		return nil
	}

	t0, k := iter.Values()
	exp.StartTime = TimeSeriePointToTime(t0)
	exp.Series = append(exp.Series, Serie{uint64(0), uint64(k / 1_000_000)})
	for {
		if !iter.Next() {
			break
		}
		t, k := iter.Values()
		// tt := TimeSeriePointToTime(t)
		// milli := tt.UnixMilli() - exp.StartTime.UnixMilli()
		exp.Series = append(exp.Series, Serie{t - t0, uint64(k / 1_000_000)})
	}
	return &exp

}

func TimeSeriePointToTime(v uint64) time.Time {
	return time.UnixMilli(int64(v)).UTC()
}

type TimeSeriesMap struct {
	Map       map[string]*TimeSeries
	StartTime time.Time
	lock      sync.RWMutex
}

func (tsm *TimeSeriesMap) Push(label string, t time.Time, value float64) {

	// Not really sure if this is beneficial, but the idea is this:
	// Normally, a new Timeseries with label x is only created only, but there are lots of .Push'es.
	// Therefire, we attempt to only use a read-lock and check if we need to add a TimeSeries, which will require locking.
	// When we issue the write-lock, we must first unlock the read-lock.
	tsm.lock.RLock()
	defer tsm.lock.RUnlock()
	if _, ok := tsm.Map[label]; !ok {
		tsm.lock.RUnlock()
		tsm.lock.Lock()
		tsm.Map[label] = NewTimeSeries(tsm.StartTime)
		tsm.lock.Unlock()
		tsm.lock.RLock()
	}
	tsm.Map[label].Push(t, value)
}

func NewTimeSeriesWithLabel(startTime time.Time) TimeSeriesMap {
	return TimeSeriesMap{
		StartTime: startTime,
		Map:       map[string]*TimeSeries{},
	}
}
func (s *TimeSeriesMap) MarshalJSON() ([]byte, error) {
	maps := s.Expand()

	if maps == nil || len(maps) == 0 {
		return nil, nil
	}
	for k, v := range maps {
		n := v.DropResolution(800)
		maps[k] = &n
	}
	return json.Marshal(maps)
}
func (tsm *TimeSeriesMap) Expand() map[string]*TimeSeriesExpanded {
	tsm.lock.RLock()
	defer tsm.lock.RUnlock()

	if len(tsm.Map) == 0 {
		return nil
	}
	maps := map[string]*TimeSeriesExpanded{}
	for k, v := range tsm.Map {
		maps[k] = v.Expand()
	}

	return maps
}

type TimeSeriesExpanded struct {
	StartTime time.Time
	Series    []Serie
}

// This is probably terrible for performance, but still better than serving 500000 items to the client...
// With a real timeseries-database this kind of operation would be a lot less expensive
func (tse TimeSeriesExpanded) DropResolution(maxResolution int) TimeSeriesExpanded {
	length := len(tse.Series)
	if length <= maxResolution {
		return tse
	}
	bucketLength := length / maxResolution

	tseNew := TimeSeriesExpanded{StartTime: tse.StartTime, Series: make([]Serie, maxResolution)}
	j := 0
	inBucket := 0
	for i := 0; i < length; i++ {
		inBucket++
		tseNew.Series[j][0] += tse.Series[i][0]
		if tseNew.Series[j][1] < tse.Series[i][1] {
			tseNew.Series[j][1] = tse.Series[i][1]
		}
		if i == length-1 || (i > 0 && i%bucketLength == 0) {
			tseNew.Series[j][0] /= uint64(inBucket)
			if j < maxResolution-1 {
				j++
				inBucket = 0
			}
		}
	}

	return tseNew
}

type Serie [2]uint64
