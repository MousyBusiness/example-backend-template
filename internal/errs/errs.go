package errs

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/mousybusiness/example-backend-template/pkg/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

// HttpError defines an error which contains
// an http status code from an API request.
// This faciliates a more accurate API response
// for the daemon API.
type HttpError interface {
	error
	Code() int
	Cause() string
}

type httpError struct {
	code int
	err  error
}

func NewHttpError(code int, b []byte, fallback string) httpError {
	e := errors.New(fallback)

	if b != nil {
		var r model.ErrorResponse
		if err := json.Unmarshal(b, &r); err == nil {
			e = errors.New(r.Error)
		}
	}

	return httpError{
		code: code,
		err:  e,
	}
}

func (e httpError) Error() string {
	return errors.Wrap(e.err, fmt.Sprintf("HttpError[%v]", e.code)).Error()
}

func (e httpError) Cause() string {
	return e.err.Error()
}

func (e httpError) Code() int {
	return e.code
}

func ExtractHttpError(err error) (HttpError, bool) {
	e, ok := err.(httpError)
	if !ok {
		return nil, false
	}
	return e, true
}

func SentryWarnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	sentry.CaptureMessage(msg)
	log.Warn(msg)
}

func SentryError(err error) {
	sentry.CaptureException(err)
	log.Error(err)
}

func SentryFatal(err error) {
	sentry.CaptureException(err)
	sentry.Flush(2 * time.Second)
	log.Fatal(err)
}
