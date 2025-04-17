package utils

type Remark struct {
	Remark ErrorDetails
}

type ErrorDetails struct {
	// Message (required) is the user-defined error message.
	// E.g. "account with no rekening not found".
	Message string

	// Code (required) is the user-defined error code string that follows
	// E.g. "ACCOUNT_NOT_FOUND".
	Code string

	// Field (optional) is the related field the error occurred on, if any.
	// field name
	Field string

	// Object (optional) is the related object of, if any.
	Object interface{}
}

func NewRemark(message, code, field string, object interface{}) *Remark {
	return &Remark{
		Remark: *NewErrorDetailsWithObject(message, code, field, object),
	}
}

// NewErrorDetailsWithObject creates a new ErrorDetails struct with an associated object.
func NewErrorDetailsWithObject(message, code, field string, object interface{}) *ErrorDetails {
	return &ErrorDetails{
		Message: message,
		Code:    code,
		Field:   field,
		Object:  object,
	}
}

func (e *Remark) Error() string {
	return e.Remark.Message
}
