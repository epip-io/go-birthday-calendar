// Package conf contains functions to configure logging, the database and router
package conf

import (
	"bufio"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/epip-io/go-birthday-calendar/pkg/models"
)

// ConfigureLogger configures and returns a logger
func ConfigureLogger(cfg *models.LoggerConfig) (*log.Entry, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	if cfg.File != "" && cfg.File != "-" {
		f, errOpen := os.OpenFile(cfg.File, os.O_RDWR|os.O_APPEND, 0660)
		if errOpen != nil {
			return nil, errOpen
		}

		log.SetOutput(bufio.NewWriter(f))
	}

	if cfg.Level != "" {
		lvl, err := log.ParseLevel(strings.ToUpper(cfg.Level))
		if err != nil {
			return nil, err
		}

		log.SetLevel(lvl)
	}

	log.SetFormatter(&log.JSONFormatter{})

	return log.StandardLogger().WithField("hostname", hostname), nil
}

// LoggingMiddleware configures and adds logging to a http.HandlerFunc
func LoggingMiddleware(inner http.Handler, name string, logger *log.Entry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		logger.WithFields(log.Fields{
			"Method":     r.Method,
			"RequestURI": r.RequestURI,
			"Route":      name,
			"Latancy":    time.Since(start),
		}).Info("handling request")
	})
}
