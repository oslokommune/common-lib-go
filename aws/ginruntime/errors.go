package ginruntime

import (
	"fmt"
	"net/http"

	"github.com/oslokommune/common-lib-go/httpcomm"
)

// Server response status code and associated error message
type ApiError struct {
	Status int
	Reason string
	Detail error
}

func (this *ApiError) Error() string {
	statusText := http.StatusText(this.Status)
	if this.Reason != "" {
		return fmt.Sprintf("%s: %s", statusText, this.Reason)
	}
	return statusText
}

// Converts `err` to an `*ApiError` if it's not nil.
func Normalize(err error) *ApiError {
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case *ApiError:
		return err
	// TODO: See if we can avoid having a depencency towards httpcomm in the future
	case *httpcomm.HTTPError:
		if err.StatusCode == http.StatusNotFound {
			return NotFound(err.Error())
		} else if err.StatusCode == http.StatusForbidden {
			return Forbidden(err.Error())
		} else {
			return FailedDependency(err.Error())
		}
	default:
		// Vi returnerer ikke feilmeldingen her for å unngå å potensielt sende sensitiv informasjon
		return InternalServerError("internal server error")
	}
}

func BadRequest(reason string) *ApiError {
	return &ApiError{Status: http.StatusBadRequest, Reason: reason}
}

func Unauthorized(reason string) *ApiError {
	return &ApiError{Status: http.StatusUnauthorized, Reason: reason}
}

func Forbidden(reason string) *ApiError {
	return &ApiError{Status: http.StatusForbidden, Reason: reason}
}

func NotFound(reason string) *ApiError {
	return &ApiError{Status: http.StatusNotFound, Reason: reason}
}

func UnprocessableEntity(reason string) *ApiError {
	return &ApiError{Status: http.StatusUnprocessableEntity, Reason: reason}
}

func FailedDependency(reason string) *ApiError {
	return &ApiError{Status: http.StatusFailedDependency, Reason: reason}
}

func InternalServerError(reason string) *ApiError {
	return &ApiError{Status: http.StatusInternalServerError, Reason: reason}
}

func NotImplemented(reason string) *ApiError {
	return &ApiError{Status: http.StatusNotImplemented, Reason: reason}
}
