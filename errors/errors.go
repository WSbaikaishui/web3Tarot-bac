package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type code uint32

func (c code) MarshalJSON() ([]byte, error) {
	return codeToID[uint32(c)], nil
}

type ErrorInfo struct {
	Code    code   `json:"code"`
	Message string `json:"message"`
	//Data    any    `json:"data"`
}

func NewErrorInfo(major, minor int, message string) ErrorInfo {
	return ErrorInfo{Code: code(uint32(major)*10000 + uint32(minor)), Message: message}
}

func (err ErrorInfo) Error() string {
	return fmt.Sprintf("code: %d, message: %s", err.Code, err.Message)
}

func (err ErrorInfo) StatusCode() int {
	return int(err.Code / 10000)
}

var codeToID = map[uint32]json.RawMessage{
	4000100: []byte(`"InvalidParameter"`),
	4000101: []byte(`"TokenExpired"`),
	4000102: []byte(`"Unauthorized"`),
	4000103: []byte(`"Forbidden"`),
	4000104: []byte(`"NotSatisfied"`),
	4000105: []byte(`"NotFound"`),
	4000107: []byte(`"InvalidSignature"`),
	5000100: []byte(`"InternalError"`),
	5000101: []byte(`"Unimplemented"`),
}

var (
	ErrInvalidParameter = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusBadRequest, 100, msg) }
	ErrTokenExpired     = func() ErrorInfo { return NewErrorInfo(http.StatusBadRequest, 101, "token is expired") }
	ErrUnauthorized     = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusBadRequest, 102, msg) }
	ErrForbidden        = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusBadRequest, 103, msg) }
	ErrNotSatisfied     = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusBadRequest, 104, msg) }
	ErrNotFound         = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusBadRequest, 105, msg) }
	ErrInvalidSignature = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusBadRequest, 107, msg) }

	ErrInternal = func(msg string) ErrorInfo {
		return NewErrorInfo(http.StatusInternalServerError, 100, msg)
	}
	ErrUnimplemented = func() ErrorInfo {
		return NewErrorInfo(http.StatusInternalServerError, 101, "unimplemented")
	}
)
