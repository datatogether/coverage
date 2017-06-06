package repositories

import (
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/repositories/archivers2"
	"github.com/archivers-space/coverage/repositories/archivers_space"
	"github.com/archivers-space/coverage/repositories/eot"
	"github.com/archivers-space/coverage/repositories/ia"
	"github.com/archivers-space/coverage/repositories/svalbard"
	"github.com/archivers-space/coverage/tree"
)

// a Repository is anything that stores data.
type Repository interface {
	GetId() string
	DataRepo() *archive.DataRepo
}

// A Coverage Repository is any service that can also provide coverage information
type CoverageRepository interface {
	Repository
	AddUrls(tree *tree.Node, src *archive.Source) error
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
