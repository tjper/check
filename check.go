package check

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// Check is a data object for checking and reporting errors.
type Check struct {
	// Message is an error message to be written by the check. This message
	// is intended to be read during debugging.
	Message string

	// ClientMessage is an error message written by a server check. This
	// message is intended to be read by a client application.
	ClientMessage string
}

// NewCheck creates a new Check object checks and reports on errors.
func NewCheck(options ...func(*Check)) *Check {
	chk := &Check{
		Message: "encountered error",
	}
	for _, option := range options {
		option(chk)
	}
	return chk
}

// WithMessage returns a Check option func that sets the Check.Message field.
func WithMessage(msg string) func(*Check) {
	return func(chk *Check) {
		chk.Message = msg
	}
}

// WithClientMessage returns a Check option func that sets the
// Check.ClientMessage field.
func WithClientMessage(msg string) func(*Check) {
	return func(chk *Check) {
		chk.ClientMessage = msg
	}
}

// Err checks an error. If the error exists, the error is logged along with the
// Check objects Message.
func (chk *Check) Err(err error) bool {
	if err != nil {
		glog.Error(errors.Wrap(err, chk.Message))
		return true
	}
	return false
}

// SrvErr checks an error. If the error exists, the error is logged with the
// Check object's Message field, and the Check object's ClientMessage field is
// written as an internal server error to the client.
func (chk *Check) SrvErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		glog.Error(errors.Wrap(err, chk.Message))
		http.Error(w, chk.ClientMessage, http.StatusInternalServerError)
		return true
	}
	return false
}
