package trace

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type KV struct {
	Key   string
	Value interface{}
}

// StartTransaction
func StartTransaction(ctx context.Context, c *gin.Context, service string) *sentry.Span {
	r := c.Request
	path := c.Request.URL.Path
	for _, param := range c.Params {
		path = strings.Replace(path, param.Value, ":"+param.Key, -1)
	}

	span := sentry.StartSpan(ctx, string(ServiceOp),
		sentry.TransactionName(fmt.Sprintf("%s %s", service, path)),
		sentry.ContinueFromRequest(r),
	)

	span.SetTag("stage", os.Getenv("STAGE"))
	span.SetTag("service", strings.ToLower(service))
	return span
}

// AddInvokerMetadata tracks user and project information if available
func AddInvokerMetadata(c *gin.Context, uid *string) {
	if hub := GetHubFromContext(c); hub != nil {
		scope := hub.Scope()
		if uid != nil {
			u := sentry.User{
				ID: *uid,
			}

			scope.SetUser(u)
		}

		//scope.SetTag("custom_tag", value)

	}
}

type Operation string

const (
	ServiceOp       Operation = "http.service"   // Measure the entire incoming http request (sentry transaction)
	IntraServiceOp  Operation = "http.intra"     // Measure an intra-service http request
	DBCreateOp      Operation = "db.create"      // Measure a DB write
	DBGetOp         Operation = "db.get"         // Measure a DB read
	DBQueryOp       Operation = "db.query"       // Measure a DB query
	DBTransactionOp Operation = "db.transaction" // Meaure a DB transaction
	CacheSet        Operation = "cache.set"      // Measure cache set
	CacheGet        Operation = "cache.get"      // Measure cache get
	CacheDel        Operation = "cache.del"      // Measure cache delete
)

func setDescription(description string) sentry.SpanOption {
	return func(s *sentry.Span) {
		s.Description = description
	}
}

func AddDBSpan(ctx context.Context, operation Operation, entityKind string) *sentry.Span {
	return AddSpan(ctx, operation, fmt.Sprintf("%v_%v", operation, strings.ToLower(entityKind)))
}

// AddSpan will add a span to an existing transaction.
// Must be called after StartTransaction.
// Example operation values: db.query, db.update, db.get, db.put, cache.get, cache.put, auth.verify
func AddSpan(ctx context.Context, operation Operation, name string) *sentry.Span {
	if c, ok := ctx.(*gin.Context); ok {
		transaction := sentry.TransactionFromContext(c.Request.Context())
		if transaction != nil {
			return sentry.StartSpan(transaction.Context(), string(operation), setDescription(name))
		}
		return sentry.StartSpan(c, string(operation), setDescription(name))
	}
	transaction := sentry.TransactionFromContext(ctx)
	if transaction != nil {
		return sentry.StartSpan(transaction.Context(), string(operation), setDescription(name))
	}
	return sentry.StartSpan(ctx, string(operation), setDescription(name))
}

func GetTransactionFromContext(ctx context.Context) *sentry.Span {
	if c, ok := ctx.(*gin.Context); ok {
		transaction := sentry.TransactionFromContext(c.Request.Context())
		if transaction != nil {
			return transaction
		}
	}
	return sentry.TransactionFromContext(ctx)
}

// SetSpanData safely sets span data, checking for nil map and duplicate keys
func SetSpanData(span *sentry.Span, key string, value interface{}) {
	if span.Data == nil {
		span.Data = make(map[string]interface{})
	}

	if _, ok := span.Data[key]; ok {
		log.Warnf("key %v already in span", key)
		return
	}

	span.Data[key] = value
}
