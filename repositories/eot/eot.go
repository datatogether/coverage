package eot

import (
	"encoding/json"
	"fmt"
	"github.com/datatogether/core"
	"github.com/datatogether/coverage/tree"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var Repository = &repository{
	Id:          "3f5b22g5-37b4-5dc3-be91-c324bf686a87",
	Title:       "EOT Nomination Tool",
	Description: "",
	Url:         "https://github.com/edgi-govdata-archiving/eot-nomination-tool",
}

type repository core.DataRepo

func (s *repository) GetId() string {
	return s.Id
}

func (r *repository) DataRepo() *core.DataRepo {
	dr := core.DataRepo(*r)
	return &dr
}

func (s *repository) AddUrls(t *tree.Node, sources ...*core.Source) error {
	rawData, err := ioutil.ReadFile(filepath.Join(os.Getenv("GOPATH"), "src/github.com/datatogether/coverage", "repositories/eot/nomination_tool_epa_primer_uncrawlables.json"))
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
