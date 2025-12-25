package errs

import "net/http"

func NewUnauthorizedError(message string, override bool) *HttpError {
	return &HttpError{
		Code:     MakeUpperCaseWithUnderScores(http.StatusText(http.StatusUnauthorized)),
		Message:  message,
		Status:   http.StatusUnauthorized,
		Override: override,
	}
}

func NewForbiddenError(message string, override bool) *HttpError {
	return &HttpError{
		Code:     MakeUpperCaseWithUnderScores(http.StatusText(http.StatusForbidden)),
		Message:  message,
		Status:   http.StatusForbidden,
		Override: override,
	}
}

func NewBadRequestError(message string, override bool, code *string, errors []FieldError, action *Action) *HttpError {
	formattedCode := MakeUpperCaseWithUnderScores(http.StatusText(http.StatusBadRequest))

	if code != nil {
		formattedCode = *code
	}

	return &HttpError{
		Code:     formattedCode,
		Message:  message,
		Status:   http.StatusBadRequest,
		Override: override,
		Errors:   errors,
		Action:   action,
	}
}

func NewNotFoundError(message string, override bool, code *string) *HttpError {
	formattedCode := MakeUpperCaseWithUnderScores(http.StatusText(http.StatusNotFound))

	if code != nil {
		formattedCode = *code
	}

	return &HttpError{
		Code:     formattedCode,
		Message:  message,
		Status:   http.StatusNotFound,
		Override: override,
	}
}

func NewInternalServerError() *HttpError {
	return &HttpError{
		Code:     MakeUpperCaseWithUnderScores(http.StatusText(http.StatusInternalServerError)),
		Message:  http.StatusText(http.StatusInternalServerError),
		Status:   http.StatusInternalServerError,
		Override: false,
	}
}

func ValidationError(err error) *HttpError {
	return NewBadRequestError("Validation failed: "+err.Error(), false, nil, nil, nil)
}
