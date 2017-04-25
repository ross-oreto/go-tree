package tree

import (
	"testing"
)

func TestBtree(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	btree := New()
	btree.Insert(1).Insert(2).Insert(3)

	if btree.Size() != 3 {
		t.Error("size of tree should be 3")
	}
	if btree.Get(2) != 2 {
		t.Error("get value should be 2")
	}
	if btree.beginning().Value != 1 {
		t.Error("beggining of tree should be 1")
	}
	if btree.root.Value != 2 {
		t.Error("root of tree should be 2")
	}
}