package main

import (
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/repositories"
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
		log.Info("error loading cached tree:", err.Error())
	}

	for _, s := range repositories.Repositories {

		if err := s.AddUrls(t, nil); err != nil {
			log.Info(s.Info()["Name"])
			log.Info(err.Error())
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

	log.Info("successfully loaded cached tree")
	*n = *t
	return nil
}

func WriteTreeCache(filename string, n *tree.Node) error {
	data, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, os.ModePerm)
}

func CoverageTree(src *archive.Source) (*tree.Node, error) {
	t := &tree.Node{
		Name: src.Title,
		Id:   src.Id,
	}

	for _, s := range repositories.Repositories {
		if err := s.AddUrls(t, src); err != nil {
			log.Info(s.Info()["Name"])
			log.Info(err.Error())
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

	return t, nil
}

func CoverageSummary(src *archive.Source) (map[string]interface{}, error) {
	t, err := CoverageTree(src)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"archived":    t.NumDescendantsArchived,
		"descendants": t.NumDescendants,
	}, nil
}
