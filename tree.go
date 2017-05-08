package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// our main tree node
var tree = &Node{Id: "root", Name: "Coverage"}

func init() {
	if err := LoadCachedTree(tree); err != nil {
		logger.Println(err.Error())
	}

	if err := addSvalbardUncrawlables(tree); err != nil {
		panic(err)
	}

	if err := addArchivers2Gotten(tree); err != nil {
		panic(err)
	}

	if err := addArchiversSpaceUncrawlables(tree); err != nil {
		panic(err)
	}

	if err := addNominationUncrawlables(tree); err != nil {
		panic(err)
	}

	if err := addIAUrls(tree); err != nil {
		panic(err)
	}

	markArchiversCompletions(tree)
	markArchivers2Completions(tree)
	markIACompletions(tree)

	tree.Walk(func(n *Node) {
		n.NumDescendants = -1
		n.NumDescendantsArchived = 0
		n.NumChildren = len(n.Children)
		n.Walk(func(d *Node) {
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

	if err := WriteTreeCache("cache.json", tree); err != nil {
		fmt.Println(err.Error())
	}
}

func LoadCachedTree(n *Node) error {
	cacheData, err := ioutil.ReadFile("cache.json")
	if err != nil {
		// not having cache data isn't an error
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	tree := &Node{}
	if err := json.Unmarshal(cacheData, tree); err != nil {
		return err
	}

	logger.Println("successfully loaded cached tree")
	*n = *tree
	return nil
}

func WriteTreeCache(filename string, n *Node) error {
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, os.ModePerm)
}
