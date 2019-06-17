/*
	Coverage is a service for mapping an archiving surface area, and tracking
	the amount of that surface area that any number of archives have covered
*/
package main

import (
	"database/sql"
	"fmt"
	"github.com/datatogether/core"
	"github.com/datatogether/coverage/tree"
	"github.com/datatogether/sql_datastore"
	"github.com/datatogether/sqlutil"
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

	// our main t node
	// TODO - remove this
	t = &tree.Node{}

	// application database connection
	appDB = &sql.DB{}

	// hoist standard datastore
	store = sql_datastore.DefaultStore
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
		// fatal if the server is missing a vital configuration detail
		log.Fatal(fmt.Errorf("server configuration error: %s", err.Error()))
	}

	go func() {
		log.Infoln("connecting to db")
		if err := sqlutil.ConnectToDb("postgres", cfg.PostgresDbUrl, appDB); err != nil {
			log.Infoln(err.Error())
			return
		}
		log.Infoln("connected to db")
		sql_datastore.SetDB(appDB)
		sql_datastore.Register(
			&core.Source{},
			&core.Primer{},
		)
		update(appDB)
	}()
	go listenRpc()

	s := &http.Server{}
	// connect mux to server
	s.Handler = NewServerRoutes()

	// fire it up!
	log.Infof("starting server on port %s in %s mode\n", cfg.Port, cfg.Mode)

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
