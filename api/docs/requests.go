// swagger:route GET /request request listRequests
// responses:
//   200: requestsResponse
//   500: apiError

// swagger:route GET /request/{id} request getRequest
// Returns a single request by id.
// responses:
//   200: requestResponse
//   404: apiError
//   500: apiError

// swagger:route POST /request/ request createRequest
// Create a new request
// responses:
//   200: okResponse
//   404: apiError
//   500: apiError
package docs

import (
	"github.com/runar-rkmedia/gabyoall/api/types"
)

// Lists requests registered
// swagger:response requestsResponse
type requestsResponse struct {
	// in:body
	Body []types.RequestEntity
}

// Returns single request
// swagger:response requestResponse
type requestResponse struct {
	// in:body
	Body types.RequestEntity
}

// swagger:parameters getRequest
type getRequestParams struct {
	// minLength: 3
	// maxLength: 40
	// in: path
	// example: abc123
	ID string `json:"id"`
}

// swagger:parameters createRequest
type createRequest struct {
	// in: body
	// required: true
	Body types.RequestPayload
}
