// swagger:route GET /endpoint endpoint listEndpoints
// Endpoints are remote or local urls with configured handling for authorization etc.
// responses:
//   200: endpointsResponse
//   500: apiError

// swagger:route GET /endpoint/{id} endpoint getEndpoint
// Returns a single endpoint by id.
// responses:
//   200: endpointResponse
//   404: apiError
//   500: apiError

// swagger:route POST /endpoint/ endpoint createEndpoint
// Create a new endpoint
// responses:
//   200: okResponse
//   404: apiError
//   500: apiError
package docs

import (
	"github.com/runar-rkmedia/gabyoall/api/types"
)

// Lists endpoints registered
// swagger:response endpointsResponse
type endpointsResponse struct {
	// in:body
	Body []types.EndpointEntity
}

// Returns single endpoint
// swagger:response endpointResponse
type endpointResponse struct {
	// in:body
	Body types.EndpointEntity
}

// swagger:parameters getEndpoint
type GetParams struct {
	// minLength: 3
	// maxLength: 40
	// in: path
	// example: abc123
	ID string `json:"id"`
}

// swagger:parameters createEndpoint
type createEndpoint struct {
	// in: body
	// required: true
	Body types.EndpointPayload
}
