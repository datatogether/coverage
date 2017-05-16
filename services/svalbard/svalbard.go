package svalbard

import (
	"bufio"
	"github.com/archivers-space/coverage/tree"
	"net/url"
	"os"
	"strings"
)

var Service = &service{
	Id:          "1a1112f4-38a8-26d3-be70-c324bf686a87",
	Name:        "Project Svalbard",
	Description: "",
	HomeUrl:     "",
}

type service struct {
	Id          string
	Name        string
	Description string
	HomeUrl     string
}

func (s *service) Info() map[string]interface{} {
	return map[string]interface{}{
		"Id":          s.Id,
		"Name":        s.Name,
		"Description": s.Description,
		"HomeUrl":     s.HomeUrl,
	}
}

func (s *service) AddUrls(t *tree.Node) error {
	f, err := os.Open("services/svalbard/svalbard_urls.txt")
	if err != nil {
		return err
	}
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		node := t
		u, err := url.Parse(sc.Text())
		if err != nil {
			return err
		}
		if u.Scheme == "" {
			u.Scheme = "http"
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
	}
	return nil
}

func (s *service) AddCoverage(t *tree.Node) {
	// no-op, need to add svalbard coverage work
}
