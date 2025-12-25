package errs

import "strings"

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ActionType string

const (
	ActionTypeRedirect ActionType = "redirect"
)

type Action struct {
	Message string     `json:"message"`
	Value   string     `json:"value"`
	Type    ActionType `json:"type"`
}

type HttpError struct {
	Code     string       `json:"code"`
	Message  string       `json:"message"`
	Status   int          `json:"status"`
	Override bool         `json:"override"`
	Errors   []FieldError `json:"errors"` // Field level errors
	Action   *Action      `json:"action"` // Action to be taken
}

func (e *HttpError) Error() string {
	return e.Message
}

func (e *HttpError) Is(target error) bool {
	_, ok := target.(*HttpError)
	return ok
}

func (e *HttpError) withMessage(message string) *HttpError {
	return &HttpError{
		Code:     e.Code,
		Message:  e.Message,
		Status:   e.Status,
		Override: e.Override,
		Errors:   e.Errors,
		Action:   e.Action,
	}
}

func MakeUpperCaseWithUnderScores(str string) string {
	return strings.ToUpper(strings.ReplaceAll(str, " ", "_"))
}
