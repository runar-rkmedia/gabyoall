package requestContext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/jmespath/go-jmespath"
	"github.com/pelletier/go-toml"
)

func UnmarshalRequestBody(r *http.Request, j interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return UnmarshalRequestBytes(r, b, j)
}
func UnmarshalRequestBytes(r *http.Request, b []byte, j interface{}) error {
	return UnmarshalWithKind(WantedOutputFormat(r), b, j)
}
func UnmarshalWithKind(kind OutputKind, b []byte, j interface{}) error {
	switch kind {
	case OutputToml:

		// TODO: test this. (probably does not work because it does not read json-tags)
		return toml.Unmarshal(b, j)
	case OutputJson:
		return json.Unmarshal(b, j)
	case OutputYaml:
		return yaml.Unmarshal(b, j)
	}
	// Fallback to a readable format.
	return toml.Unmarshal(b, j)
}

func WriteOutput(isError bool, statusCode int, output interface{}, r *http.Request, rw http.ResponseWriter) error {
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
				return err
			}
			var JSON map[string]interface{}
			err = json.Unmarshal(b, &JSON)
			if err != nil {
				WriteErr(err, CodeErrUnmarshal, r, rw)
				return err
			}
			result, err := jmespath.Search(jmesPath, JSON)
			if err != nil {
				WriteErr(fmt.Errorf("failed in jmes-path '%s': %w", jmesPath, err), CodeErrJmesPath, r, rw)
				return err
			}

			if o == OutputToml {
				// Toml does not support outputting primitives, so we cheat a bit.
				switch result.(type) {
				case string:
					// rw.Write([]byte(result.(string)))
					// return
				case int, int8, int16, int32, int64:
					rw.Write([]byte(fmt.Sprintf("%d", result)))
					return err
				}
			}
			// Technically not an error, but we dont want to run jmes-path-again
			return WriteOutput(true, statusCode, result, r, rw)
		}
	}
	if statusCode >= 100 {
		rw.WriteHeader(statusCode)
	}
	switch o {
	case OutputJson:
		b, err := json.Marshal(output)
		if err != nil {
			WriteErr(err, CodeErrMarhal, r, rw)
			return err
		}
		rw.Write(b)
	case OutputYaml:
		b, err := yaml.Marshal(output)
		if err != nil {
			WriteErr(err, CodeErrMarhal, r, rw)
			return err
		}
		rw.Write(b)
	case OutputToml:
		// toml does not use json-tags.
		// This is basically the same as what yaml does
		// E.g. it first uses json-marshaller/unmarshal then toml-marshal.
		jb, err := json.Marshal(output)
		if err != nil {
			WriteErr(err, CodeErrMarhal, r, rw)
			return err
		}

		var JSON map[string]interface{}
		err = json.Unmarshal(jb, &JSON)
		if err != nil {
			WriteErr(err, CodeErrMarhal, r, rw)
			return err
		}
		b, err := toml.Marshal(JSON)
		if err != nil {
			WriteErr(err, CodeErrMarhal, r, rw)
			return err
		}
		rw.Write(b)
	}
	return nil
}
