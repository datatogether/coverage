package repositories

import (
	"github.com/archivers-space/archive"
	"github.com/archivers-space/errors"
)

type RepositoryRequests int

type RepositoryListParams struct {
	Orderby string
	Limit   int
	Offset  int
}

func (r RepositoryRequests) List(p *RepositoryListParams, res *[]*archive.DataRepo) error {
	repos := make([]*archive.DataRepo, len(Repositories))
	for i, repo := range Repositories {
		repos[i] = repo.DataRepo()
	}
	*res = repos
	return nil
}

type RepositoryGetParams struct {
	Id   string
	Name string
}

func (r RepositoryRequests) Get(p *RepositoryGetParams, res *archive.DataRepo) error {
	for _, repo := range Repositories {
		if repo.GetId() == p.Id {
			*res = *repo.DataRepo()
			return nil
		}
	}
	return errors.ErrNotFound
}
