package main

import (
	"github.com/archivers-space/coverage/repositories"
	"github.com/archivers-space/coverage/tree"
	"io"
	"net/http"
	"strconv"
)

func reqParamInt(key string, r *http.Request) (int, error) {
	i, err := strconv.ParseInt(r.FormValue(key), 10, 0)
	return int(i), err
}

func reqParamBool(key string, r *http.Request) (bool, error) {
	return strconv.ParseBool(r.FormValue(key))
}

func RootNodeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		EmptyOkHandler(w, r)
	case "GET":
		writeResponse(w, tree.CopyToDepth(t, 1))
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
			writeErrResponse(w, http.StatusNotFound, ErrNotFound)
			return
		}
		writeResponse(w, tree.CopyToDepth(node, 1))
	default:
		NotFoundHandler(w, r)
	}
}

func ListRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		EmptyOkHandler(w, r)
	case "GET":
		ListRepositories(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func ListRepositories(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, repositories.Repositories)
}

func RepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		EmptyOkHandler(w, r)
	case "GET":
		GetRepository(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func GetRepository(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/repositories/"):]

	for _, r := range repositories.Repositories {
		if r.GetId() == id {
			writeResponse(w, r)
			return
		}
	}
	NotFoundHandler(w, r)
}

func FullTreeHandler(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, t)
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
