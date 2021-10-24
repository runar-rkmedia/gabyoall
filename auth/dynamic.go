package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/runar-rkmedia/gabyoall/utils"
)

type DynamicAuth struct {
	Requests  []DynamicRequest `json:"requests" validate:"required,min=1,max=100"`
	HeaderKey string           `json:"headerKey"`
}

type DynamicRequest struct {
	Method         string            `json:"method,omitempty" validate:"required"`
	Uri            string            `json:"uri,omitempty" validate:"required"`
	Headers        map[string]string `json:"headers,omitempty"`
	JsonRequest    bool              `json:"json_request,omitempty"`
	JsonResponse   bool              `json:"json_response,omitempty"`
	ResultJmesPath string            `json:"result_jmes_path,omitempty"`
	Body           interface{}       `json:"body,omitempty"`
}
type DynamicResponse struct {
	result   interface{}
	response *DynamicHttpResponse
}

type DynamicHttpResponse struct {
	Headers    http.Header
	StatusCode int
	BodyRaw    []byte
	BodyJson   map[string]interface{}
}
type DynamicAuthResult struct {
	Token     string
	HeaderKey string
	Responses []*DynamicHttpResponse
}

func (da *DynamicAuth) Retrieve() (DynamicAuthResult, error) {
	var dyn DynamicAuthResult
	if len(da.Requests) == 0 {
		return dyn, fmt.Errorf("a minimum of one requests are required for dynamicAuth")
	}
	results := make([]DynamicResponse, len(da.Requests))
	dyn.Responses = make([]*DynamicHttpResponse, len(da.Requests))
	// TODO: these should be piped into eachother somehow, so that the response and result can be used in templating or something.
	for i := 0; i < len(da.Requests); i++ {
		dynRes, err := da.Requests[i].Do()
		if dynRes != nil && dynRes.response != nil {
			dyn.Responses[i] = dynRes.response
		}
		if err != nil {
			return dyn, err
		}
		results[i] = *dynRes
	}
	last := results[len(results)-1]
	if last.result == nil {
		return dyn, fmt.Errorf("result was nil")
	}
	switch last.result.(type) {
	case string:
		dyn.Token = last.result.(string)
	case []byte:
		dyn.Token = string(last.result.([]byte))
	case int:
		dyn.Token = fmt.Sprintf("%d", last.result.(int))
	case float64:
		dyn.Token = fmt.Sprintf("%f", last.result.(float64))
	default:
		return dyn, fmt.Errorf("Unhandled result in dynamic auth: %#v", last.result)
	}
	dyn.HeaderKey = da.HeaderKey
	if dyn.Token == "" {
		return dyn, fmt.Errorf("result was empty")
	}
	return dyn, nil
}

func (dr DynamicRequest) Do() (*DynamicResponse, error) {
	var reader io.Reader
	if dr.Method == "" {
		return nil, fmt.Errorf("missing field method")
	}
	if dr.Uri == "" {
		return nil, fmt.Errorf("missing field uri")
	}
	if dr.JsonRequest && dr.Body != nil {
		var b []byte
		var err error
		switch dr.Body.(type) {
		case string:
			b = []byte(dr.Body.(string))
		case map[string]interface{}:
			b, err = json.Marshal(dr.Body)
		default:
			fmt.Println("type")
		}
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
		response: &DynamicHttpResponse{
			Headers:    res.Header,
			StatusCode: res.StatusCode,
		},
	}
	if res.StatusCode >= 400 {
		return &dyn, fmt.Errorf("%d %s %s", res.StatusCode, res.Status, contentType)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &dyn, fmt.Errorf("failed to read body of request: %w", err)
	}
	defer res.Body.Close()
	dyn.response.BodyRaw = body
	if strings.Contains(contentType, "json") {
		if err := json.Unmarshal(body, &dyn.response.BodyJson); err != nil {
			return &dyn, fmt.Errorf("failed to unmarshal json %w", err)
		}
	}
	if dr.ResultJmesPath != "" {
		vars := struct {
			Response DynamicHttpResponse
		}{
			Response: *dyn.response,
		}
		templateResult := utils.RunTemplating(logger.GetLogger("dynamic-request-templating"), dr.ResultJmesPath, "jmes-path", vars)
		dyn.result = templateResult
	}
	return &dyn, nil

}
