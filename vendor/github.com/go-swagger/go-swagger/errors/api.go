package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Error represents a error interface all swagger framework errors implement
type Error interface {
	error
	Code() int32
}

type apiError struct {
	code    int32
	message string
}

func (a *apiError) Error() string {
	return a.message
}

func (a *apiError) Code() int32 {
	return a.code
}

// New creates a new API error with a code and a message
func New(code int32, message string, args ...interface{}) Error {
	if len(args) > 0 {
		return &apiError{code, fmt.Sprintf(message, args...)}
	}
	return &apiError{code, message}
}

// NotFound creates a new not found error
func NotFound(message string, args ...interface{}) Error {
	if message == "" {
		message = "Not found"
	}
	return New(http.StatusNotFound, fmt.Sprintf(message, args...))
}

// NotImplemented creates a new not implemented error
func NotImplemented(message string) Error {
	return New(http.StatusNotImplemented, message)
}

// MethodNotAllowedError represents an error for when the path matches but the method doesn't
type MethodNotAllowedError struct {
	code    int32
	Allowed []string
	message string
}

func (m *MethodNotAllowedError) Error() string {
	return m.message
}

// Code the error code
func (m *MethodNotAllowedError) Code() int32 {
	return m.code
}

func errorAsJSON(err Error) []byte {
	b, _ := json.Marshal(struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	}{err.Code(), err.Error()})
	return b
}

func flattenComposite(errs *CompositeError) *CompositeError {
	var res []error
	for _, er := range errs.Errors {
		switch e := er.(type) {
		case *CompositeError:
			flat := flattenComposite(e)
			res = append(res, flat.Errors...)
		default:
			res = append(res, e)
		}
	}
	return CompositeValidationError(res...)
}

// MethodNotAllowed creates a new method not allowed error
func MethodNotAllowed(requested string, allow []string) Error {
	msg := fmt.Sprintf("method %s is not allowed, but [%s] are", requested, strings.Join(allow, ","))
	return &MethodNotAllowedError{code: http.StatusMethodNotAllowed, Allowed: allow, message: msg}
}

// ServeError the error handler interface implemenation
func ServeError(rw http.ResponseWriter, r *http.Request, err error) {
	switch e := err.(type) {
	case *CompositeError:
		er := flattenComposite(e)
		ServeError(rw, r, er.Errors[0])
	case *MethodNotAllowedError:
		rw.Header().Add("Allow", strings.Join(err.(*MethodNotAllowedError).Allowed, ","))
		rw.WriteHeader(int(e.Code()))
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(e))
		}
	case Error:
		rw.WriteHeader(int(e.Code()))
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(e))
		}
	default:
		rw.WriteHeader(http.StatusInternalServerError)
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(New(http.StatusInternalServerError, err.Error())))
		}
	}

}