package requestContext

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/jmespath/go-jmespath"
	"github.com/naoina/toml"
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
func WriteOutput(isError bool, output interface{}, r *http.Request, rw http.ResponseWriter) {
	o := WantedOutputFormat(r)
	switch o {
	case OutputJson:
		rw.Header().Set("Content-Type", "application/json")
	case OutputToml:
		rw.Header().Set("Content-Type", "application/toml")
	case OutputYaml:
		rw.Header().Set("Content-Type", "text/vnd.yaml")
	}
	if !isError {
		jmesPath := r.Header.Get("JMES-path")
		if jmesPath != "" {
			b, err := json.Marshal(output)
			if err != nil {
				WriteErr(err, CodeErrMarhal, r, rw)
				return
			}
			var JSON map[string]interface{}
			err = json.Unmarshal(b, &JSON)
			if err != nil {
				WriteErr(err, CodeErrUnmarshal, r, rw)
				return
			}
			result, err := jmespath.Search(jmesPath, JSON)
			if err != nil {
				WriteErr(fmt.Errorf("failed in jmes-path '%s': %w", jmesPath, err), CodeErrJmesPath, r, rw)
				return
			}

			if o == OutputToml {
				// Toml does not support outputting primitives, so we cheat a bit.
				switch result.(type) {
				case string:
					// rw.Write([]byte(result.(string)))
					// return
				case int, int8, int16, int32, int64:
					rw.Write([]byte(fmt.Sprintf("%d", result)))
					return
				}
			}
			// Technically not an error, but we dont want to run jmes-path-again
			WriteOutput(true, result, r, rw)
			return
		}
	}
	switch o {
	case OutputJson:
		b, err := json.Marshal(output)
		if err != nil {
			WriteErr(err, CodeErrMarhal, r, rw)
			return
		}
		rw.Write(b)
	case OutputYaml:
		b, err := yaml.Marshal(output)
		if err != nil {
			WriteErr(err, CodeErrMarhal, r, rw)
			return
		}
		rw.Write(b)
	case OutputToml:
		b, err := toml.Marshal(output)
		if err != nil {
			WriteErr(err, CodeErrMarhal, r, rw)
			return
		}
		rw.Write(b)
	}
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
