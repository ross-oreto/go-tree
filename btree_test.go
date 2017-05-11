package go_tree

import (
	"testing"
	gtree "github.com/google/btree"
	"flag"
	"math/rand"
	"reflect"
)

func btreeInOrder(n int) *Btree {
	btree := New()
	for i := 1; i <= n; i++ {
		btree.Insert(i)
	}
	return btree
}

func btreeFixed(values []interface{}) *Btree {
	btree := New()
	btree.InsertAll(values)
	return btree
}

const benchLen = 1000000
var btreeDegree = flag.Int("degree", 32, "B-Tree degree")

func TestBtree_Get(t *testing.T) {
	values := []interface{}{ 9,4,2,6,8,0,3,1,7,5 }
	btree := btreeFixed(values).InsertAll(values)

	expect, actual := len(values), btree.Len()
	if actual != expect {
		t.Error("length should equal", expect, "actual", actual)
	}

	expect = 2
	if btree.Get(expect) == nil {
		t.Error("value should equal", expect)
	}
}

func TestBtree_Contains(t *testing.T) {
	btree := btreeInOrder(1000)

	test := 1
	if !btree.Contains(test) { t.Error("tree should contain", test) }

	test2 := []interface{}{1, 2, 3, 4}
	if !btree.ContainsAll(test2) { t.Error("tree should contain", test2) }

	test2 = []interface{}{ 5 }
	if !btree.ContainsAny(test2) { t.Error("tree should contain", test2) }

	test2 = []interface{}{ 5000, 2000 }
	if btree.ContainsAny(test2) { t.Error("tree should not contain any", test2) }
}

func TestBtree_String(t *testing.T) {
	btree := btreeFixed([]interface{}{1, 2, 3, 4, 5, 6})
	s1 := btree.String()
	s2 := "[1 2 3 4 5 6]"
	if s1 != s2 {
		t.Error(s1, "tree string representation should equal", s2)
	}
}

func TestBtree_Values(t *testing.T) {
	const capacity = 3
	btree := btreeFixed([]interface{}{1,2})

	b := btree.Values()
	c := []interface{}{1, 2}
	if !reflect.DeepEqual(c , b) {
		t.Error(c, "should equal", b)
	}
	btree.Insert(3)

	desc := [capacity]int{}
	btree.Descend(func(n *Node, i int) bool {
		desc[i] = n.Value.(int)
		return true
	})
	d := [capacity]int{3, 2, 1}
	if !reflect.DeepEqual(desc , d) {
		t.Error(desc, "should equal", d)
	}

	e := []interface{}{1, 2, 3}
	for i, v := range btree.Values() {
		if e[i] != v {
			t.Error(e[i], "should equal", v)
		}
	}
}

func TestBtree_Delete(t *testing.T) {
	test := []interface{}{1, 2, 3}
	btree := btreeFixed(test)

	btree.DeleteAll(test)

	if !btree.Empty() {
		t.Error("tree should be empty")
	}

	btree = btreeFixed(test)
	pop := btree.Pop()
	if pop != 3 {
		t.Error(pop, "should be 3")
	}
	pull := btree.Pull()
	if pull != 1 {
		t.Error(pop, "should be 3")
	}
	if !btree.Delete(btree.Pop()).Empty() {
		t.Error("tree should be empty")
	}
	btree.Pop()
	btree.Pull()
}

func TestBtree_HeadTail(t *testing.T) {
	btree := btreeFixed([]interface{}{1, 2, 3})
	if btree.Head() != 1 {
		t.Error("head element should be 1")
	}
	if btree.Tail() != 3 {
		t.Error("head element should be 3")
	}
	btree.Init()
	if btree.Head() != nil {
		t.Error("head element should be nil")
	}
}

