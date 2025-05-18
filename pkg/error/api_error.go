package error

import "errors"

type Error struct {
	Err  error `json:"error"`
	Code int   `json:"code"`
}

func New(msg string, code int) Error {
	return Error{Err: errors.New(msg), Code: code}
}

func (e Error) Error() string {
	return e.Err.Error()
}

var (
	ErrNotFound              = New("not found", 400)
	ErrUnauthorized          = New("unauthorized", 401)
	ErrForbidden             = New("forbidden", 403)
	ErrInternal              = New("internal server error", 500)
	ErrUserNotFound          = New("user not found", 404)
	ErrInvalidEmail          = New("invalid email", 422)
	ErrUserAlreadyExists     = New("user already exists", 409)
	ErrUsernameAlreadyExists = New("username already exists", 422)
	ErrEmailAlreadyExists    = New("email already exists", 422)
	ErrInvalidPassword       = New("invalid password", 422)
	ErrInvalidToken          = New("invalid token", 422)
)
