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
	// Root    url.Url
	// Depth   int
	Sources []*archive.Source
	Repos   []repositories.CoverageRepository
}

// NewCoverageGenerator creates a CoverageGenerator with the default
// properties
func NewCoverageGenerator(repoIds []string, patterns []string) *CoverageGenerator {
	var sources []*archive.Source
	if patterns != nil {
		sources := make([]*archive.Source, len(patterns))
		for i, pattern := range patterns {
			sources[i] = &archive.Source{
				Url: pattern,
			}
		}
	}

	repos := repositories.Repositories
	if repoIds != nil {
		r := []repositories.CoverageRepository{}
		for _, id := range repoIds {
			for _, repo := range repos {
				if repo.GetId() == id {
					r = append(r, repo)
				}
			}
		}
		repos = r
	}
	return &CoverageGenerator{
		Sources: sources,
		Repos:   repos,
	}
}

func (c CoverageGenerator) Tree() (*tree.Node, error) {
	t := &tree.Node{
		Name: "coverage",
		Id:   "root",
	}

	// if len(c.Sources) == 1 {
	// 	// TODO - should this be like this?
	// 	t.Name = c.Sources[0].Title
	// 	t.Id = c.Sources[0].Id
	// }

	for _, s := range c.Repos {
		if err := s.AddUrls(t, c.Sources...); err != nil {
			// log.Info(s.Info()["Name"])
			// log.Info(err.Error())
		}

		s.AddCoverage(t)
	}

	t.Walk(func(n *tree.Node) {
		// n.NumDescendants = -1
		// n.NumDescendantsArchived = 0
		// n.NumLeaves = 0
		n.NumChildren = len(n.Children)
		n.Walk(func(d *tree.Node) {
			n.SortChildren()
			// n.NumDescendants++
			// if d.Archived {
			// 	n.NumDescendantsArchived++
			// }
			if len(d.Children) == 0 {
				n.NumLeaves++
				// if d.Archived {
				// 	n.NumLeavesArchived++
				// }
				for _, c := range d.Coverage {
					if c.Archived {
						n.NumLeavesArchived++
						break
					}
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

func (c CoverageGenerator) Summary() (*Summary, error) {
	t, err := c.Tree()
	if err != nil {
		return nil, err
	}

	return &Summary{
		Archived:    t.NumDescendantsArchived,
		Descendants: t.NumDescendants,
	}, nil
}
