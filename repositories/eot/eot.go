package eot

import (
	"encoding/json"
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/tree"
	"io/ioutil"
	"net/url"
	"strings"
)

var Repository = &repository{
	Id:          "3f5b22g5-37b4-5dc3-be91-c324bf686a87",
	Title:       "EOT Nomination Tool",
	Description: "",
	Url:         "https://github.com/edgi-govdata-archiving/eot-nomination-tool",
}

type repository archive.DataRepo

func (s *repository) GetId() string {
	return s.Id
}

func (r *repository) DataRepo() *archive.DataRepo {
	dr := archive.DataRepo(*r)
	return &dr
}

func (s *repository) AddUrls(t *tree.Node, src *archive.Source) error {
	rawData, err := ioutil.ReadFile("repositories/eot/nomination_tool_epa_primer_uncrawlables.json")
	if err != nil {
		return err
	}

	info := []struct {
		Url string
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
			if c.RepositoryId == s.Id {
				continue
			}
		}

		node.Coverage = append(node.Coverage, &tree.Coverage{
			RepositoryId: s.Id,
			Uncrawlable:  true,
		})
	}

	return nil
}

func (s *repository) AddCoverage(t *tree.Node) {
	// noop b/c eot doesn't provide coverage info
}
