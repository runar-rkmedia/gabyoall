package types

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/runar-rkmedia/gabyoall/requests"
)

type Entity struct {
	// Time of which the entity was created in the database
	// Required: true
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// Time of which the entity was updated, if any
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	// Unique identifier of the entity
	// Required: true
	ID string `json:"id,omitempty"`
	// If set, the item is considered deleted. The item will normally not get deleted from the database,
	// but it may if cleanup is required.
	Deleted *time.Time `json:"deleted,omitempty"`
}

type EndpointPayload struct {
	// required: true
	// example: https://example.com
	Url     string              `json:"url,omitempty" validate:"required,uri"`
	Headers map[string][]string `json:"headers,omitempty" validate:"dive,max=1000"`
	Config  *Config             `json:"config,omitempty"`
}
type RequestPayload struct {
	Body          string                 `json:"body,omitempty"`
	Query         string                 `json:"query,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty" validate:"dive,max=1000"`
	Headers       map[string]string      `json:"headers,omitempty" validate:"dive,max=1000"`
	OperationName string                 `json:"operationName,required"`
	Method        string                 `json:"method"`
	Config        Config                 `json:"config,omitempty"`
}

type EndpointEntity struct {
	Entity
	Endpoint
}
type Endpoint struct {
	requests.Endpoint
	Config *Config `json:"config,omitempty"`
}
type RequestEntity struct {
	requests.Request
	Entity
	Config *Config `json:"config,omitempty"`
}

type ScheduleEntity struct {
	Entity
	Schedule
	Config *Config `json:"config,omitempty"`
}
type StatEntity struct {
	Entity
	requests.CompactRequestStatistics
}
type StatPayload = requests.CompactRequestStatistics

type Schedule struct {
	// These are calculated in create/update. These are used for faster lookups.
	// Should be ordered Ascending, e.g. the first element
	Dates []time.Time `json:"dates"`
	// FIXME: Should not use schedulePaylaod directly here.
	SchedulePayload
	// From these, the dates above can be calculated
	LastRun   *time.Time `json:"lastRun,omitempty"`
	LastError string     `json:"lastError,omitempty"`

	// TODO: implement runs, which hold historical information about a run.
	// These should hold a reference to the object of which is was created with,
	// But since that object might change in the future, we want to store those
	// parameters normalized in a Run too.
	// RunIDS []string
}

// Gets the next calculated run. Does not include the ForcedStartDate
func (s Schedule) NextRun() *time.Time {
	lastRun := s.LastRun
	if lastRun == nil {
		lastRun = &s.StartDate
	}
	if n := s.ScheduleWeek.NextRun(*lastRun, s.EndDate); n != nil {
		return n
	}
	backOffTime := 1 * time.Minute
	// TODO: this should calculate based on other parameters
	t := s.StartDate
	now := time.Now()
	if s.LastRun != nil && s.LastRun.After(t) {
		if s.LastError != "" {
			if now.Sub(*s.LastRun) < backOffTime {
				return nil
			}
		} else {
			// TODO: make sure this logic works before attempting retries
		}
		return nil
	}
	return &t
}

// Indicates whether this schedule should be run based on run-parameters
func (s Schedule) ShouldRun() bool {
	if s.ForcedStartDate != nil && time.Now().After(*s.ForcedStartDate) {
		if s.LastRun == nil || s.LastRun.Before(*s.ForcedStartDate) {

			return true
		}
	}
	nextRun := s.NextRun()

	if nextRun == nil {
		return false
	}

	sh := time.Now().After(*nextRun)
	// fmt.Println("should run", sh, s.Label, nextRun.String())
	// return false
	return sh
}

type SchedulePayload struct {
	RequestID  string `json:"requestID" validate:"required"`
	EndpointID string `json:"endpointID" validate:"required"`
	// If set to a positive value, the scheduler will not schedule more than this total concurrency
	// when starting this job, and when it is running.
	//
	// Some jobs might have been configured to run very slowly, with low concurrency,
	// high wait-times and can therefore run alongside other such jobs.
	MaxInterJobConcurrency bool       `json:"maxInterJobConcurrency"`
	Label                  string     `json:"label" validate:"required"`
	StartDate              time.Time  `json:"start_date" validate:"required"`
	ForcedStartDate        *time.Time `json:"forced_start_date"`
	EndDate                *time.Time `json:"end_date"`
	Frequency              Frequency  `json:"frequency,omitempty"`
	Multiplier             float64    `json:"multiplier,omitempty"`
	// TODO: either document what this value is or drop it. I dont remember why I added this.
	// I am sure there was a reason, though...
	Offsets []int   `json:"offsets,omitempty"`
	Config  *Config `json:"config,omitempty"`
	ScheduleWeek
}

