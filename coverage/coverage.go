package coverage

import (
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/repositories"
	"github.com/archivers-space/coverage/tree"

	"encoding/json"
	"io/ioutil"
	"os"
)

// our main t node
var t = &tree.Node{Id: "root", Name: "Coverage"}

func InitTree(t *tree.Node) error {
	if err := LoadCachedTree(t); err != nil {
		// log.Info("error loading cached tree:", err.Error())
		return err
	}

	for _, s := range repositories.Repositories {

		if err := s.AddUrls(t, nil); err != nil {
			// log.Info(s.Info()["Name"])
			// log.Info(err.Error())
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
		return err
	}
	return nil
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

	// log.Info("successfully loaded cached tree")
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

// CovreageGen holds configuration for coverage analysis
type CoverageGenerator struct {
}

// NewCoverageGenerator creates a CoverageGenerator with the default
// properties
func NewCoverageGenerator() *CoverageGenerator {
	return &CoverageGenerator{}
}

func (c CoverageGenerator) Tree(src *archive.Source) (*tree.Node, error) {
	t := &tree.Node{
		Name: src.Title,
		Id:   src.Id,
	}

	for _, s := range repositories.Repositories {
		if err := s.AddUrls(t, src); err != nil {
			// log.Info(s.Info()["Name"])
			// log.Info(err.Error())
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

type Summary struct {
	Archived    int
	Descendants int
}

func (c CoverageGenerator) Summary(src *archive.Source) (*Summary, error) {
	t, err := c.Tree(src)
	if err != nil {
		return nil, err
	}

	return &Summary{
		Archived:    t.NumDescendantsArchived,
		Descendants: t.NumDescendants,
	}, nil
}
