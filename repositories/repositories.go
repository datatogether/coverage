package repositories

import (
	"github.com/datatogether/core"
	"github.com/datatogether/coverage/repositories/archivers2"
	"github.com/datatogether/coverage/repositories/archivers_space"
	"github.com/datatogether/coverage/repositories/eot"
	"github.com/datatogether/coverage/repositories/ia"
	"github.com/datatogether/coverage/repositories/svalbard"
	"github.com/datatogether/coverage/tree"
)

// a Repository is anything that stores data.
type Repository interface {
	GetId() string
	DataRepo() *core.DataRepo
}

// A Coverage Repository is any service that can also provide coverage information
type CoverageRepository interface {
	Repository
	AddUrls(tree *tree.Node, sources ...*core.Source) error
	AddCoverage(tree *tree.Node)
	// UrlCoverage(rawurl string) (*tree.Coverage, error)
}

var Repositories = []CoverageRepository{
	archivers2.Repository,
	archivers_space.Repository,
	eot.Repository,
	ia.Repository,
	svalbard.Repository,
}
