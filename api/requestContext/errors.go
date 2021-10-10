package requestContext

import "fmt"

type ErrorCodes string
type Error struct {
	Code    ErrorCodes
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

func NewError(msg string, code ErrorCodes) Error {
	return Error{code, msg}
}

var (
	ErrIDNonValid = NewError("Id was not valid", CodeErrIDNonValid)
	ErrIDTooLong  = NewError("Id was too long", CodeErrIDTooLong)
	ErrIDEmpty    = NewError("Id was empty", CodeErrIDEmpty)
)

const (
	CodeErrEndpoint         ErrorCodes = "Endpoint"
	CodeErrRequest          ErrorCodes = "Request"
	CodeErrMethodNotAllowed ErrorCodes = "MethodNotAllowed"
	CodeErrNoRoute          ErrorCodes = "NoRoute"
	CodeErrReadBody         ErrorCodes = "ReadBody"
	CodeErrMarhal           ErrorCodes = "Marhal"
	CodeErrUnmarshal        ErrorCodes = "Unmarshal"
	CodeErrJmesPath         ErrorCodes = "JmesPath"
	CodeErrJmesPathMarshal  ErrorCodes = "JmesPathMarshal"
	CodeErrDBCreateEndpoint ErrorCodes = "DBCreateEndpoint"
	CodeErrSchedule         ErrorCodes = "DBCreateSchedule"
	CodeErrDBCreateRequest  ErrorCodes = "DBCreateRequest"
	CodeErrInputValidation  ErrorCodes = "InputValidation"
	CodeErrIDNonValid       ErrorCodes = "IDNonValid"
	CodeErrIDTooLong        ErrorCodes = "IDTooLong"
	CodeErrIDEmpty          ErrorCodes = "IDEmpty"
)

type ApiError struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}
