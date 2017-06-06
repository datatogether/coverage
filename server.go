/*
	Coverage is a service for mapping an archiving surface area, and tracking
	the amount of that surface area that any number of archives have covered
*/
package main

import (
	"database/sql"
	"fmt"
	"github.com/archivers-space/coverage/tree"
	"github.com/archivers-space/sqlutil"
	_ "github.com/lib/pq"
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

	// application database connection
	appDB = &sql.DB{}

	// our main t node
	t = &tree.Node{}
)

func init() {
	// configure logger
	log.Out = os.Stderr
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

	go func() {
		if err := sqlutil.ConnectToDb("postgres", cfg.PostgresDbUrl, appDB); err != nil {
			log.Infoln(err.Error())
			return
		}
		update(appDB)
	}()
	go listenRpc()

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
	m.Handle("/healthcheck", middleware(HealthCheckHandler))
	m.Handle("/", middleware(NotFoundHandler))

	m.Handle("/repositories", middleware(ListRepositoriesHandler))
	m.Handle("/repositories/", middleware(RepositoriesHandler))

	m.Handle("/fulltree", middleware(FullTreeHandler))

	m.Handle("/coverage", middleware(CoverageHandler))
	m.Handle("/tree", middleware(CoverageTreeHandler))

	return m
}
