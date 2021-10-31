package types

import (
	"errors"
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
func ParseDurationOrDie(s string) *time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return &d
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
	locOslo   = ParseLocationOrDie("Europe/Oslo")      // +02:00
	locHobart = ParseLocationOrDie("Australia/Hobart") // +11:00
)

func TestScheduleWeek_NextRun(t *testing.T) {
	type fields struct {
		Location  *time.Location
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
				location:  tt.fields.Location,
				Monday:    tt.fields.Monday,
				Tuesday:   tt.fields.Tuesday,
				Wednesday: tt.fields.Wednesday,
				Thursday:  tt.fields.Thursday,
				Friday:    tt.fields.Friday,
				Saturday:  tt.fields.Saturday,
				Sunday:    tt.fields.Sunday,
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

func TestSchedulePayload_Prepare(t *testing.T) {
	type fields struct {
		RequestID              string
		EndpointID             string
		MaxInterJobConcurrency bool
		Label                  string
		StartDate              time.Time
		EndDate                *time.Time
		Frequency              Frequency
		Multiplier             float64
		Offsets                []int
		Config                 *Config
		ScheduleWeek           ScheduleWeek
	}
	tests := []struct {
		name    string
		fields  fields
		want    SchedulePayload
		wantErr error
	}{
		{
			"should set location to Local of no location-data is present from start-date",
			fields{StartDate: *now},
			SchedulePayload{
				ScheduleWeek: ScheduleWeek{
					location: ParseLocationOrDie("Local"),
					Location: "Local",
				},
				StartDate: *now,
			},
			nil,
		},
		{
			"should fill location correctly from location-object",
			fields{StartDate: *now, ScheduleWeek: ScheduleWeek{location: ParseLocationOrDie("Europe/Oslo")}},
			SchedulePayload{
				ScheduleWeek: ScheduleWeek{
					location: ParseLocationOrDie("Europe/Oslo"),
					Location: "Europe/Oslo",
				},
				StartDate: *now,
			},
			nil,
		},
		{
			"should fill location correctly from Location-string",
			fields{StartDate: *now, ScheduleWeek: ScheduleWeek{Location: "Europe/Oslo"}},
			SchedulePayload{
				ScheduleWeek: ScheduleWeek{
					location: ParseLocationOrDie("Europe/Oslo"),
					Location: "Europe/Oslo",
				},
				StartDate: *now,
			},
			nil,
		},
		{
			"should err on invalid Location-string",
			fields{StartDate: *now, ScheduleWeek: ScheduleWeek{Location: "Haua"}},
			SchedulePayload{
				ScheduleWeek: ScheduleWeek{
					location: ParseLocationOrDie("Europe/Oslo"),
					Location: "Europe/Oslo",
				},
				StartDate: *now,
			},
			errors.New("Location was incorrect: unknown time zone Haua"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := SchedulePayload{
				RequestID:              tt.fields.RequestID,
				EndpointID:             tt.fields.EndpointID,
				MaxInterJobConcurrency: tt.fields.MaxInterJobConcurrency,
				Label:                  tt.fields.Label,
				StartDate:              tt.fields.StartDate,
				EndDate:                tt.fields.EndDate,
				Frequency:              tt.fields.Frequency,
				Multiplier:             tt.fields.Multiplier,
				Offsets:                tt.fields.Offsets,
				Config:                 tt.fields.Config,
				ScheduleWeek:           tt.fields.ScheduleWeek,
			}
			got, err := sp.Prepare()
			errStr := nilErrString(err)
			wantErrStr := nilErrString(tt.wantErr)
			if errStr != wantErrStr {
				t.Errorf("Mismatched error: \ngot:  %s\nwant: %s", errStr, wantErrStr)
			}
			if tt.wantErr != nil {
				return
			}
			gotStr := yamlString(got)
			wantStr := yamlString(tt.want)
			t.Log(got.location, tt.want.location)
			if gotStr != wantStr {
				t.Errorf("SchedulePayload.Prepare() = \n\tgot\n%v \n\twant\n%v", gotStr, wantStr)
			}
		})
	}
}

func yamlString(j interface{}) string {
	b, err := yaml.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(b)
}
