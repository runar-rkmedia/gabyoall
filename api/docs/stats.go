// swagger:route GET /stat stat listStats
// Stats are remote or local urls with configured handling for authorization etc.
// responses:
//   200: statsResponse
//   500: apiError

// swagger:route GET /stat/{id} stat getStat
// Returns a single stat by id.
// responses:
//   200: statResponse
//   404: apiError
//   500: apiError

package docs

import (
	"github.com/runar-rkmedia/gabyoall/api/types"
)

// Lists stats registered
// swagger:response statsResponse
type statsResponse struct {
	// in:body
	Body []types.CompactRequestStatisticsEntity
}

// Returns single stat
// swagger:response statResponse
type statResponse struct {
	// in:body
	Body types.CompactRequestStatisticsEntity
}

// swagger:parameters getStat
type getStatParams struct {
	// minLength: 3
	// maxLength: 40
	// in: path
	// example: abc123
	ID string `json:"id"`
}