type TestKey1 struct {
	Name string
}
func (testkey TestKey1) Comp(val interface{}) int8 {
	var c int8 = 0
	tk := val.(TestKey1)
	if testkey.Name > tk.Name {
		c = 1
	} else if testkey.Name < tk.Name {
		c = -1
	}
	return c
}
type TestKey2 struct {
	Name string
}
func (testkey TestKey2) String() string {
	return testkey.Name
}

func TestBtree_CustomKey(t *testing.T) {
	btree := New()
	btree.InsertAll([]interface{}{TestKey1{Name: "Ross"}, TestKey1{Name: "Michael"},
		TestKey1{Name: "Angelo"}, TestKey1{Name: "Jason"}})

	rootName := btree.root.Value.(TestKey1).Name
	if btree.root.Value.(TestKey1).Name != "Michael" {
		t.Error(rootName, "should equal Michael")
	}
	btree.Init()
	btree.InsertAll([]interface{}{TestKey2{Name: "Ross"}, TestKey2{Name: "Michael"},
		TestKey2{Name: "Angelo"}, TestKey2{Name: "Jason"}})
	btree.Debug()
	s := btree.String()
	test := "[Angelo Jason Michael Ross]"
	if s != test {
		t.Error(s, "should equal", test)
	}

	btree.Delete(TestKey2{Name: "Michael"})
	if btree.Len() != 3 {
		t.Error("tree length should be 3")
	}
	test = "Jason"
	if btree.root.Value.(TestKey2).Name != test {
		t.Error(btree.root.Value, "root of the tree should be", test)
	}
	for !btree.Empty() {
		btree.Delete(btree.root.Value)
	}
	btree.Debug()
}

func TestBtree_Duplicates(t *testing.T) {
	btree := New()
	btree.InsertAll([]interface{}{0, 2, 5, 10, 15, 20, 12, 14, 13, 25, 0, 2, 5, 10, 15, 20, 12, 14, 13, 25})
	test := 10
	length := btree.Len()
	if length != test {
		t.Error(length, "tree length should be", test)
	}
}

// benchmark tests comparing google btree

var bt *Btree
var gt *gtree.BTree
var btPerm []int

func BenchmarkInsertBtree(b *testing.B)  {
	btree := New()
	for i := 0; i < b.N; i++ {
		for i := 0; i < benchLen; i++ {
			btree.Insert(i)
		}
	}
}

func BenchmarkInsertGtree(b *testing.B)  {
	btree := gtree.New(*btreeDegree)
	for i := 0; i < b.N; i++ {
		for i := gtree.Int(0); i < benchLen; i++ {
			btree.ReplaceOrInsert(i)
		}
	}
}

func BenchmarkInsertRandomBtree(b *testing.B)  {
	bt = New()
	btPerm = rand.Perm(benchLen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range btPerm {
			bt.Insert(v)
		}
	}
}

func BenchmarkInsertRandomGtree(b *testing.B)  {
	gt = gtree.New(*btreeDegree)
	for i := 0; i < b.N; i++ {
		for _, v := range btPerm {
			gt.ReplaceOrInsert(gtree.Int(v))
		}
	}
}

func BenchmarkGetBtree(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		for _, v := range btPerm {
			bt.Get(v)
		}
	}
}

func BenchmarkGetGtree(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		for _, v := range btPerm {
			gt.Get(gtree.Int(v))
		}
	}
}

func BenchmarkIterationBtree(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		length := len(bt.Values())
		for i := 0; i < length; i++ {
			if bt.values[i] != nil {}
		}
	}
}

func BenchmarkIterationGtree(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		gt.Ascend(func(a gtree.Item) bool {
			return true
		})
	}
}

func BenchmarkLenBtree(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		bt.Len()
	}
}

func BenchmarkLenGtree(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		gt.Len()
	}
}

func BenchmarkDeleteBtree(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		for _, v := range btPerm {
			bt.Delete(v)
		}
	}
}

func BenchmarkDeleteGtree(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		for _, v := range btPerm {
			gt.Delete(gtree.Int(v))
		}
	}
}

