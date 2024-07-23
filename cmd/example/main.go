package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/mousybusiness/example-backend-template/internal/middleware"
	"github.com/mousybusiness/example-backend-template/internal/trace"
	"github.com/mousybusiness/example-backend-template/pkg/secrets"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"net/http"
	"os"
	"time"
)

const (
	serverLabel         = "EXAMPLE"
	adminAPIKeySecretID = "admin-api-key"
)

var (
	userAPIKey  string
	adminAPIKey string
)

// getSecrets gets API keys from GCP Secret Manager
func getSecrets() error {
	ctx := context.Background()

	// load Push Service API key from GCP Secret manager
	var err error
	adminAPIKey, err = secrets.Vault.GetSecret(ctx, adminAPIKeySecretID)
	if err != nil {
		return err
	}

	return nil
}

// HealthcheckHandler handles requests to /healthcheck
// used by GCP Load Balancer and Managed Instance Groups
func HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Message": "Okay"})
	}
}

func setupRouter() *gin.Engine {
	// setup GIN
	r := gin.New()
	_ = r.SetTrustedProxies([]string{"130.211.0.0/22", "35.191.0.0/16"}) // https://cloud.google.com/load-balancing/docs/https
	r.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/healthcheck", "/"),
		gin.Recovery(),
	)
	r.Use(trace.New(serverLabel, trace.Options{}))
	r.Use(middleware.Trace(serverLabel))
	r.Use(middleware.LogError(serverLabel))

	ctx := context.Background()
	serviceAccount := os.Getenv("FIREBASE_CONFIG_FILE")
	log.Debugf("service account file: %v", serviceAccount)
	opt := option.WithCredentialsFile(serviceAccount)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing auth, creating auth client"))
	}

	// setup healthcheck endpoint for managed instance groups
	healthcheck := r.Group("/healthcheck")
	healthcheck.GET("", HealthCheckHandler())

	// use JWT auth middleware
	authed := r.Group("/user/:uid")
	authed.Use(middleware.AuthZ(authClient))
	authed.GET("", HealthCheckHandler())

	admin := r.Group("/admin")
	admin.Use(middleware.AuthAPIKey(userAPIKey))
	admin.GET("", HealthCheckHandler())
	return r
}

func main() {
	defer sentry.Flush(2 * time.Second)

	if err := getSecrets(); err != nil {
		log.Fatalln("couldn't get secrets from Secret Manager", err)
	}

	// start API
	log.Fatalln(setupRouter().Run(":80"))
}
