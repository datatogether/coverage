package main

import (
	"bufio"
	"net/url"
	"os"
	"strings"
)

var svalbardService = &Service{
	Id:   "1a1112f4-38a8-26d3-be70-c324bf686a87",
	Name: "Project Svalbard",
}

func addSvalbardUncrawlables(tree *Node) error {
	f, err := os.Open("sources/svalbard_urls.txt")
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
	}
	return nil
}
