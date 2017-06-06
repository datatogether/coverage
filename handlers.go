package main

import (
	"github.com/archivers-space/api/apiutil"
	"github.com/archivers-space/coverage/tree"
	"github.com/archivers-space/errors"
	"io"
	"net/http"
)

func RootNodeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		EmptyOkHandler(w, r)
	case "GET":
		apiutil.WriteResponse(w, tree.CopyToDepth(t, 1))
	default:
		NotFoundHandler(w, r)
	}
}

func NodeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		EmptyOkHandler(w, r)
	case "GET":
		id := r.URL.Path[len("/tree/"):]
		node := t.Find(id)
		if node == nil {
			apiutil.WriteErrResponse(w, http.StatusNotFound, errors.ErrNotFound)
			return
		}
		apiutil.WriteResponse(w, tree.CopyToDepth(node, 1))
	default:
		NotFoundHandler(w, r)
	}
}

func FullTreeHandler(w http.ResponseWriter, r *http.Request) {
	apiutil.WriteResponse(w, t)
}

// HealthCheckHandler is a basic "hey I'm fine" for load balancers & co
// TODO - add Database connection & proper configuration checks here for more accurate
// health reporting
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status" : 200 }`))
}

// EmptyOkHandler is an empty 200 response, often used
// for OPTIONS requests that responds with headers set in addCorsHeaders
func EmptyOkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// CertbotHandler pipes the certbot response for manual certificate generation
func CertbotHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, cfg.CertbotResponse)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{ "status" :  "not found" }`))
}
