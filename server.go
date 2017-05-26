/*
	Coverage is a service for mapping an archiving surface area, and tracking
	the amount of that surface area that any number of archives have covered
*/
package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	// cfg is the global configuration for the server. It's read in at startup from
	// the config.json file and enviornment variables, see config.go for more info.
	cfg *config

	// log output handled by logrus package
	log = logrus.New()
)

func init() {
	// configure logger
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
}

func main() {
	var err error
	cfg, err = initConfig(os.Getenv("GOLANG_ENV"))
	if err != nil {
		// panic if the server is missing a vital configuration detail
		panic(fmt.Errorf("server configuration error: %s", err.Error()))
	}

	s := &http.Server{}
	// connect mux to server
	s.Handler = NewServerRoutes()

	// fire it up!
	fmt.Println("starting server on port", cfg.Port)

	// start server wrapped in a log.Fatal b/c http.ListenAndServe will not
	// return unless there's an error
	log.Fatal(StartServer(cfg, s))
}

// NewServerRoutes returns a Muxer that has all API routes.
// This makes for easy testing using httptest package
func NewServerRoutes() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/.well-known/acme-challenge/", CertbotHandler)
	m.Handle("/", middleware(HealthCheckHandler))

	m.Handle("/services", middleware(ListServicesHandler))
	m.Handle("/services/", middleware(ServicesHandler))

	m.Handle("/fulltree", middleware(FullTreeHandler))

	// m.Handle("/coverage", middleware(CoverageHandler))
	// m.Handle("/coverage/", middleware(CoverageHandler))

	m.Handle("/tree", middleware(RootNodeHandler))
	m.Handle("/tree/", middleware(NodeHandler))

	return m
}
