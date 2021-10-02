package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/jmespath/go-jmespath"
)

type DynamicAuth struct {
	Requests  []DynamicRequest
	HeaderKey string
}

type DynamicRequest struct {
	Method         string
	Uri            string
	Headers        map[string]string
	JsonRequest    bool
	JsonResponse   bool
	ResultJmesPath string
	Body           interface{}
}
type DynamicResponse struct {
	result   interface{}
	response *http.Response
}
type DynamicAuthResult struct {
	Token     string
	HeaderKey string
}

func (da *DynamicAuth) Retrieve() (DynamicAuthResult, error) {
	results := make([]DynamicResponse, len(da.Requests))
	var dyn DynamicAuthResult
	// TODO: these should be piped into eachother somehow, so that the response and result can be used in templating or something.
	for i := 0; i < len(da.Requests); i++ {
		dynRes, err := da.Requests[i].Do()
		if err != nil {
			return dyn, err
		}
		results[i] = *dynRes
	}
	last := results[len(results)-1]
	if last.result == nil {
		return dyn, fmt.Errorf("result was nil")
	}
	if str, ok := last.result.(string); ok {
		dyn = DynamicAuthResult{
			HeaderKey: da.HeaderKey,
			Token:     str,
		}
		return dyn, nil
	}
	return dyn, fmt.Errorf("Unhandled result in dynamic auth")
}

func (dr DynamicRequest) Do() (*DynamicResponse, error) {
	var reader io.Reader
	if dr.Method == "" {
		return nil, fmt.Errorf("missing field method")
	}
	if dr.Uri == "" {
		return nil, fmt.Errorf("missing field uri")
	}
	b, _ := yaml.Marshal(dr)
	fmt.Println(string(b))
	// return nil, fmt.Errorf("FFFFFFFFFFF")
	if dr.JsonRequest && dr.Body != nil {
		JSON, ok := dr.Body.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to convert body to map[string]interface: %s", dr.Body)
		}
		b, err := json.Marshal(JSON)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		reader = bytes.NewReader(b)

	}
	r, err := http.NewRequest(dr.Method, dr.Uri, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()
	contentType := res.Header.Get("Content-Type")
	dyn := DynamicResponse{
		response: res,
	}
	if res.StatusCode >= 400 {
		return &dyn, fmt.Errorf("%d %s %s", res.StatusCode, res.Status, contentType)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &dyn, fmt.Errorf("failed to read body of request: %w", err)
	}
	var JSON map[string]interface{}
	if strings.Contains(contentType, "json") {
		if err := json.Unmarshal(body, &JSON); err != nil {
			return &dyn, fmt.Errorf("failed to unmarshal json %w", err)
		}
	}
	if dr.ResultJmesPath != "" {
		result, err := jmespath.Search(dr.ResultJmesPath, JSON)
		if err != nil {
			return &dyn, fmt.Errorf("failed to perform jmesPath %s: %w", dr.ResultJmesPath, err)
		}
		dyn.result = result
	}
	return &dyn, nil

}
