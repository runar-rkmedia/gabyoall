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
	CodeErrEndpoint         ErrorCodes = "Error: Regarding type Endpoint"
	CodeErrRequest          ErrorCodes = "Error: Regarding type Request"
	CodeErrMethodNotAllowed ErrorCodes = "Error: HTTP-Method is not allowed"
	CodeErrNoRoute          ErrorCodes = "Error: No route matched for this http-path"
	CodeErrReadBody         ErrorCodes = "Error: Failed to read body"
	CodeErrMarhal           ErrorCodes = "Error: Failed to marshal"
	CodeErrUnmarshal        ErrorCodes = "Error: Failed to unmarshal"
	CodeErrJmesPath         ErrorCodes = "Error: JmesPath"
	CodeErrJmesPathMarshal  ErrorCodes = "Error: JmesPathMarshal"

	CodeErrInputValidation ErrorCodes = "Error: General input validation"
	CodeErrIDNonValid      ErrorCodes = "Error: ID not valid"
	CodeErrIDTooLong       ErrorCodes = "Error: ID is too long"
	CodeErrIDEmpty         ErrorCodes = "Error: ID was Empty"

	CodeErrDBUpdateSchedule ErrorCodes = "Error: Database Update Schedule"
	CodeErrDBUpdateEndpoint ErrorCodes = "Error: Database Update Endpoint"
	CodeErrDBUpdateRequest  ErrorCodes = "Error: Database Update Request"
	CodeErrDBDeleteEndpoint ErrorCodes = "Error: Database Delete Endpoint"
	CodeErrDBDeleteRequest  ErrorCodes = "Error: Database Delete Request"
	CodeErrDBDeleteSchedule ErrorCodes = "Error: Database Delete Schedule"
	CodeErrDBCreateEndpoint ErrorCodes = "Error: Database Create Endpoint"
	CodeErrSchedule         ErrorCodes = "Error: Database Create Schedule"
	CodeErrDBCreateRequest  ErrorCodes = "Error: Database Create Request"
	CodeErrDBCreateSchedule ErrorCodes = "Error: Database Create Schedule"
)

type ApiError struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}
