package main

import (
	"github.com/archivers-space/api/apiutil"
	"github.com/archivers-space/coverage/repositories"
	"net/http"
)

// Concrete CoverateRequests instance
var RepositoryRequests = new(repositories.RepositoryRequests)

// func CoverageSummaryHandler(w http.ResponseWriter, r *http.Request) {
//   args := &coverage.CoverageSummaryParams{
//     Pattern: r.FormValue("pattern"),
//   }

//   res := &coverage.Summary{}
//   err := CoverageRequests.Summary(args, res)
//   if err != nil {
//     log.Info(err.Error())
//     apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
//     return
//   }
//   apiutil.WriteResponse(w, res)
// }

// func GetCoverageTreeHandler(w http.ResponseWriter, r *http.Request) {
//   args := &coverage.CoverageTreeParams{
//     Pattern: r.FormValue("pattern"),
//   }
//   res := &tree.Node{}
//   err := CoverageRequests.Tree(args, res)
//   if err != nil {
//     log.Info(err.Error())
//     apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
//     return
//   }
//   apiutil.WriteResponse(w, res)
// }

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
	apiutil.WriteResponse(w, repositories.Repositories)
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
			apiutil.WriteResponse(w, r)
			return
		}
	}
	NotFoundHandler(w, r)
}
