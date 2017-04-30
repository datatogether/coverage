package main

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"
)

var eotService = &Service{
	Id:          "3f5b22g5-37b4-5dc3-be91-c324bf686a87",
	Name:        "EOT Nomination Tool",
	Description: "",
	HomeUrl:     "https://github.com/edgi-govdata-archiving/eot-nomination-tool",
}

func addNominationUncrawlables(tree *Node) error {
	rawData, err := ioutil.ReadFile("sources/nomination_tool_epa_primer_uncrawlables.json")
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
			ServiceId:   eotService.Id,
			Uncrawlable: true,
		})
	}

	return nil
}
