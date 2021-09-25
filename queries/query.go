package queries

type GraphQLQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
	// For some reason, the server does not like operationName.
	OperationName string `json:"-"` //`json:"operationName"`
}
