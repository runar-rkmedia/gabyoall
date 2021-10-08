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

	// Schedules() (es map[string]ScheduleEntity, err error)
	// Schedule(id string) (ScheduleEntity, error)
	// CreateSchedule(e SchedulePayload) (ScheduleEntity, error)
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

type EndpointToRequestEntity struct {
	Entity
	EndpointToRequestPayload
}

// This should perhaps be a configuration, with run-settings?
type EndpointToRequestPayload struct {
	RequestID  string `json:"requestID"`
	EndpointID string `json:"endpointID"`
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
	ScheduleInput SchedulePayload
	// TODO: implement runs, which hold historical information about a run.
	// These should hold a reference to the object of which is was created with,
	// But since that object might change in the future, we want to store those
	// parameters normalized in a Run too.
	// RunIDS []string
}
type SchedulePayload struct {
	EndpointToRequestPayload
	// If set to a positive value, the scheduler will not schedule more than this total concurrency
	// when starting this job, and when it is running.
	//
	// Some jobs might have been configured to run very slowly, with low concurrency,
	// high wait-times and can therefore run alongside other such jobs.
	MaxInterJobConcurrency bool      `json:"maxInterJobConcurrency"`
	Label                  string    `json:"label"`
	StartDate              time.Time `json:"start_date" validate:"required"`
	Frequency              Frequency `json:"frequency,omitempty"`
	Multiplier             int       `json:"multiplier,omitempty"`
	Offsets                []int     `json:"offsets,omitempty"`
}

type Frequency int8

const (
	FrequencyMinute Frequency = iota
	FrequencyHour   Frequency = iota
	FrequencyDay    Frequency = iota
	FrequencyWeek   Frequency = iota
	FrequencyMonth  Frequency = iota
)
