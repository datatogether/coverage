package ia

import (
	"encoding/json"
	"github.com/archivers-space/coverage/tree"
	"io/ioutil"
	"net/url"
	"strings"
)

var Service = &service{
	Id:          "5d5583g5-38a8-26d3-be70-c324bf686a87",
	Name:        "Internet Archive",
	Description: "the internet archive",
	HomeUrl:     "https://archive.org",
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

func (s *service) AddCompletions(t *tree.Node) {
	t.Walk(func(n *tree.Node) {
		if n.Archived == false && n.Coverage != nil {
			for _, c := range n.Coverage {
				if c.ServiceId == s.Id {
					n.Archived = c.Archived
				}
			}
		}
	})
}

func (s *service) AddUrls(t *tree.Node) error {
	rawData, err := ioutil.ReadFile("services/ia/ia_urls.json")
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
			if c.ServiceId == s.Id {
				continue
			}
		}

		node.Coverage = append(node.Coverage, &tree.Coverage{
			ServiceId:  s.Id,
			Archived:   i.Available,
			ArchiveUrl: i.WaybackUrl,
		})
	}

	return nil
}

func (s *service) AddCoverage(t *tree.Node) {
	// no-op for now
}