type ScheduleWeek struct {
	LocationStr string    `json:"location"`
	Monday      *Duration `json:"monday,omitempty"`
	Tuesday     *Duration `json:"tuesday,omitempty"`
	Wednesday   *Duration `json:"wednesday,omitempty"`
	Thursday    *Duration `json:"thursday,omitempty"`
	Friday      *Duration `json:"friday,omitempty"`
	Saturday    *Duration `json:"saturday,omitempty"`
	Sunday      *Duration `json:"sunday,omitempty"`
}

func (sw ScheduleWeek) NextRun(start time.Time, end *time.Time) *time.Time {
	if sw.Monday == nil && sw.Tuesday == nil && sw.Wednesday == nil && sw.Thursday == nil && sw.Friday == nil && sw.Saturday == nil && sw.Sunday == nil {
		return nil
	}
	loc, err := time.LoadLocation(sw.LocationStr)
	if err != nil {
		l := logger.GetLogger("schedule-week")
		l.Error().Str("location-string", sw.LocationStr).Err(err).Msg("failed to load location-string")
		return nil
	}
	t := start.In(loc)
	dow := t.Weekday()
	var d *Duration
	switch dow {
	case time.Monday:
		d = sw.Monday
	case time.Tuesday:
		d = sw.Tuesday
	case time.Wednesday:
		d = sw.Wednesday
	case time.Thursday:
		d = sw.Thursday
	case time.Friday:
		d = sw.Friday
	case time.Saturday:
		d = sw.Saturday
	case time.Sunday:
		d = sw.Sunday
	}
	if d == nil {
		return nil
	}
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	n := midnight.Add(d.Duration)
	if n.Before(t) {
		return nil
	}
	if end != nil && n.After(*end) {
		return nil
	}
	return &n
}

type Duration struct{ time.Duration }

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

func (s SchedulePayload) nextRun(startTime time.Time, x time.Duration) (ts []time.Time) {
	switch s.Frequency {
	// TODO: Up until and uncluding HOUR, we can use modulo
	case FrequencySecond:
		ts = append(ts, startTime.Add(time.Second*x))
	case FrequencyMinute:
		ts = append(ts, startTime.Add(time.Minute*x))
	case FrequencyHour:
		ts = append(ts, startTime.Add(time.Hour*x))
	case FrequencyDay:
		// TODO: For these, I think we need to use an actual calendar, and loop through them from the start-date until we get to somewhere next week.
		ts = append(ts, startTime.Add(time.Hour*x*24))
	case FrequencyWeek:
		// TODO: this should be mapped by a real week
		ts = append(ts, startTime.Add(time.Hour*x*24*7))
	case FrequencyMonth:
		// TODO: this should be mapped by a real month
		ts = append(ts, startTime.Add(time.Hour*x*24*30))
	}
	return
}
func (s SchedulePayload) CalculateNextRuns(maxCount int) (ts []time.Time) {

	if s.Frequency == FrequencyNull {
		return ts
	}
	x := time.Duration(s.Multiplier)
	if x == 0 {
		x = 1
	}
	// TODO: this value should be calculated after last run
	startTime := s.StartDate
	for i := 0; i < maxCount; i++ {

		switch s.Frequency {
		case FrequencySecond:
			ts = append(ts, startTime.Add(time.Second*x))
		case FrequencyMinute:
			ts = append(ts, startTime.Add(time.Minute*x))
		case FrequencyHour:
			ts = append(ts, startTime.Add(time.Hour*x))
		case FrequencyDay:
			ts = append(ts, startTime.Add(time.Hour*x*24))
		case FrequencyWeek:
			// TODO: this should be mapped by a real week
			ts = append(ts, startTime.Add(time.Hour*x*24*7))
		case FrequencyMonth:
			// TODO: this should be mapped by a real month
			ts = append(ts, startTime.Add(time.Hour*x*24*30))
		}
		//if s.Multiplier != 0 {
		//	// Don't really care much if we round a microsecond off.
		//	d = time.Duration(s.Multiplier * float64(d))
		//}
	}
	return ts
}

type Frequency int8

const (
	FrequencyNull   Frequency = iota
	FrequencySecond Frequency = iota
	FrequencyMinute Frequency = iota
	FrequencyHour   Frequency = iota
	FrequencyDay    Frequency = iota
	FrequencyWeek   Frequency = iota
	FrequencyMonth  Frequency = iota
)

type ServerInfo struct {
	// When the server was started
	ServerStartedAt time.Time
	// Short githash for current commit
	GitHash string
	// Version-number for commit
	Version string
	// Date of build
	BuildDate time.Time

	// Size of database.
	DatabaseSize    int64
	DatabaseSizeStr string
}
