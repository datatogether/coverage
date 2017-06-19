package ia

import (
	"encoding/json"
	"fmt"
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/tree"
	"io/ioutil"
	"net/url"
	"strings"
)

var Repository = &repository{
	Id:          "5d5583g5-38a8-26d3-be70-c324bf686a87",
	Title:       "Internet Archive",
	Description: "the internet archive",
	Url:         "https://archive.org",
}

type repository archive.DataRepo

func (s *repository) GetId() string {
	return s.Id
}

func (r *repository) DataRepo() *archive.DataRepo {
	dr := archive.DataRepo(*r)
	return &dr
}

func (s *repository) AddCompletions(t *tree.Node) {
	t.Walk(func(n *tree.Node) {
		if n.Archived == false && n.Coverage != nil {
			for _, c := range n.Coverage {
				if c.RepositoryId == s.Id {
					n.Archived = c.Archived
				}
			}
		}
	})
}

func (s *repository) AddUrls(t *tree.Node, sources ...*archive.Source) error {
	rawData, err := ioutil.ReadFile("repositories/ia/ia_urls.json")
	if err != nil {
		return err
	}

	info := []struct {
		Url          string
		Available    bool
		WaybackUrl   string
		Timestamp    string
		Status       string
		NumSnapshots int
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

		// node = node.Child(u.Scheme).Child(u.Host)
		node = node.Child(fmt.Sprintf("%s://%s", u.Scheme, u.Host))
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

		node.Archived = node.Archived || i.Available
		node.Coverage = append(node.Coverage, &tree.Coverage{
			RepositoryId: s.Id,
			Archived:     i.Available,
			ArchiveUrl:   i.WaybackUrl,
		})
	}

	return nil
}

func (s *repository) AddCoverage(t *tree.Node) {
	// no-op for now
}
