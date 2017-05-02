package tree

import (
	"testing"
	//"reflect"
	gtree "github.com/google/btree"
	"flag"
	"math/rand"
)

func btree123() *Btree {
	btree := New()
	btree.Insert(1).Insert(2).Insert(3)
	return btree
}

func btreeRandom(n int) *Btree {
	btree := New()
	for _, v := range rand.Perm(n) {
		btree.Insert(v)
	}
	return btree
}

func TestBtree1(t *testing.T) {
	btree := btreeRandom(1000)

	btree.Debug()
}

//func TestBtree2(t *testing.T) {
//	btree := btree123().InsertAll([]interface{}{4, 5, 6})
//	s1 := btree.String()
//	s2 := "[1 2 3 4 5 6 a b c]"
//	if s1 != s2 {
//		t.Error(s1, "tree string representation should equal", s2)
//	}
//}
//
//func TestIterateBtree(t *testing.T) {
//	const capacity = 3
//	btree := btree123()
//
//	var b [capacity]int
//	btree.Ascend(func(n *Node, i int) bool {
//		b[i] = n.Value.(int)
//		return true
//	})
//	c := [capacity]int{1, 2, 3}
//	if !reflect.DeepEqual(c , b) {
//		t.Error(c, "should equal", b)
//	}
//
//	btree.Descend(func(n *Node, i int) bool {
//		b[i] = n.Value.(int)
//		return true
//	})
//	d := [capacity]int{3, 2, 1}
//	if !reflect.DeepEqual(b , d) {
//		t.Error(b, "should equal", d)
//	}
//
//	l := btree.Values()
//	for i, v := range l {
//		if c[i] != v {
//			t.Error(c[i], "should equal", v)
//		}
//	}
//}

//func TestRemoveBtree(t *testing.T) {
//	btree := btree123()
//	test := []interface{}{1, 2, 3}
//	if !btree.ContainsAll(test) {
//		t.Error("tree should contain all of", test)
//	}
//	if !btree.ContainsAny(test) {
//		t.Error("tree should contain one of", test)
//	}
//
//	btree.DeleteAll(test)
//
//	if !btree.Empty() {
//		t.Error("tree should be empty")
//	}
//	if btree.ContainsAny(test) {
//		t.Error("tree should not contain any of", test)
//	}
//	btree.Init()
//	if btree.Contains(1) ||  btree.Get(1) != nil {
//		t.Error("tree should be empty")
//	}
//
//	btree = btree123()
//	pop := btree.Pop()
//	if pop != 3 {
//		t.Error(pop, "should be 3")
//	}
//	pull := btree.Pull()
//	if pull != 1 {
//		t.Error(pop, "should be 3")
//	}
//	if !btree.Delete(btree.Pop()).Empty() {
//		t.Error("tree should be empty")
//	}
//}

//type TestKey1 struct {
//	Name string
//}
//func (testkey TestKey1) Comp(val interface{}) int {
//	var c int = 0
//	tk := val.(TestKey1)
//	if testkey.Name > tk.Name {
//		c = 1
//	} else if testkey.Name < tk.Name {
//		c = -1
//	}
//	return c
//}
//type TestKey2 struct {
//	Name string
//}
//func (testkey TestKey2) String() string {
//	return testkey.Name
//}
//
//func TestCustomKeyBtree(t *testing.T) {
//	btree := New()
//	btree.InsertAll([]interface{}{TestKey1{Name: "Ross"}, TestKey1{Name: "Michael"},
//		TestKey1{Name: "Angelo"}, TestKey1{Name: "Jason"}})
//
//	rootName := btree.root.Value.(TestKey1).Name
//	if btree.root.Value.(TestKey1).Name != "Michael" {
//		t.Error(rootName, "should equal Michael")
//	}
//	btree.Init()
//	btree.InsertAll([]interface{}{TestKey2{Name: "Ross"}, TestKey2{Name: "Michael"},
//		TestKey2{Name: "Angelo"}, TestKey2{Name: "Jason"}})
//	btree.Debug()
//	s := btree.String()
//	test := "[Angelo Jason Michael Ross]"
//	if s != test {
//		t.Error(s, "should equal", test)
//	}
//
//	btree.Delete(TestKey2{Name: "Michael"})
//	if btree.Len() != 3 {
//		t.Error("tree length should be 3")
//	}
//	test = "Jason"
//	if btree.root.Value.(TestKey2).Name != test {
//		t.Error(btree.root.Value, "root of the tree should be", test)
//	}
//	for !btree.Empty() {
//		btree.Delete(btree.root.Value)
//	}
//	btree.Debug()
//}
//
//func TestCustomDuplicatesBtree(t *testing.T) {
//	btree := New()
//	btree.InsertAll([]interface{}{0, 2, 5, 10, 15, 20, 12, 14, 13, 25, 0, 2, 5, 10, 15, 20, 12, 14, 13, 25})
//	test := 10
//	if btree.Len() != 10 {
//		t.Error("tree length should be", test)
//	}
//}

// benchmark tests comparing google btree
const benchLen = 1000
var btreeDegree = flag.Int("degree", 32, "B-Tree degree")
// perm returns a random permutation of n Int items in the range [0, n).
func perm(n int) (out []gtree.Item) {
	for _, v := range rand.Perm(n) {
		out = append(out, gtree.Int(v))
	}
	return
}

var bt *Btree
var gt *gtree.BTree
var btPerm []int
var gtPerm []gtree.Item

func BenchmarkInsertBtree(b *testing.B)  {
	btree := New()
	for i := 0; i < benchLen; i++ {
		btree.Insert(i)
	}
}

func BenchmarkInsertGtree(b *testing.B)  {
	btree := gtree.New(*btreeDegree)
	for i := gtree.Int(0); i < benchLen; i++ {
		btree.ReplaceOrInsert(i)
	}
}

func BenchmarkInsertRandomBtree(b *testing.B)  {
	bt = New()
	btPerm = rand.Perm(benchLen)
	for _, v := range btPerm {
		bt.Insert(v)
	}
}

func BenchmarkInsertRandomGtree(b *testing.B)  {
	gt = gtree.New(*btreeDegree)
	gtPerm = perm(benchLen)
	for _, v := range gtPerm {
		gt.ReplaceOrInsert(v)
	}
}

//func BenchmarkGetBtree(b *testing.B)  {
//	for _, v := range btPerm {
//		bt.Get(v)
//	}
//}
//
//func BenchmarkGetGtree(b *testing.B)  {
//	for _, v := range gtPerm {
//		gt.Get(v)
//	}
//}

//func BenchmarkIterationBtree(b *testing.B)  {
//	bt.Ascend(func(n *Node, i int) bool {
//		return true
//	})
//}
//
//func BenchmarkIterationGtree(b *testing.B)  {
//	gt.Ascend(func(a gtree.Item) bool {
//		return true
//	})
//}


//func BenchmarkDeleteBtree(b *testing.B)  {
//	for _, v := range btPerm {
//		bt.Remove(v)
//	}
//}
//
//func BenchmarkDeleteGtree(b *testing.B)  {
//	for _, v := range gtPerm {
//		gt.Delete(v)
//	}
//}

