// swagger:route POST /dry-dynamic dynamic dryDynamic
// Runs a dry-run for dynamic requests
// responses:
//   200: dryDynamicResponse
//   500: apiError

package docs

import (
	"github.com/runar-rkmedia/gabyoall/auth"
)

// Result of dry-response
// swagger:response dryDynamicResponse
type dryDynamicResponse struct {
	// in:body
	Body struct {
		Result auth.DynamicAuth `json:"result"`
		Error  string           `json:"error"`
	}
}

// swagger:parameters dryDynamic
type dryDynamic struct {
	// in:body
	Body auth.DynamicAuth
}
