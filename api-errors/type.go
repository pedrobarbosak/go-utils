package errors

type Type int

const (
	FatalError Type = iota
	InputError
	NotFoundError
	ForbiddenError
	UnauthorizedError
	UnprocessableEntityError
	ConflictError
	UnsupportedMediaType
)

func NewInput(args ...any) error {
	return _new(InputError, args...)
}

func NewNotFound(args ...any) error {
	return _new(NotFoundError, args...)
}

func NewFatal(args ...any) error {
	return _new(FatalError, args...)
}

func NewConflict(args ...any) error {
	return _new(ConflictError, args...)
}

func NewForbidden(args ...any) error {
	return _new(ForbiddenError, args...)
}

func NewUnauthorized(args ...any) error {
	return _new(UnauthorizedError, args...)
}

func NewUnsupported(args ...any) error {
	return _new(UnsupportedMediaType, args...)
}

func NewUnprocessable(args ...any) error {
	return _new(UnprocessableEntityError, args...)
}
