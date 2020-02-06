// Package conf contains functions to configure logging, the database and router
package conf

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/epip-io/go-birthday-calendar/pkg/models"
)

// Route defines a single managed route
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

// Routes is the list of managed routes
type Routes []Route

// ConfigureRouter configures and returns a http server
func ConfigureRouter(cfg *models.Config, db *gorm.DB, logger *log.Entry) *mux.Router {
	routes := Routes{
		Route{
			Name:        "Health",
			Method:      "GET",
			Path:        "/healthz",
			HandlerFunc: Health(db, logger),
		},
		Route{
			Name:        "BirthdayService",
			Method:      "GET",
			Path:        "/",
			HandlerFunc: BirthdayService(),
		},
		Route{
			Name:        "BirthdayMessage",
			Method:      "GET",
			Path:        "/{Name}",
			HandlerFunc: BirthdayMessage(db, logger),
		},
		Route{
			Name:        "PersonsBirthday",
			Method:      "PUT",
			Path:        "/{Name}",
			HandlerFunc: PersonsBirthday(db, logger),
		},
	}

	router := mux.NewRouter().StrictSlash(true)

	for _, r := range routes {
		var handler http.Handler
		var scheme string

		handler = r.HandlerFunc
		// handler = LoggingMiddleware(handler, r.Name, logger)
		scheme = "http"
		if cfg.TLS.Enabled {
			scheme = "https"
		}

		logger.WithFields(log.Fields{
			"method":     r.Method,
			"pathPrefix": cfg.Path,
			"path":       r.Path,
			"name":       r.Name,
			"schemes":    scheme,
		}).Debug("adding route")
		router.
			PathPrefix(cfg.Path).
			Path(r.Path).
			Methods(r.Method).
			Name(r.Name).
			Handler(handler)
	}

	if cfg.Redirect && cfg.TLS.Enabled {
		router.Use(RedirectMiddleware)
	}

	return router
}

// RedirectMiddleware adds a HTTP to HTTPS redirect http.HandlerFunc
func RedirectMiddleware(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.Header.Get("x-forward-proto")

		if strings.ToLower(p) == "http" {
			http.Redirect(w, r, fmt.Sprintf("https://%s%s", r.Host, r.URL), http.StatusPermanentRedirect)
			return
		}

		inner.ServeHTTP(w, r)
	})
}
