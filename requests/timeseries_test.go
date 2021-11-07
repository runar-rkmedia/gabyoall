package requests

import (
	"reflect"
	"testing"
	"time"
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
		name string
		args args
		want []time.Time
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTimeSeries(tt.args.start)
			for _, v := range tt.args.times {
				t.Log("adding", v)
				ts.Push(v, 1.0)
			}
			ts.Finish()
			iter := ts.Iter()
			var gotTimes []time.Time
			for {
				if !iter.Next() {
					break
				}
				tk, _ := iter.Values()
				gotTimes = append(gotTimes, time.Unix(int64(tk), 0))

			}
			if !reflect.DeepEqual(gotTimes, tt.want) {

				t.Errorf("ScheduleWeek.NextRun() = %v want %v", gotTimes, tt.want)
			}
		})
	}
}
func ParseTimeOrDie(s string) *time.Time {
	n, err := time.Parse("2006-01-02T15:04:05Z0700", s)
	if err != nil {
		panic(err)
	}
	return &n
}
