package logging

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	DefaultLogLevel = log.DebugLevel
)

type Redirect int

const (
	StdErr Redirect = iota
	StdOut
)

// ConfigureLogging sets up logging using logrus.
// Logrus adds log level and time to logs.
// By default logs will print to StdErr.
func ConfigureLogging() {
	var err error
	logLevel := DefaultLogLevel
	ll := viper.GetString("log_level")
	if ll != "" {
		logLevel, err = log.ParseLevel(strings.ToLower(ll))
		if err != nil {
			log.Fatal(err)
		}
	}
	log.SetLevel(logLevel)

	logJson := viper.GetBool("log_json")
	if logJson {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{DisableColors: true})
	}

	r := viper.GetInt("log_redirect")
	if r < 0 || r > 2 {
		log.Fatal("invalid redirect value")
	}
	logRedirect := Redirect(r)
	switch logRedirect {
	case StdErr:
		log.SetOutput(os.Stderr)
	case StdOut:
		log.SetOutput(os.Stdout)
	}
}
