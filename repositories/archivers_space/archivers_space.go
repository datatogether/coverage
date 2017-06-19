package archivers_space

import (
	"encoding/json"
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/tree"
	"io/ioutil"
	"net/url"
	"strings"
)

var Repository = &repository{
	Id:          "4c0122g5-38a8-40b3-be91-c324bf686a87",
	Title:       "archivers.space",
	Description: "",
	Url:         "https://www.archivers.space",
}

type repository archive.DataRepo

func (s *repository) GetId() string {
	return s.Id
}

func (r *repository) DataRepo() *archive.DataRepo {
	dr := archive.DataRepo(*r)
	return &dr
}

func (s *repository) AddUrls(t *tree.Node, sources ...*archive.Source) error {
	rawData, err := ioutil.ReadFile("repositories/archivers_space/archivers.space_urls.json")
	if err != nil {
		return err
	}

	info := []struct {
		HarvestUrl string `json:"harvest_url"`
		Url        string
		Priority   int
		BagDone    bool   `json:"bag_done"`
		BagUrl     string `json:"bag_url"`
	}{}

	if err := json.Unmarshal(rawData, &info); err != nil {
		return err
	}

	for _, i := range info {
		node := t
		u, err := url.Parse(i.Url)
		if err != nil {
			return err
		}
		if u.Scheme == "" {
			u.Scheme = "http"
		}

		// skip this url if it doesn't match the passed in Source filter
		if len(sources) > 0 {
			match := false
			for _, src := range sources {
				if src != nil && src.MatchesUrl(u.String()) {
					match = true
				}
			}
			if !match {
				continue
			}
		}

		node = node.Child(u.Scheme).Child(u.Host)
		components := strings.Split(u.Path, "/")

		for _, c := range components {
			if c != "" {
				node = node.Child(c)
			}
		}
		if u.RawQuery != "" {
			node = node.Child(u.RawQuery)
		}

		for _, c := range node.Coverage {
			if c.RepositoryId == s.Id {
				continue
			}
		}

		node.Coverage = append(node.Coverage, &tree.Coverage{
			RepositoryId: s.Id,
			Priority:     i.Priority,
			Archived:     i.HarvestUrl != "",
		})
	}

	return nil
}

func (s *repository) AddCoverage(t *tree.Node) {
	t.Walk(func(n *tree.Node) {
		if n.Coverage != nil {
			for _, c := range n.Coverage {
				if c.RepositoryId == s.Id {
					n.Walk(func(an *tree.Node) {
						for _, c := range an.Coverage {
							if c.RepositoryId == s.Id {
								return
							}
						}
						an.Coverage = append(an.Coverage, &tree.Coverage{
							RepositoryId: s.Id,
							Archived:     true,
						})
					})
				}
			}
		}
	})
}
