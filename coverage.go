package main

import (
	"github.com/archivers-space/coverage/services"
	"github.com/archivers-space/coverage/tree"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// our main t node
var t = &tree.Node{Id: "root", Name: "Coverage"}

func init() {
	if err := LoadCachedTree(t); err != nil {
		logger.Println(err.Error())
	}

	for _, s := range services.Services {

		if err := s.AddUrls(t); err != nil {
			logger.Println(s.Info()["Name"])
			logger.Println(err.Error())
		}

		s.AddCoverage(t)
	}

	t.Walk(func(n *tree.Node) {
		n.NumDescendants = -1
		n.NumDescendantsArchived = 0
		n.NumChildren = len(n.Children)
		n.Walk(func(d *tree.Node) {
			n.SortChildren()
			n.NumDescendants++
			if d.Archived {
				n.NumDescendantsArchived++
			}
			if d.Children == nil {
				n.NumLeaves++
				if d.Archived {
					n.NumLeavesArchived++
				}
			}
		})
	})

	if err := WriteTreeCache("cache.json", t); err != nil {
		fmt.Println(err.Error())
	}
}

func LoadCachedTree(n *tree.Node) error {
	cacheData, err := ioutil.ReadFile("cache.json")
	if err != nil {
		// not having cache data isn't an error
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	t := &tree.Node{}
	if err := json.Unmarshal(cacheData, t); err != nil {
		return err
	}

	logger.Println("successfully loaded cached tree")
	*n = *t
	return nil
}

func WriteTreeCache(filename string, n *tree.Node) error {
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, os.ModePerm)
}
