package requestContext

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/naoina/toml"
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
		return toml.Unmarshal(b, j)
	case OutputJson:
		return json.Unmarshal(b, j)
	case OutputYaml:
		return yaml.Unmarshal(b, j)
	}
	// Fallback to a readable format.
	return toml.Unmarshal(b, j)
}
