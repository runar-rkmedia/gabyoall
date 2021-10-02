package requests

type Request struct {
	// Will only be used if Query is unset.
	Body      interface{}            `json:"-"`
	Query     string                 `json:"query,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
	Headers   map[string]string      `json:"-"`
	// For some reason, the server does not like operationName.
	OperationName string `json:"-"` //`json:"operationName"`
	Method        string `json:"-"`
}
