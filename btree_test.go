package tree

import (
	"testing"
	"reflect"
)

func btree123() *Btree {
	btree := New()
	btree.Insert(1).Insert(2).Insert(3)
	return btree
}

func TestBtree1(t *testing.T) {
	btree := btree123()

	if btree.Size() != 3 {
		t.Error("size of tree should be 3")
	}
	if btree.Get(2) != 2 {
		t.Error("get value should be 2")
	}
	if btree.Head().Value != 1 {
		t.Error("beggining of tree should be 1")
	}
	if btree.Root.Value != 2 {
		t.Error("Root of tree should be 2")
	}
}

func TestBtree2(t *testing.T) {
	btree := btree123().InsertAll([]interface{}{4, 5, 6})
	btree.PutAll(map[interface{}]interface{}{ 7:"a", 8:"b", 9:"c" })
	s1 := btree.String()
	s2 := "[1 2 3 4 5 6 a b c]"
	if s1 != s2 {
		t.Error(s1, "tree string representation should equal", s2)
	}
}

func TestIterateBtree(t *testing.T) {
	const capacity = 3
	btree := btree123()

	var b [capacity]int
	for i, n := 0, btree.Head(); n != nil; i, n = i + 1, n.Next() {
		b[i] = n.Value.(int)
	}
	c := [capacity]int{1, 2, 3}
	if !reflect.DeepEqual(c , b) {
		t.Error(c, "should equal", b)
	}

	for i, n := 0, btree.Tail(); n != nil; i, n = i + 1, n.Prev() {
		b[i] = n.Value.(int)
	}
	d := [capacity]int{3, 2, 1}
	if !reflect.DeepEqual(b , d) {
		t.Error(b, "should equal", d)
	}

	l := btree.Values()
	for i, v := range l {
		if c[i] != v {
			t.Error(c[i], "should equal", v)
		}
	}
}

func TestRemoveBtree(t *testing.T) {
	btree := btree123()
	btree.Remove(1).Remove(2).Remove(3)

	if !btree.Empty() {
		t.Error("tree should be empty")
	}
}