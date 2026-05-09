// pkg/errors/errors.go
package errors

import "errors"

// Business error codes
const (
	CodeSuccess       = 20000
	CodeBadParams     = 40000
	CodeUnauthorized  = 40100
	CodeForbidden     = 40300
	CodeNotFound      = 50004
	CodeConflict      = 50008
	CodeInternalError = 50000
)

// Domain errors — compare with errors.Is
var (
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrConflict       = errors.New("resource conflict")
	ErrBadParams      = errors.New("bad params")
	ErrVersionInUse   = errors.New("resume version is referenced by an application")
	ErrEmailCodeWrong = errors.New("email verification code is incorrect or expired")
	ErrDuplicate      = errors.New("duplicate entry")
)
