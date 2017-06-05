package coverage

import (
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/tree"
)

type CoverageRequests int

type CoverageTreeParams struct {
	Pattern string
}

func (p *CoverageTreeParams) Validate() error {
	return nil
}

func (CoverageRequests) Tree(p *CoverageTreeParams, res *tree.Node) error {
	root, err := NewCoverageGenerator().Tree(&archive.Source{Url: p.Pattern})
	if err != nil {
		return err
	}

	*res = *root
	return nil
}

type CoverageSummaryParams struct {
	Pattern string
}

func (p *CoverageSummaryParams) Validate() error {
	return nil
}

func (CoverageRequests) Summary(p *CoverageSummaryParams, res *Summary) error {
	summary, err := NewCoverageGenerator().Summary(&archive.Source{Url: p.Pattern})
	if err != nil {
		return err
	}

	*res = *summary
	return nil
}
