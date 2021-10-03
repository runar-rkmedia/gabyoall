// Package classification Gobyoall API.
//
// Api fo gobyoall-api.
// <a href="https://insomnia.rest/run/?label=&uri=" target="_blank"><img src="https://insomnia.rest/images/run.svg" alt="Run in Insomnia"></a>
//     Title: "Gobyoall"
//     Schemes: http
//     BasePath: /api
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Host: localhost
//
//     Consumes:
//     - application/json
//     - text/vnd.yaml
//     - application/toml
//
//     Produces:
//     - application/json
//     - text/vnd.yaml
//     - application/toml
//
//
// swagger:meta
package docs

// swagger:response apiError
type apiError struct {
	// in:body
	Body ApiError
}

// swagger:response okResponse
type okResponse struct {
	// in:body
	Body OkResponse
}

// swagger:response createResponse
type createResponse struct {
	// in:body
	Body CreateResponse
}

// FIXME: duplicates, use import
type ApiError struct {
	Message string
	Code    int
}
type OkResponse struct {
	Ok bool `json:"ok"`
}
type CreateResponse struct {
	OkResponse
	ID string `json:"id"`
}
