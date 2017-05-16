package eot

import (
	"encoding/json"
	"github.com/archivers-space/coverage/tree"
	"io/ioutil"
	"net/url"
	"strings"
)

var Service = &service{
	Id:          "3f5b22g5-37b4-5dc3-be91-c324bf686a87",
	Name:        "EOT Nomination Tool",
	Description: "",
	HomeUrl:     "https://github.com/edgi-govdata-archiving/eot-nomination-tool",
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
	rawData, err := ioutil.ReadFile("services/eot/nomination_tool_epa_primer_uncrawlables.json")
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
			ServiceId:   s.Id,
			Uncrawlable: true,
		})
	}

	return nil
}

func (s *service) AddCoverage(t *tree.Node) {
	// noop b/c eot doesn't provide coverage info
}
