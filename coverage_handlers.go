package main

import (
	"github.com/archivers-space/coverage/tree"
	"net/http"
)

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
	args := &CoverageSummaryArgs{
		Pattern: r.FormValue("pattern"),
	}

	res := &CoverageSummary{}
	err := new(Coverage).Summary(args, res)
	if err != nil {
		log.Info(err.Error())
		writeErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeResponse(w, res)
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
	args := &CoverageTreeArgs{
		Pattern: r.FormValue("pattern"),
	}
	res := &tree.Node{}
	err := new(Coverage).Tree(args, res)
	if err != nil {
		log.Info(err.Error())
		writeErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeResponse(w, res)
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
