package archivers2

import (
	"bufio"
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/tree"
	"net/url"
	"os"
	"strings"
)

// Concrete Archivers2 instance
var Repository = &repository{
	Id:          "8d7e22g5-38a8-40b3-be91-c324bf686a87",
	Title:       "archivers 2.0",
	Description: "",
	Url:         "https://alpha.archivers.space",
}

type repository archive.DataRepo

func (s *repository) GetId() string {
	return s.Id
}

func (r *repository) Info() map[string]interface{} {
	return map[string]interface{}{
		"Id":          r.Id,
		"Title":       r.Title,
		"Description": r.Description,
		"Url":         r.Url,
	}
}

func (a *repository) AddCoverage(t *tree.Node) {
	t.Walk(func(n *tree.Node) {
		if n.Archived == false && n.Coverage != nil {
			for _, c := range n.Coverage {
				if c.RepositoryId == a.Id {
					n.Archived = c.Archived
					if c.Archived {
						n.ArchiveCount++
					}
					break
				}
			}
		}
	})
}

func (a *repository) AddUrls(t *tree.Node, src *archive.Source) error {
	f, err := os.Open("repositories/archivers2/archivers_2_downloaded_epa_content_urls.txt")
	if err != nil {
		return err
	}
	s := bufio.NewScanner(f)

	for s.Scan() {
		node := t
		u, err := url.Parse(s.Text())
		if err != nil {
			return err
		}
		if u.Scheme == "" {
			u.Scheme = "http"
		}

		// skip this url if it doesn't match the passed in Source filter
		if src != nil && !src.MatchesUrl(u.String()) {
			continue
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
			if c.RepositoryId == a.Id {
				continue
			}
		}

		node.Coverage = append(node.Coverage, &tree.Coverage{
			// Url:       u.String(),
			RepositoryId: a.Id,
			Archived:     true,
		})
	}
	return nil
}
