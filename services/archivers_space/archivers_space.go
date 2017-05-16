package archivers_space

import (
	"encoding/json"
	"github.com/archivers-space/coverage/tree"
	"io/ioutil"
	"net/url"
	"strings"
)

var Service = &service{
	Id:          "4c0122g5-38a8-40b3-be91-c324bf686a87",
	Name:        "archivers.space",
	Description: "",
	HomeUrl:     "https://www.archivers.space",
}

type service struct {
	Id          string
	Name        string
	Description string
	HomeUrl     string
}

func (s *service) Info() map[string]interface{} {
	return map[string]interface{}{
		"Id":         s.Id,
		"Name":       s.Name,
		"Decription": s.Description,
		"HomeUrl":    s.HomeUrl,
	}
}

func (s *service) AddUrls(t *tree.Node) error {
	rawData, err := ioutil.ReadFile("services/archivers_space/archivers.space_urls.json")
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
			ServiceId: s.Id,
			Priority:  i.Priority,
			Archived:  i.HarvestUrl != "",
		})
	}

	return nil
}

func (s *service) AddCoverage(t *tree.Node) {
	t.Walk(func(n *tree.Node) {
		if n.Coverage != nil {
			for _, c := range n.Coverage {
				if c.ServiceId == s.Id {
					n.Walk(func(an *tree.Node) {
						for _, c := range an.Coverage {
							if c.ServiceId == s.Id {
								return
							}
						}
						an.Coverage = append(an.Coverage, &tree.Coverage{
							ServiceId: s.Id,
							Archived:  true,
						})
					})
				}
			}
		}
	})
}
