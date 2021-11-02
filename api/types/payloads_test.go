package types

import (
	"reflect"
	"testing"
	"time"

	"github.com/ghodss/yaml"
)

func ParseTimeOrDie(s string) *time.Time {
	n, err := time.Parse("2006-01-02T15:04:05Z0700", s)
	if err != nil {
		panic(err)
	}
	return &n
}
func ParseDurationOrDie(s string) *Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return &Duration{d}
}
func ParseLocationOrDie(s string) *time.Location {
	d, err := time.LoadLocation(s)
	if err != nil {
		panic(err)
	}
	return d
}

var (
	now       = ParseTimeOrDie("2021-10-29T15:58:10+0200")
	locOslo   = "Europe/Oslo"      // +02:00
	locHobart = "Australia/Hobart" // +11:00
)

func TestScheduleWeek_NextRun(t *testing.T) {
	type fields struct {
		Location  string
		Monday    *Duration
		Tuesday   *Duration
		Wednesday *Duration
		Thursday  *Duration
		Friday    *Duration
		Saturday  *Duration
		Sunday    *Duration
	}

	type args struct {
		start time.Time
		end   *time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *time.Time
	}{
		{
			"should return correctly for friday.",
			fields{
				Location: locOslo,
				Friday:   ParseDurationOrDie("19h30m"),
			},
			args{*now, nil},
			ParseTimeOrDie("2021-10-29T19:30:00+0200"),
		},
		{
			"should return nil if time is passed",
			fields{
				Location: locOslo,
				Friday:   ParseDurationOrDie("13h30m"),
			},
			args{*now, nil},
			nil,
		},
		{
			"should return nil if endTime is passed",
			fields{
				Location: locOslo,
				Friday:   ParseDurationOrDie("19h30m"),
			},
			args{
				start: *now,
				end:   ParseTimeOrDie("2021-10-29T18:58:10+0200"),
			},
			nil,
		},
		{
			"should return work correctly with other timezone",
			fields{
				Location: locHobart,
				Thursday: ParseDurationOrDie("21h30m"),
				Friday:   ParseDurationOrDie("20h30m"),
				Saturday: ParseDurationOrDie("19h30m"), // This should be the one that is "picked"
			},
			args{*now, nil},
			ParseTimeOrDie("2021-10-30T19:30:00+1100"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := ScheduleWeek{
				LocationStr: tt.fields.Location,
				Monday:      tt.fields.Monday,
				Tuesday:     tt.fields.Tuesday,
				Wednesday:   tt.fields.Wednesday,
				Thursday:    tt.fields.Thursday,
				Friday:      tt.fields.Friday,
				Saturday:    tt.fields.Saturday,
				Sunday:      tt.fields.Sunday,
			}
			if got := sw.NextRun(tt.args.start, tt.args.end); !reflect.DeepEqual(nilTimeString(got), nilTimeString(tt.want)) {
				t.Errorf("ScheduleWeek.NextRun() = %v (%v), want %v (%v)", got, nilTimeString(got), tt.want, nilTimeString(tt.want))
			}
		})
	}
}

func nilTimeString(s *time.Time) string {
	if s == nil {
		return ""
	}
	return s.UTC().String()
}
func nilErrString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func yamlString(j interface{}) string {
	b, err := yaml.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(b)
}
