package main

import (
	"bufio"
	"net/url"
	"os"
	"strings"
)

var archivers2Service = &Service{
	Id:          "8d7e22g5-38a8-40b3-be91-c324bf686a87",
	Name:        "archivers 2.0",
	Description: "",
	HomeUrl:     "https://alpha.archivers.space",
}

func markArchivers2Completions(tree *Node) {
	tree.Walk(func(n *Node) {
		if n.Archived == false && n.Coverage != nil {
			for _, c := range n.Coverage {
				if c.ServiceId == archivers2Service.Id {
					n.Archived = c.Archived
				}
			}
		}
	})
}

func addArchivers2Gotten(tree *Node) error {
	f, err := os.Open("sources/archivers_2_downloaded_epa_content_urls.txt")
	if err != nil {
		return err
	}
	s := bufio.NewScanner(f)

	for s.Scan() {
		node := tree
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

		node.Coverage = append(node.Coverage, &Coverage{
			// Url:       u.String(),
			ServiceId: archivers2Service.Id,
			Archived:  true,
		})
	}
	return nil
}
