package requests

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/runar-rkmedia/gabyoall/internal"
)

var (
	now = ParseTimeOrDie("2021-10-29T15:58:10+0200")
)

func TestTimeSeries(t *testing.T) {

	type args struct {
		start  time.Time
		times  []time.Time
		values []float64
	}
	tests := []struct {
		name    string
		args    args
		want    []time.Time
		wantExp TimeSeriesExpanded
	}{
		{
			"should return correctly for friday.",
			args{
				start: *now,
				times: []time.Time{
					*ParseTimeOrDie("2021-10-29T16:03:10+0200"),
					*ParseTimeOrDie("2021-10-29T16:04:10+0200"),
				},
			},
			[]time.Time{
				*ParseTimeOrDie("2021-10-29T16:03:10+0200"),
				*ParseTimeOrDie("2021-10-29T16:04:10+0200"),
			},
			TimeSeriesExpanded{
				StartTime: *ParseTimeOrDie("2021-10-29T16:03:10+0200"),
				Series: []Serie{
					[2]uint64{0, 1000},
					[2]uint64{60000, 2000},
				},
			},
		},
		{
			"should still work if times are not ordered.",
			args{
				start: *now,
				times: []time.Time{
					*ParseTimeOrDie("2021-10-29T16:03:10+0200"),
					*ParseTimeOrDie("2021-10-29T16:04:10+0200"),
					*ParseTimeOrDie("2021-10-29T16:01:10+0200"),
					*ParseTimeOrDie("2021-10-29T14:54:10+0200"),
				},
			},
			[]time.Time{
				*ParseTimeOrDie("2021-10-29T14:54:10+0200"),
				*ParseTimeOrDie("2021-10-29T16:01:10+0200"),
				*ParseTimeOrDie("2021-10-29T16:03:10+0200"),
				*ParseTimeOrDie("2021-10-29T16:04:10+0200"),
			},
			TimeSeriesExpanded{
				StartTime: *ParseTimeOrDie("2021-10-29T16:03:10+0200"),
				Series: []Serie{
					[2]uint64{0, 1000},
					[2]uint64{60000, 2000},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTimeSeries(tt.args.start)
			i := 0.0
			for _, v := range tt.args.times {
				i++
				// t.Log("adding", v)
				ts.Push(v, i*float64(time.Second))
			}
			ts.Finish()
			iter := ts.Iter()
			var gotTimes []time.Time
			for {
				if !iter.Next() {
					break
				}
				tk, _ := iter.Values()
				gotTimes = append(gotTimes, time.UnixMilli(int64(tk)).UTC())

			}
			if err := internal.Compare("ScheduleWeek.NextRun()", gotTimes, tt.want); err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(gotTimes, tt.want) {
				t.Errorf("ScheduleWeek.NextRun() = \n%v\nwant\n%v", gotTimes, tt.want)
			}
			exp := *ts.Expand()
			if diff := deep.Equal(exp, tt.wantExp); diff != nil {
				t.Errorf("ts.Expand() = \n%v\ndiff:\n%s\nwant:\n%v", exp, diff, tt.wantExp)
			}
			if !reflect.DeepEqual(exp, tt.wantExp) {
				t.Errorf("ts.Expand() = \n%v\nwant\n%v", exp, tt.wantExp)
			}

		})
	}
}
func ParseTimeOrDie(s string) *time.Time {
	n, err := time.Parse("2006-01-02T15:04:05Z0700", s)
	if err != nil {
		panic(err)
	}
	utc := n.UTC()
	return &utc
}
