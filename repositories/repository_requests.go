package repositories

import (
	"fmt"
	"github.com/datatogether/core"
)

type RepositoryRequests int

type RepositoryListParams struct {
	Orderby string
	Limit   int
	Offset  int
}

func (r RepositoryRequests) List(p *RepositoryListParams, res *[]*core.DataRepo) error {
	repos := make([]*core.DataRepo, len(Repositories))
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

func (r RepositoryRequests) Get(p *RepositoryGetParams, res *core.DataRepo) error {
	for _, repo := range Repositories {
		if repo.GetId() == p.Id {
			*res = *repo.DataRepo()
			return nil
		}
	}
	return fmt.Errorf("Not Found")
}
