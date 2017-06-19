package coverage

import (
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/tree"
)

type CoverageRequests int

type CoverageTreeParams struct {
	Patterns []string
}

func (p *CoverageTreeParams) Validate() error {
	return nil
}

func (CoverageRequests) Tree(p *CoverageTreeParams, res *tree.Node) error {
	sources := make([]*archive.Source, len(p.Patterns))
	for i, p := range p.Patterns {
		sources[i] = &archive.Source{Url: p}
	}
	root, err := NewCoverageGenerator().Tree(sources...)
	if err != nil {
		return err
	}

	*res = *root
	return nil
}

type CoverageSummaryParams struct {
	Patterns []string
}

func (p *CoverageSummaryParams) Validate() error {
	return nil
}

func (CoverageRequests) Summary(p *CoverageSummaryParams, res *Summary) error {
	sources := make([]*archive.Source, len(p.Patterns))
	for i, p := range p.Patterns {
		sources[i] = &archive.Source{Url: p}
	}
	summary, err := NewCoverageGenerator().Summary(sources...)
	if err != nil {
		return err
	}

	*res = *summary
	return nil
}
