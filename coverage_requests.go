package main

import (
	"github.com/archivers-space/archive"
	"github.com/archivers-space/coverage/tree"
)

type Coverage int

type CoverageTreeArgs struct {
	Pattern string
}

func (args *CoverageTreeArgs) Validate() error {
	return nil
}

func (Coverage) Tree(args *CoverageTreeArgs, res *tree.Node) error {
	root, err := NewCoverageGenerator().Tree(&archive.Source{Url: args.Pattern})
	if err != nil {
		return err
	}

	*res = *root
	return nil
}

type CoverageSummaryArgs struct {
	Pattern string
}

func (args *CoverageSummaryArgs) Validate() error {
	return nil
}

func (Coverage) Summary(args *CoverageSummaryArgs, res *CoverageSummary) error {
	summary, err := NewCoverageGenerator().Summary(&archive.Source{Url: args.Pattern})
	if err != nil {
		return err
	}

	*res = *summary
	return nil
}
