package types

import (
	"time"

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

func (s Schedule) NextRun() *time.Time {
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
	MaxInterJobConcurrency bool      `json:"maxInterJobConcurrency"`
	Label                  string    `json:"label" validate:"required"`
	StartDate              time.Time `json:"start_date" validate:"required"`
	Frequency              Frequency `json:"frequency,omitempty"`
	Multiplier             float64   `json:"multiplier,omitempty"`
	Offsets                []int     `json:"offsets,omitempty"`
	Config                 *Config   `json:"config,omitempty"`
	ScheduleWeek
}

func (sp SchedulePayload) Prepare() (SchedulePayload, error) {
	if sp.location == nil {
		if sp.Location != "" {
			l, err := time.LoadLocation(sp.Location)
			if err != nil {
				return sp, fmt.Errorf("Location was incorrect: %w", err)
			}
			sp.location = l
		} else {
			sp.location = sp.StartDate.Location()
		}
	}
	if sp.Location == "" {
		sp.Location = fmt.Sprintf("%s", sp.location)
	}
	return sp, nil
}

type ScheduleWeek struct {
	location  *time.Location
	Location  string    `json:"location"`
	Monday    *Duration `json:"monday,omitempty"`
	Tuesday   *Duration `json:"tuesday,omitempty"`
	Wednesday *Duration `json:"wednesday,omitempty"`
	Thursday  *Duration `json:"thursday,omitempty"`
	Friday    *Duration `json:"friday,omitempty"`
	Saturday  *Duration `json:"saturday,omitempty"`
	Sunday    *Duration `json:"sunday,omitempty"`
}

func (sw ScheduleWeek) NextRun(now time.Time) *time.Time {
	t := now.In(sw.location)
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
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, sw.location)
	n := midnight.Add(*d)
	if n.Before(t) {
		fmt.Printf("%v is before %v (%s) %v \n\n", n, t, dow, midnight)
		// fmt.Printf("midnight: %v  \n", midnight, n, t, dow, sw.Location)
		return nil
	}
	return &n
}

type Duration = time.Duration

}

type Frequency int8

const (
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
