package errors

import (
	goerrors "errors"

	"utils/errors"
)

type Error struct {
	Type Type
	err  error
}

func (e Error) Error() string {
	return e.err.Error()
}

func Is(err, target error) bool {
	return goerrors.Is(err, target)
}

func Join(errs ...error) error {
	return goerrors.Join(errs...)
}

func _new(eType Type, args ...any) error {
	return &Error{err: errors.NewCustom(2, args...), Type: eType}
}

func GetCode(err error) Type {
	if err == nil {
		return FatalError
	}

	var ourErr *Error
	if ok := goerrors.As(err, &ourErr); !ok {
		return FatalError
	}

	return ourErr.Type
}
