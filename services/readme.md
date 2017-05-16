# Services

This services directory is intended to give a clean overview of what each service brings to the table in terms of coverage. Each service lives in a folder that implements the `CoverageService` interface:

```go
// A Coverage Service is any service that can also provide coverage information
type CoverageService interface {
  Service
  AddUrls(tree *tree.Node) error
  AddCoverage(tree *tree.Node)
}
```

When the server starts a function in `coverage.go` builds the coverage tree. It starts an empty tree, then iterates through each service, calling both `AddUrls` & `AddCoverage` so each service gets a chance to add it's info.

The `AddUrls` function takes a tree of urls, and adds nodes for all urls that the service may know about. We're selective about what is included here for the sake of sanity, and each service makes it's own "decisions" about what is & isn't a viable url for inclusion. It's worth noting that this is expected to be a list of urls that the service "knows about", not coverage information. That's the job of `AddCoverage`

`AddCoverage` attaches coverage information to the tree.

Once each service has added it's information, the function in `coverage.go` walks the completed tree, tabulating coverage information.