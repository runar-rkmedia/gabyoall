package types

import (
	"time"

	"github.com/runar-rkmedia/gabyoall/requests"
)

type Storage interface {
	Endpoints() (es map[string]EndpointEntity, err error)
	Endpoint(id string) (EndpointEntity, error)
	CreateEndpoint(e EndpointPayload) (EndpointEntity, error)

	Requests() (es map[string]RequestEntity, err error)
	Request(id string) (RequestEntity, error)
	CreateRequest(e RequestPayload) (RequestEntity, error)

	Schedules() (es map[string]ScheduleEntity, err error)
	Schedule(id string) (ScheduleEntity, error)
	CreateSchedule(e SchedulePayload) (ScheduleEntity, error)
	UpdateSchedule(id string, p Schedule) (ScheduleEntity, error)
}

type Entity struct {
	// Time of which the entity was created in the database
	// Required: true
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// Time of which the entity was updated, if any
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	// Unique identifier of the entity
	// Required: true
	ID string `json:"id,omitempty"`
}

type EndpointPayload struct {
	// required: true
	// example: https://example.com
	Url     string              `json:"url,omitempty" validate:"required,uri"`
	Headers map[string][]string `json:"headers,omitempty" validate:"dive,max=1000"`
}
type RequestPayload struct {
	Body          string                 `json:"body,omitempty"`
	Query         string                 `json:"query,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty" validate:"dive,max=1000"`
	Headers       map[string]string      `json:"headers,omitempty" validate:"dive,max=1000"`
	OperationName string                 `json:"operationName,required"`
	Method        string                 `json:"method"`
}

type EndpointEntity struct {
	requests.Endpoint
	Entity
}
type RequestEntity struct {
	requests.Request
	Entity
}

type ScheduleEntity struct {
	Entity
	Schedule
}
type Schedule struct {
	// These are calculated in create/update. These are used for faster lookups.
	// Should be ordered Ascending, e.g. the first element
	Dates []time.Time `json:"dates"`
	// From these, the dates above can be calculated
	SchedulePayload
	LastRun   *time.Time `json:"lastRun,omitempty"`
	LastError error      `json:"lastError,omitempty"`

	// TODO: implement runs, which hold historical information about a run.
	// These should hold a reference to the object of which is was created with,
	// But since that object might change in the future, we want to store those
	// parameters normalized in a Run too.
	// RunIDS []string
}

func (s Schedule) NextRun() *time.Time {
	// TODO: this should calculate based on other parameters
	t := s.StartDate
	if s.LastRun != nil && s.LastRun.After(t) {
		return nil
	}
	return &t
}

// Indicates whether this schedule should be run based on run-parameters
func (s Schedule) ShouldRun() bool {
	t := s.NextRun()
	if t == nil {
		return false
	}

	now := time.Now()
	if s.LastError != nil {
		// If there was an error, we postpone the run a bit.
		now = now.Add(-time.Minute)
	}
	return t.Before(now)
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
