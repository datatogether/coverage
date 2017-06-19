package main

import (
	"github.com/archivers-space/api/apiutil"
	"github.com/archivers-space/coverage/coverage"
	"github.com/archivers-space/coverage/tree"
	"net/http"
	"strconv"
	"strings"
)

// Concrete CoverateRequests instance
var CoverageRequests = new(coverage.CoverageRequests)

func CoverageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		EmptyOkHandler(w, r)
	case "GET":
		CoverageSummaryHandler(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func CoverageSummaryHandler(w http.ResponseWriter, r *http.Request) {
	args := &coverage.CoverageSummaryParams{
		Patterns: strings.Split(r.FormValue("patterns"), ","),
		RepoIds:  strings.Split(r.FormValue("repos"), ","),
	}

	res := &coverage.Summary{}
	err := CoverageRequests.Summary(args, res)
	if err != nil {
		log.Info(err.Error())
		apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	apiutil.WriteResponse(w, res)
}

func CoverageTreeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		EmptyOkHandler(w, r)
	case "GET":
		GetCoverageTreeHandler(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func GetCoverageTreeHandler(w http.ResponseWriter, r *http.Request) {
	depth, err := strconv.ParseInt(r.FormValue("depth"), 10, 0)
	if err != nil {
		depth = 0
	}

	args := &coverage.CoverageTreeParams{
		Root:     r.FormValue("root"),
		Depth:    int(depth),
		Patterns: strings.Split(r.FormValue("patterns"), ","),
		RepoIds:  strings.Split(r.FormValue("repos"), ","),
	}
	res := &tree.Node{}
	err = CoverageRequests.Tree(args, res)
	if err != nil {
		log.Info(err.Error())
		apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	apiutil.WriteResponse(w, res)
}

func DownloadCoverageTreeHandler(w http.ResponseWriter, r *http.Request) {
	// u := &archive.Url{Id: r.URL.Path[len("/urls/"):]}
	// if err := u.Read(appDB); err != nil {
	//  w.WriteHeader(http.StatusNotFound)
	//  io.WriteString(w, err.Error())
	//  return
	// }

	// f, err := u.File()
	// if err != nil {
	//  w.WriteHeader(http.StatusNotFound)
	//  io.WriteString(w, err.Error())
	//  return
	// }

	// if err := f.GetS3(); err != nil {
	//  w.WriteHeader(http.StatusNotFound)
	//  io.WriteString(w, err.Error())
	//  return
	// }

	// w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, u.FileName))
	// w.Write(f.Data)
}
