package tree

import (
	"fmt"
	"github.com/pborman/uuid"
	"sort"
)

// Node is a url component in the coverage tree
type Node struct {
	Id                     string      `json:"id"`
	Name                   string      `json:"name"`
	NumDescendants         int         `json:"numDescendants,omitempty"`
	NumDescendantsArchived int         `json:"numDescendantsArchived,omitempty"`
	NumLeaves              int         `json:"numLeaves,omitempty"`
	NumLeavesArchived      int         `json:"numLeavesArchived,omitempty"`
	NumChildren            int         `json:"numChildren,omitempty"`
	Archived               bool        `json:"archived,omitempty"`
	Children               []*Node     `json:"children,omitempty"`
	Coverage               []*Coverage `json:"coverage,omitempty"`
}

func (n *Node) Child(name string) *Node {
	for _, c := range n.Children {
		if c.Name == name {
			return c
		}
	}

	c := &Node{Id: uuid.New(), Name: name}
	n.Children = append(n.Children, c)
	return c
}

func (n *Node) Copy() *Node {
	return &Node{
		Id:                     n.Id,
		Name:                   n.Name,
		NumDescendants:         n.NumDescendants,
		NumDescendantsArchived: n.NumDescendantsArchived,
		NumLeaves:              n.NumLeaves,
		NumLeavesArchived:      n.NumLeavesArchived,
		NumChildren:            n.NumChildren,
		Archived:               n.Archived,
		Children:               n.Children,
		Coverage:               n.Coverage,
	}
}

func (n *Node) SortChildren() {
	sort.Slice(n.Children, func(i, j int) bool { return n.Children[i].Name < n.Children[j].Name })
}

func (n *Node) Find(id string) (found *Node) {
	n.Walk(func(c *Node) {
		if c.Id == id {
			found = c
		}
	})
	return
}

func (n *Node) Walk(visit func(*Node)) {
	visit(n)
	for _, c := range n.Children {
		c.Walk(visit)
	}
}

func (n *Node) PrintTree(depth, maxDepth int, initial, indent string) {
	if depth == maxDepth {
		return
	}
	fmt.Printf("%s%s:%s\n", indent, n.Id, n.Name)
	for _, c := range n.Children {
		c.PrintTree(depth+1, maxDepth, initial, indent+initial)
	}
}

func CopyToDepth(node *Node, depth int) *Node {
	n := node.Copy()
	if depth == 0 {
		n.Children = nil
		return n
	} else if n.Children != nil {
		n.Children = make([]*Node, len(node.Children))
		copy(n.Children, node.Children)
		for i, c := range n.Children {
			n.Children[i] = CopyToDepth(c, depth-1)
		}
	}
	return n
}
