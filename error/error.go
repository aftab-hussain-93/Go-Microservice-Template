// Custom error package
package err

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errCode string

const (
	ErrInvalid         errCode = "INVALID_INPUT"
	ErrUnAuthenticated errCode = "NOT_FOUND"
	ErrForbidden       errCode = "FORBIDDEN"
	ErrInternal        errCode = "INTERNAL"
)

var errToHTTPStatus = map[errCode]int{
	ErrInvalid:         http.StatusBadRequest,
	ErrUnAuthenticated: http.StatusUnauthorized,
	ErrForbidden:       http.StatusForbidden,
	ErrInternal:        http.StatusInternalServerError,
}

type Err struct {
	// inner error, optional
	Err error
	// error code explaining the error
	Code       errCode
	Message    string
	Resolution string
}

func (e *Err) Error() (errStr string) {
	errStr = fmt.Sprintf("code: %s|message: %s|resolution: %s", e.Code, e.Message, e.Resolution)
	if e.Err != nil {
		errStr += "inner: " + e.Err.Error()
	}
	return
}

func (e *Err) GetHTTPStatusCode() (s int) {
	s = errToHTTPStatus[e.Code]
	if s == 0 {
		s = http.StatusInternalServerError
	}
	return
}

func (e *Err) MarshalJSON() ([]byte, error) {
	r := map[string]any{
		"code":       e.Code,
		"message":    e.Message,
		"resolution": e.Resolution,
	}
	if e.Err != nil {
		r["error"] = e.Err.Error()
	}

	return json.Marshal(r)
}
