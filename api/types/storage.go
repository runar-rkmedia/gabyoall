package types

import (
	"time"
)

type Storage interface {
	Endpoints() (es map[string]EndpointEntity, err error)
	Endpoint(id string) (EndpointEntity, error)
	CreateEndpoint(e EndpointPayload) (EndpointEntity, error)
	UpdateEndpoint(id string, p EndpointPayload) (EndpointEntity, error)
	SoftDeleteEndpoint(id string) (EndpointEntity, error)

	Requests() (es map[string]RequestEntity, err error)
	Request(id string) (RequestEntity, error)
	CreateRequest(e RequestPayload) (RequestEntity, error)
	SoftDeleteRequest(id string) (RequestEntity, error)

	Schedules() (es map[string]ScheduleEntity, err error)
	Schedule(id string) (ScheduleEntity, error)
	CreateSchedule(e SchedulePayload) (ScheduleEntity, error)
	UpdateSchedule(id string, p Schedule) (ScheduleEntity, error)
	SoftDeleteSchedule(id string) (ScheduleEntity, error)

	CompactStats() (es map[string]StatEntity, err error)
	CleanCompactStats() (err error)

	Size() (int64, error)

	// These are only used internally.

	CreateCompactStats(id string, createdAt time.Time, p StatPayload) error

	UpdateCompactStats(id string, createdAt time.Time, p StatPayload) error
}
