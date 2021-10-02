package requests

type ErrorType string

var (
	GQLError        ErrorType = "GQLError"
	NonOK           ErrorType = "NonOK"
	ServerTestError ErrorType = "ServerTestError"
	Unknwon         ErrorType = "UnknownError"
)
