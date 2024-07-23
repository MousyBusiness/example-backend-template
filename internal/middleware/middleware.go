package middleware

import (
	"bytes"
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/mousybusiness/example-backend-template/internal/errs"
	"github.com/mousybusiness/example-backend-template/internal/trace"
	. "github.com/mousybusiness/example-backend-template/pkg/model"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	apiKeyHeader              = "X-Api-Key"
	BearerAuthorizationHeader = "Authorization" // Bearer token, ID Token or Access Token
	UIDContextKey             = "uid"
	EmailContextKey           = "email"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Trace(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("stacktrace from panic: %v", string(debug.Stack()))
				errs.SentryError(errors.Errorf("panic in Trace middleware, %v", r))
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					})
			}
		}()

		path := c.Request.URL.Path
		for _, param := range c.Params {
			//if pkRegex.MatchString(param.Value) {
			//	path = strings.Replace(path, param.Value, ":pk", 1) // :pk includes :pid so needs first action to avoid a bad replace
			//} else {
			path = strings.Replace(path, param.Value, ":"+param.Key, 1)
			//}
		}

		c.Set("Log", log.WithField("path", path))
		// ignore healthcheck
		if strings.HasSuffix(path, "/healthcheck") {
			c.Next()
			return
		}

		c.Next()

		uid, _ := c.Value(UIDContextKey).(*string)
		trace.AddInvokerMetadata(c, uid)
	}
}

func Log(c context.Context) *log.Entry {
	l, ok := c.Value("Log").(*log.Entry)
	if !ok {
		return log.WithField("path", "unknown")
	}

	return l
}

// LogError middleware ensures any non-200 codes
// are logged for visibility on issues.
func LogError(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("stacktrace from panic: %v", string(debug.Stack()))
				errs.SentryError(errors.Errorf("panic in LogError middleware, %v", r))
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					})
			}
		}()

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		s := fmt.Sprintf("[%v] reponse, uri: %v, method: %v, code: %v, body: %v",
			service, c.Request.RequestURI, c.Request.Method, c.Writer.Status(), w.body.String())
		code := c.Writer.Status()

		// ignores 404s because constant bot traffic
		if code >= 400 {
			log.Errorf(s)
			if code != 404 && code != 401 {
				sentry.CaptureException(errors.New(s))
			}
		} else if code < 400 && code != 200 && code != 204 {
			log.Warnf(s)
			sentry.CaptureMessage(s)
		}
	}
}

func AuthZ(auth *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var t = strings.Replace(c.Request.Header.Get(BearerAuthorizationHeader), "Bearer ", "", 1)

		jwt, err := auth.VerifyIDToken(c, t)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if uid, ok := c.Params.Get("uid"); ok {
			if uid != jwt.UID {
				log.Errorf("UID mismatch, uid: %v, token: %v", uid, jwt.UID)
				c.Status(http.StatusForbidden)
				return
			}
			log.Debug("UID matches path")
		}

		c.Set(UIDContextKey, &jwt.UID)
		iarr, ok := jwt.Firebase.Identities["email"].([]interface{})
		if ok {
			if email, ok := iarr[0].(string); ok {
				c.Set(EmailContextKey, &email)
			}
		}

		c.Next()
	}
}

func AuthAPIKey(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get(apiKeyHeader)

		if key == "" || secret != key {
			log.Errorf("api key mismatch!")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				ErrorResponse{
					Error: http.StatusText(http.StatusUnauthorized),
				})
			return
		}

		c.Next()
	}
}
