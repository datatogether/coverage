package archivers2

import (
	"bufio"
	"github.com/archivers-space/coverage/tree"
	"net/url"
	"os"
	"strings"
)

// Concrete Archivers2 instance
var Service = &service{
	Id:          "8d7e22g5-38a8-40b3-be91-c324bf686a87",
	Name:        "archivers 2.0",
	Description: "",
	HomeUrl:     "https://alpha.archivers.space",
}

type service struct {
	Id          string
	Name        string
	Description string
	HomeUrl     string
}

func (a *service) Info() map[string]interface{} {
	return map[string]interface{}{
		"Id":          a.Id,
		"Name":        a.Name,
		"Description": a.Description,
		"HomeUrl":     a.HomeUrl,
	}
}

func (a *service) AddCoverage(t *tree.Node) {
	t.Walk(func(n *tree.Node) {
		if n.Archived == false && n.Coverage != nil {
			for _, c := range n.Coverage {
				if c.ServiceId == a.Id {
					n.Archived = c.Archived
				}
			}
		}
	})
}

func (a *service) AddUrls(t *tree.Node) error {
	f, err := os.Open("services/archivers2/archivers_2_downloaded_epa_content_urls.txt")
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
			if c.ServiceId == a.Id {
				continue
			}
		}

		node.Coverage = append(node.Coverage, &tree.Coverage{
			// Url:       u.String(),
			ServiceId: a.Id,
			Archived:  true,
		})
	}
	return nil
}
