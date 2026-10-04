package domain

import "errors"

var (
	ErrValidation = errors.New("validation failed")
	ErrNotFound   = errors.New("not found")
)

// ValidationError marks an error without changing its human-readable message.
func ValidationError(err error) error {
	return &classifiedError{kind: ErrValidation, cause: err}
}

// NotFoundError preserves the underlying error, including sql.ErrNoRows.
func NotFoundError(err error) error {
	return &classifiedError{kind: ErrNotFound, cause: err}
}

type classifiedError struct {
	kind  error
	cause error
}

func (e *classifiedError) Error() string   { return e.cause.Error() }
func (e *classifiedError) Unwrap() []error { return []error{e.kind, e.cause} }
