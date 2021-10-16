// swagger:route GET /serverInfo serverInfo serverInfo
// Various information about the server.
// responses:
//   200: serverInfoResponse
//   500: apiError
package docs

import (
	"github.com/runar-rkmedia/gabyoall/api/types"
)

// Server info
// swagger:response serverInfoResponse
type serverInfoResponse struct {
	// in:body
	Body []types.ServerInfo
}
