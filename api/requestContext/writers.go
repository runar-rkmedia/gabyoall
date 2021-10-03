package requestContext

import (
	"net/http"
	"strings"
)

type OutputKind = int

const (
	OutputJson OutputKind = iota + 1
	OutputYaml
	OutputToml
)

func WriteAuto(output interface{}, err error, errCode ErrorCodes, r *http.Request, rw http.ResponseWriter) {
	if err != nil {
		WriteError(err.Error(), CodeErrEndpoint, r, rw)
		return
	}

	WriteOutput(false, output, r, rw)
}
func WriteErr(err error, code ErrorCodes, r *http.Request, rw http.ResponseWriter) {
	WriteError(err.Error(), code, r, rw)
	return
}
func WriteError(msg string, code ErrorCodes, r *http.Request, rw http.ResponseWriter) {
	ae := ApiError{msg, string(code)}
	WriteOutput(true, ae, r, rw)
	switch code {
	case CodeErrMethodNotAllowed:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

	return

}

// attempts to guess at what kind of output the user most likely wants
func WantedOutputFormat(r *http.Request) OutputKind {
	if o := contentType(r.Header.Get("Accept")); o > 0 {
		return o
	}
	// If the content-type is set, the user probably wants the same type.
	if o := contentType(r.Header.Get("Content-Type")); o > 0 {
		return o
	}
	// Fallback to a readable format.
	return OutputToml
}

func contentType(kind string) OutputKind {
	switch {
	case strings.Contains(kind, "application/json"):
		return OutputJson
	case strings.Contains(kind, "yaml"):
		return OutputYaml
	case strings.Contains(kind, "toml"):
		return OutputToml
	}
	return 0
}
