package main

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"
)

var archiversService = &Service{
	Id:          "4c0122g5-38a8-40b3-be91-c324bf686a87",
	Name:        "archivers.space",
	Description: "",
	HomeUrl:     "https://www.archivers.space",
}

func addArchiversSpaceUncrawlables(tree *Node) error {
	rawData, err := ioutil.ReadFile("sources/archivers.space_urls.json")
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

		for _, c := range node.Coverage {
			if c.ServiceId == archiversService.Id {
				continue
			}
		}

		node.Coverage = append(node.Coverage, &Coverage{
			ServiceId: archiversService.Id,
			Priority:  i.Priority,
			Archived:  i.HarvestUrl != "",
		})
	}

	return nil
}

func markArchiversCompletions(tree *Node) {
	tree.Walk(func(n *Node) {
		if n.Coverage != nil {
			for _, c := range n.Coverage {
				if c.ServiceId == archiversService.Id {
					n.Walk(func(an *Node) {
						for _, c := range an.Coverage {
							if c.ServiceId == archiversService.Id {
								return
							}
						}
						an.Coverage = append(an.Coverage, &Coverage{
							ServiceId: archiversService.Id,
							Archived:  true,
						})
					})
				}
			}
		}
	})
}
