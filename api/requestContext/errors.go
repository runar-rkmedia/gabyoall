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
	CodeErrEndpoint         ErrorCodes = "ErrEndpoint"
	CodeErrRequest          ErrorCodes = "ErrRequest"
	CodeErrMethodNotAllowed ErrorCodes = "ErrMethodNotAllowed"
	CodeErrNoRoute          ErrorCodes = "ErrNoRoute"
	CodeErrReadBody         ErrorCodes = "ErrReadBody"
	CodeErrMarhal           ErrorCodes = "ErrMarshal"
	CodeErrUnmarshal        ErrorCodes = "ErrUnmarshal"
	CodeErrJmesPath         ErrorCodes = "ErrJmesPath"
	CodeErrJmesPathMarshal  ErrorCodes = "ErrJmesPathMarshal"
	CodeErrDBCreateEndpoint ErrorCodes = "ErrDBCreateEndpoint"
	CodeErrSchedule         ErrorCodes = "ErrDBCreateSchedule"
	CodeErrDBCreateRequest  ErrorCodes = "ErrDBCreateRequest"
	CodeErrDBCreateSchedule ErrorCodes = "ErrDBCreateSchedule"
	CodeErrInputValidation  ErrorCodes = "ErrInputValidation"
	CodeErrIDNonValid       ErrorCodes = "ErrIDNonValid"
	CodeErrIDTooLong        ErrorCodes = "ErrIDTooLong"
	CodeErrIDEmpty          ErrorCodes = "ErrIDEmpty"
)

type ApiError struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}
