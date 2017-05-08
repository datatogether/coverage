package main

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"
)

var iaService = &Service{
	Id:          "5d5583g5-38a8-26d3-be70-c324bf686a87",
	Name:        "Internet Archive",
	Description: "the internet archive",
	HomeUrl:     "https://archive.org",
}

func markIACompletions(tree *Node) {
	tree.Walk(func(n *Node) {
		if n.Archived == false && n.Coverage != nil {
			for _, c := range n.Coverage {
				if c.ServiceId == iaService.Id {
					n.Archived = c.Archived
				}
			}
		}
	})
}

func addIAUrls(tree *Node) error {
	rawData, err := ioutil.ReadFile("sources/ia_urls.json")
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
		node := tree
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

		node.Coverage = append(node.Coverage, &Coverage{
			ServiceId: iaService.Id,
			Archived:  i.Available,
		})
	}

	return nil
}
