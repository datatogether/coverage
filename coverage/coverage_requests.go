package coverage

import (
	"fmt"
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/tree"
	"net/url"
	"strings"
)

type CoverageRequests int

type CoverageTreeParams struct {
	// root url to work from
	Root string
	// depth of tree to return
	Depth int
	// patterns to filter results against, optional
	Patterns []string
	// ids of repositories to limit query to, default is all
	RepoIds []string
}

func (p *CoverageTreeParams) Validate() error {
	return nil
}

func (CoverageRequests) Tree(p *CoverageTreeParams, res *tree.Node) error {
	root, err := NewCoverageGenerator(p.RepoIds, p.Patterns).Tree()
	if err != nil {
		return err
	}

	if p.Root != "" {
		u, err := url.Parse(p.Root)
		if err != nil {
			return err
		}
		if u.Scheme == "" {
			u.Scheme = "http"
		}

		root = root.Child(fmt.Sprintf("%s://%s", u.Scheme, u.Host))
		components := strings.Split(u.Path, "/")
		for _, c := range components {
			if c != "" {
				root = root.Child(c)
			}
		}
	}

	if p.Depth != 0 {
		root = tree.CopyToDepth(root, p.Depth)
	}

	*res = *root
	return nil
}

type CoverageSummaryParams struct {
	// root url to work from
	// TODO - support root param
	// Root string
	// patterns to filter results against, optional
	Patterns []string
	// ids of repositories to limit query to, default is all
	RepoIds []string
}

func (p *CoverageSummaryParams) Validate() error {
	return nil
}

func (CoverageRequests) Summary(p *CoverageSummaryParams, res *Summary) error {
	sources := make([]*archive.Source, len(p.Patterns))
	for i, p := range p.Patterns {
		sources[i] = &archive.Source{Url: p}
	}
	summary, err := NewCoverageGenerator(p.RepoIds, p.Patterns).Summary()
	if err != nil {
		return err
	}

	*res = *summary
	return nil
}
