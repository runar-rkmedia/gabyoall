// swagger:route GET /schedule schedule listSchedules
// responses:
//   200: schedulesResponse
//   500: apiError

// swagger:route GET /schedule/{id} schedule getSchedule
// Returns a single schedule by id.
// responses:
//   200: scheduleResponse
//   404: apiError
//   500: apiError

// swagger:route POST /schedule/ schedule createSchedule
// Create a new schedule
// responses:
//   200: scheduleResponse
//   404: apiError
//   500: apiError

// swagger:route PUT /schedule/{id} schedule updateSchedule
// Updates a schedule
// responses:
//   200: scheduleResponse
//   404: apiError
//   500: apiError
package docs

import (
	"github.com/runar-rkmedia/gabyoall/api/types"
)

// Lists schedules registered
// swagger:response schedulesResponse
type schedulesResponse struct {
	// in:body
	Body []types.ScheduleEntity
}

// Returns single schedule
// swagger:response scheduleResponse
type scheduleResponse struct {
	// in:body
	Body types.ScheduleEntity
}

// swagger:parameters getSchedule
type getScheduleParams struct {
	// minLength: 3
	// maxLength: 40
	// in: path
	// example: abc123
	ID string `json:"id"`
}

// swagger:parameters createSchedule
type createSchedule struct {
	// in: body
	// required: true
	Body types.SchedulePayload
}

// swagger:parameters updateSchedule
type updateSchedule struct {
	// in: body
	// required: true
	Body types.SchedulePayload
}
