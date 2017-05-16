package tree

import (
	"fmt"
	"testing"
)

var testTree = &Node{
	Id:   "_",
	Name: "test",
	Children: []*Node{
		&Node{
			Id: "a",
			Children: []*Node{
				&Node{Id: "aa",
					Children: []*Node{
						&Node{Id: "aaa"},
						&Node{Id: "aab"},
					},
				},
				&Node{Id: "ab"},
			}},
		&Node{Id: "b"},
		&Node{Id: "c"},
	},
}

func TestNodeCopy(t *testing.T) {
	a := testTree

	b := a.Copy()
	b.Name = "new name"
	b.Children = nil

	if a.Name == "new name" {
		t.Errorf("copy didn't dissociate pointer")
	}
	if a.Children == nil {
		t.Errorf("copy didn't dissociate child slice")
	}
}

func maxDepth(n *Node) int {
	max := 0
	n.Walk(func(node *Node) {
		if len(node.Id) > max {
			max = len(node.Id)
		}
	})

	return max
}

func TestCopyToDepth(t *testing.T) {
	cp := CopyToDepth(testTree, 1)
	if maxDepth(cp) != 1 {
		t.Errorf("copy depth mismatch. expected: %d, got: %d", 2, maxDepth(cp))
	}
	if maxDepth(testTree) != 3 {
		t.Errorf("copyToDepth was destuctive. maxDepth: %d", maxDepth)
		// testTree.PrintTree(0, 5, " ", " ")
		// fmt.Println("------")
		// cp.PrintTree(0, 5, " ", " ")
	}
}

func CompareNodes(a, b *Node, strict bool) error {
	if a.Id != b.Id {
		return fmt.Errorf("Id mismatch. %s != %s", a.Id, b.Id)
	}
	if a.Name != b.Name {
		return fmt.Errorf("Name mismatch. %s != %s", a.Name, b.Name)
	}
	if a.NumDescendants != b.NumDescendants {
		return fmt.Errorf("NumDescendants mismatch. %d != %d", a.NumDescendants, b.NumDescendants)
	}
	if a.NumDescendantsArchived != b.NumDescendantsArchived {
		return fmt.Errorf("NumDescendantsArchived mismatch. %d != %d", a.NumDescendantsArchived, b.NumDescendantsArchived)
	}
	if a.NumChildren != b.NumChildren {
		return fmt.Errorf("NumChildren mismatch. %d != %d", a.NumChildren, b.NumChildren)
	}
	if a.Archived != b.Archived {
		return fmt.Errorf("Archived mismatch. %d != %d", a.Archived, b.Archived)
	}
	return nil
}

func CompareNodeSlices(a, b []*Node) error {
	if a == nil && b != nil || a != nil && b == nil {
		return fmt.Errorf("nil mismatch %s != %s", a, b)
	}
	if len(a) != len(b) {
		return fmt.Errorf("node length mismatch: %d != %d", len(a), len(b))
	}

	return nil
}
