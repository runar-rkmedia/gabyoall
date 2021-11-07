package requests

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/tsenart/go-tsz"
)

type TimeSeries struct {
	tsz.Series
}

func NewTimeSeries(startTime time.Time) *TimeSeries {
	st := uint64(startTime.UnixMilli())
	return &TimeSeries{
		*tsz.New(st),
	}
}

// Each push must be in order!
func (ts *TimeSeries) Push(t time.Time, v float64) {
	st := uint64(t.UnixMilli())
	ts.Series.Push(st, v)
	return
}
func (ts *TimeSeries) Expand() *TimeSeriesExpanded {
	if ts == nil {
		return nil
	}
	iter := ts.Iter()
	exp := TimeSeriesExpanded{}
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
	return time.UnixMilli(int64(v))
}

type TimeSeriesMap struct {
	Map       map[string]*TimeSeries
	StartTime time.Time
	lock      sync.RWMutex
}

func (tsm *TimeSeriesMap) Push(label string, t time.Time, value float64) {

	// Not really sure if this is beneficial, but the idea is this:
	// Normally, a new Timeseries with lable x is only created only, but there are lots of .Push'es.
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
