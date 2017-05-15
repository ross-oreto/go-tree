/*
Copyright (c) 2017 Ross Oreto

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package tree

import (
	"fmt"
)

// Btree represents an AVL tree
type Btree struct {
	root    *Node
	values  []interface{}
	len     int
	compare Comp
}

// CompareTo interface to define the compare method used to insert and find values
type CompareTo interface {
	Comp(val interface{}) int8
}

// Comp compare function expressed as a type
type Comp func(v1, v2 interface{}) int8

// Node represents a node in the tree with a value, left and right children, and a height/balance of the node.
type Node struct {
	Value       interface{}
	left, right *Node
	height      int8
}

// New returns a new btree which expects types that implement the CompareTo or Stringer Interfaces
func New() *Btree { return new(Btree).Init() }

// NewInt returns a new btree which expects int types
func NewInt() *Btree { return new(Btree).InitWithCompare(intComp) }

// NewString returns a new btree which expects string types
func NewString() *Btree { return new(Btree).InitWithCompare(stringComp) }

// NewUint returns a new btree which expects uint types
func NewUint() *Btree { return new(Btree).InitWithCompare(uintComp) }

// NewFloat32 returns a new btree which expects float32 types
func NewFloat32() *Btree { return new(Btree).InitWithCompare(float32Comp) }

// NewFloat64 returns a new btree which expects float32 types
func NewFloat64() *Btree { return new(Btree).InitWithCompare(float64Comp) }

// NewUintptr returns a new btree which expects uintptr types
func NewUintptr() *Btree { return new(Btree).InitWithCompare(uintptrComp) }

// NewRune returns a new btree which expects rune types
func NewRune() *Btree { return new(Btree).InitWithCompare(runeComp) }

// NewByte returns a new btree which expects byte types
func NewByte() *Btree { return new(Btree).InitWithCompare(byteComp) }

// NewComplex64 returns a new btree which expects complex64 types
func NewComplex64() *Btree { return new(Btree).InitWithCompare(complex64Comp) }

// NewComplex128 returns a new btree which expects complex128 types
func NewComplex128() *Btree { return new(Btree).InitWithCompare(complex128Comp) }

// NewStringPtr returns a new btree which expects *string types
func NewStringPtr() *Btree { return new(Btree).InitWithCompare(stringPtrComp) }

// NewUintPtr returns a new btree which expects *uint types
func NewUintPtr() *Btree { return new(Btree).InitWithCompare(uintPtrComp) }

// NewIntPtr returns a new btree which expects *int types
func NewIntPtr() *Btree { return new(Btree).InitWithCompare(intPtrComp) }

// NewBytePtr returns a new btree which expects *byte types
func NewBytePtr() *Btree { return new(Btree).InitWithCompare(bytePtrComp) }

// NewRunePtr returns a new btree which expects *rune types
func NewRunePtr() *Btree { return new(Btree).InitWithCompare(runePtrComp) }

// NewFloat32Ptr returns a new btree which expects *flost32 types
func NewFloat32Ptr() *Btree { return new(Btree).InitWithCompare(float32PtrComp) }

// NewFloat64Ptr returns a new btree which expects *flost64 types
func NewFloat64Ptr() *Btree { return new(Btree).InitWithCompare(float64PtrComp) }

// NewComplex32Ptr returns a new btree which expects *complex32 types
func NewComplex32Ptr() *Btree { return new(Btree).InitWithCompare(complex32PtrComp) }

// NewComplex64Ptr returns a new btree which expects *complex64 types
func NewComplex64Ptr() *Btree { return new(Btree).InitWithCompare(complex64PtrComp) }

// Init initializes all values/clears the tree using the default compare method and returns the tree pointer
func (t *Btree) Init() *Btree {
	t.root = nil
	t.values = nil
	t.len = 0
	t.compare = comp
	return t
}

// InitWithCompare initializes all values/clears the tree using the specified compare method and returns the tree pointer
func (t *Btree) InitWithCompare(compare Comp) *Btree {
	t.Init()
	t.compare = compare
	return t
}

// String returns a string representation of the tree values
func (t *Btree) String() string {
	return fmt.Sprint(t.Values())
}

// Empty returns true if the tree is empty
func (t *Btree) Empty() bool {
	return t.root == nil
}

// NotEmpty returns true if the tree is not empty
func (t *Btree) NotEmpty() bool {
	return t.root != nil
}

func (t *Btree) balance() int8 {
	if t.root != nil {
		return balance(t.root)
	}
	return 0
}

// Insert inserts a new value into the tree and returns the tree pointer
func (t *Btree) Insert(value interface{}) *Btree {
	added := false
	t.root = insert(t.root, value, &added, t.compare)
	if added {
		t.len++
	}
	t.values = nil
	return t
}

func insert(n *Node, value interface{}, added *bool, compare Comp) *Node {
	if n == nil {
		*added = true
		return (&Node{Value: value}).Init()
	}
	c := compare(value, n.Value)
	if c > 0 {
		n.right = insert(n.right, value, added, compare)
	} else if c < 0 {
		n.left = insert(n.left, value, added, compare)
	} else {
		n.Value = value
		*added = false
		return n
	}

	n.height = n.maxHeight() + 1
	c = balance(n)

	if c > 1 {
		c = compare(value, n.left.Value)
		if c < 0 {
			return n.rotateRight()
		} else if c > 0 {
			n.left = n.left.rotateLeft()
			return n.rotateRight()
		}
	} else if c < -1 {
		c = compare(value, n.right.Value)
		if c > 0 {
			return n.rotateLeft()
		} else if c < 0 {
			n.right = n.right.rotateRight()
			return n.rotateLeft()
		}
	}
	return n
}

// InsertAll inserts all the values into the tree and returns the tree pointer
func (t *Btree) InsertAll(values []interface{}) *Btree {
	for _, v := range values {
		t.Insert(v)
	}
	return t
}

// Contains returns true if the tree contains the specified value
func (t *Btree) Contains(value interface{}) bool {
	return t.Get(value) != nil
}

// ContainsAny returns true if the tree contains any of the values
func (t *Btree) ContainsAny(values []interface{}) bool {
	for _, v := range values {
		if t.Contains(v) {
			return true
		}
	}
	return false
}

// ContainsAll returns true if the tree contains all of the values
func (t *Btree) ContainsAll(values []interface{}) bool {
	for _, v := range values {
		if !t.Contains(v) {
			return false
		}
	}
	return true
}

// Get returns the node value associated with the search value
func (t *Btree) Get(value interface{}) interface{} {
	var node *Node
	if t.root != nil {
		node = t.root.get(value, t.compare)
	}
	if node != nil {
		return node.Value
	}
	return nil
}

// Len return the number of nodes in the tree
func (t *Btree) Len() int {
	return t.len
}

// Head returns the first value in the tree
func (t *Btree) Head() interface{} {
	if t.root == nil {
		return nil
	}
	var beginning = t.root
	for beginning.left != nil {
		beginning = beginning.left
	}
	if beginning == nil {
		for beginning.right != nil {
			beginning = beginning.right
		}
	}
	if beginning != nil {
		return beginning.Value
	}
	return nil
}

// Tail returns the last value in the tree
func (t *Btree) Tail() interface{} {
	if t.root == nil {
		return nil
	}
	var beginning = t.root
	for beginning.right != nil {
		beginning = beginning.right
	}
	if beginning == nil {
		for beginning.left != nil {
			beginning = beginning.left
		}
	}
	if beginning != nil {
		return beginning.Value
	}
	return nil
}

// Values returns a slice of all the values in tree in order
func (t *Btree) Values() []interface{} {
	if t.values == nil {
		t.values = make([]interface{}, t.len)
		t.Ascend(func(n *Node, i int) bool {
			t.values[i] = n.Value
			return true
		})
	}
	return t.values
}

// Delete deletes the node from the tree associated with the search value
func (t *Btree) Delete(value interface{}) *Btree {
	deleted := false
	t.root = deleteNode(t.root, value, &deleted, t.compare)
	if deleted {
		t.len--
	}
	t.values = nil
	return t
}

// DeleteAll deletes the nodes from the tree associated with the search values
func (t *Btree) DeleteAll(values []interface{}) *Btree {
	for _, v := range values {
		t.Delete(v)
	}
	return t
}

func deleteNode(n *Node, value interface{}, deleted *bool, compare Comp) *Node {
	if n == nil {
		return n
	}

	c := compare(value, n.Value)

	if c < 0 {
		n.left = deleteNode(n.left, value, deleted, compare)
	} else if c > 0 {
		n.right = deleteNode(n.right, value, deleted, compare)
	} else {
		if n.left == nil {
			t := n.right
			n.Init()
			return t
		} else if n.right == nil {
			t := n.left
			n.Init()
			return t
		}
		t := n.right.min()
		n.Value = t.Value
		n.right = deleteNode(n.right, t.Value, deleted, compare)
		*deleted = true
	}

	//re-balance
	if n == nil {
		return n
	}
	n.height = n.maxHeight() + 1
	bal := balance(n)
	if bal > 1 {
		if balance(n.left) >= 0 {
			return n.rotateRight()
		}
		n.left = n.left.rotateLeft()
		return n.rotateRight()
	} else if bal < -1 {
		if balance(n.right) <= 0 {
			return n.rotateLeft()
		}
		n.right = n.right.rotateRight()
		return n.rotateLeft()
	}

	return n
}

// Pop deletes the last node from the tree and returns its value
func (t *Btree) Pop() interface{} {
	value := t.Tail()
	if value != nil {
		t.Delete(value)
	}
	return value
}

// Pull deletes the first node from the tree and returns its value
func (t *Btree) Pull() interface{} {
	value := t.Head()
	if value != nil {
		t.Delete(value)
	}
	return value
}

// NodeIterator expresses the iterator function used for traversals
type NodeIterator func(n *Node, i int) bool

// Ascend performs an ascending order traversal of the tree calling the iterator function on each node
// the iterator will continue as long as the NodeIterator returns true
func (t *Btree) Ascend(iterator NodeIterator) {
	var i int
	if t.root != nil {
		t.root.iterate(iterator, &i, true)
	}
}

// Descend performs a descending order traversal of the tree using the iterator
// the iterator will continue as long as the NodeIterator returns true
func (t *Btree) Descend(iterator NodeIterator) {
	var i int
	if t.root != nil {
		t.root.rIterate(iterator, &i, true)
	}
}

// Debug prints out useful debug information about the tree for debugging purposes
func (t *Btree) Debug() {
	fmt.Println("----------------------------------------------------------------------------------------------")
	if t.Empty() {
		fmt.Println("tree is empty")
	} else {
		fmt.Println(t.Len(), "elements")
	}

	t.Ascend(func(n *Node, i int) bool {
		if t.root.Value == n.Value {
			fmt.Print("ROOT ** ")
		}
		n.Debug()
		return true
	})
	fmt.Println("----------------------------------------------------------------------------------------------")
}

// Init initializes the values of the node or clears the node and returns the node pointer
func (n *Node) Init() *Node {
	n.height = 1
	n.left = nil
	n.right = nil
	return n
}

// String returns a string representing the node
func (n *Node) String() string {
	return fmt.Sprint(n.Value)
}

// Debug prints out useful debug information about the tree node for debugging purposes
func (n *Node) Debug() {
	var children string
	if n.left == nil && n.right == nil {
		children = "no children |"
	} else if n.left != nil && n.right != nil {
		children = fmt.Sprint("left child:", n.left.String(), " right child:", n.right.String())
	} else if n.right != nil {
		children = fmt.Sprint("right child:", n.right.String())
	} else {
		children = fmt.Sprint("left child:", n.left.String())
	}

	fmt.Println(n.String(), "|", "height", n.height, "|", "balance", balance(n), "|", children)
}

func height(n *Node) int8 {
	if n != nil {
		return n.height
	}
	return 0
}

func balance(n *Node) int8 {
	if n == nil {
		return 0
	}
	return height(n.left) - height(n.right)
}

func (n *Node) get(val interface{}, compare Comp) *Node {
	var node *Node
	c := compare(val, n.Value)
	if c < 0 {
		if n.left != nil {
			node = n.left.get(val, compare)
		}
	} else if c > 0 {
		if n.right != nil {
			node = n.right.get(val, compare)
		}
	} else {
		node = n
	}
	return node
}

func (n *Node) rotateRight() *Node {
	l := n.left
	// Rotation
	l.right, n.left = n, l.right

	// update heights
	n.height = n.maxHeight() + 1
	l.height = l.maxHeight() + 1

	return l
}

func (n *Node) rotateLeft() *Node {
	r := n.right
	// Rotation
	r.left, n.right = n, r.left

	// update heights
	n.height = n.maxHeight() + 1
	r.height = r.maxHeight() + 1

	return r
}

func (n *Node) iterate(iterator NodeIterator, i *int, cont bool) {
	if n != nil && cont {
		n.left.iterate(iterator, i, cont)
		cont = iterator(n, *i)
		*i++
		n.right.iterate(iterator, i, cont)
	}
}

func (n *Node) rIterate(iterator NodeIterator, i *int, cont bool) {
	if n != nil && cont {
		n.right.iterate(iterator, i, cont)
		cont = iterator(n, *i)
		*i++
		n.left.iterate(iterator, i, cont)
	}
}

func (n *Node) min() *Node {
	current := n
	for current.left != nil {
		current = current.left
	}
	return current
}

func (n *Node) maxHeight() int8 {
	rh := height(n.right)
	lh := height(n.left)
	if rh > lh {
		return rh
	}
	return lh
}

func intComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(int), v2.(int)
	if t1 > t2 {
		return 1
	} else if t1 < t2 {
		return -1
	} else {
		return 0
	}
}
func stringComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(string), v2.(string)
	if t1 > t2 {
		return 1
	} else if t1 < t2 {
		return -1
	} else {
		return 0
	}
}
func uintComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(uint), v2.(uint)
	if t1 > t2 {
		return 1
	} else if t1 < t2 {
		return -1
	} else {
		return 0
	}
}
func float32Comp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(float32), v2.(float32)
	if t1 > t2 {
		return 1
	} else if t1 < t2 {
		return -1
	} else {
		return 0
	}
}
func float64Comp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(float64), v2.(float64)
	if t1 > t2 {
		return 1
	} else if t1 < t2 {
		return -1
	} else {
		return 0
	}
}
func uintptrComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(uintptr), v2.(uintptr)
	if t1 > t2 {
		return 1
	} else if t1 < t2 {
		return -1
	} else {
		return 0
	}
}
func byteComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(byte), v2.(byte)
	if t1 > t2 {
		return 1
	} else if t1 < t2 {
		return -1
	} else {
		return 0
	}
}
func runeComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(rune), v2.(rune)
	if t1 > t2 {
		return 1
	} else if t1 < t2 {
		return -1
	} else {
		return 0
	}
}
func complex64Comp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(float64), v2.(float64)
	if t1 > t2 {
		return 1
	} else if t1 < t2 {
		return -1
	} else {
		return 0
	}
}
func complex128Comp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(float64), v2.(float64)
	if t1 > t2 {
		return 1
	} else if t1 < t2 {
		return -1
	} else {
		return 0
	}
}
func stringPtrComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(*string), v2.(*string)
	if *t1 > *t2 {
		return 1
	} else if *t1 < *t2 {
		return -1
	} else {
		return 0
	}
}
func intPtrComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(*int), v2.(*int)
	if *t1 > *t2 {
		return 1
	} else if *t1 < *t2 {
		return -1
	} else {
		return 0
	}
}
func uintPtrComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(*uint), v2.(*uint)
	if *t1 > *t2 {
		return 1
	} else if *t1 < *t2 {
		return -1
	} else {
		return 0
	}
}
func bytePtrComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(*byte), v2.(*byte)
	if *t1 > *t2 {
		return 1
	} else if *t1 < *t2 {
		return -1
	} else {
		return 0
	}
}
func runePtrComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(*rune), v2.(*rune)
	if *t1 > *t2 {
		return 1
	} else if *t1 < *t2 {
		return -1
	} else {
		return 0
	}
}
func float32PtrComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(*float32), v2.(*float32)
	if *t1 > *t2 {
		return 1
	} else if *t1 < *t2 {
		return -1
	} else {
		return 0
	}
}
func float64PtrComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(*float64), v2.(*float64)
	if *t1 > *t2 {
		return 1
	} else if *t1 < *t2 {
		return -1
	} else {
		return 0
	}
}
func complex32PtrComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(*float64), v2.(*float64)
	if *t1 > *t2 {
		return 1
	} else if *t1 < *t2 {
		return -1
	} else {
		return 0
	}
}
func complex64PtrComp(v1, v2 interface{}) int8 {
	t1, t2 := v1.(*float64), v2.(*float64)
	if *t1 > *t2 {
		return 1
	} else if *t1 < *t2 {
		return -1
	} else {
		return 0
	}
}

func comp(v1, v2 interface{}) int8 {
	var c int8
	switch v1.(type) {
	case CompareTo:
		c = v1.(CompareTo).Comp(v2)
	case fmt.Stringer:
		s1, s2 := v1.(fmt.Stringer).String(), v2.(fmt.Stringer).String()
		if s1 > s2 {
			c = 1
		} else if s1 < s2 {
			c = -1
		} else {
			c = 0
		}
	}
	return c
}
