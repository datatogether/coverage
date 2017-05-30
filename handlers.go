package main

import (
	"github.com/archivers-space/coverage/services"
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

func CoverageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		EmptyOkHandler(w, r)
	case "GET":
		ReadCoverage(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func ReadCoverage(w http.ResponseWriter, r *http.Request) {
	// u := &archive.Url{Id: r.URL.Path[len("/urls/"):]}
	// if err := u.Read(appDB); err != nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	io.WriteString(w, err.Error())
	// 	return
	// }

	// f, err := u.File()
	// if err != nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	io.WriteString(w, err.Error())
	// 	return
	// }

	// if err := f.GetS3(); err != nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	io.WriteString(w, err.Error())
	// 	return
	// }

	// w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, u.FileName))
	// w.Write(f.Data)
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

func ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		EmptyOkHandler(w, r)
	case "GET":
		ListServices(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func ListServices(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, services.Services)
}

func ServicesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		EmptyOkHandler(w, r)
	case "GET":
		GetService(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func GetService(w http.ResponseWriter, r *http.Request) {
	// TODO
	// s := &Service{
	// 	Id: r.URL.Path[len("/services/"):],
	// }

	// if err := s.Read(); err != nil {
	// 	writeErrResponse(w, http.StatusNotFound, err)
	// 	return
	// }

	// writeResponse(w, s)
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
