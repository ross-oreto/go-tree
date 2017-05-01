package tree

import (
	"fmt"
)

type Btree struct {
	root *Node
}

type Comparer interface {
	Comp(val interface{}) int
}

type Node struct {
	Value interface{}
	left, right *Node
	balance int
}

func New() *Btree { return new(Btree).Init() }

// btree methods
func (t *Btree) Init() *Btree {
	t.root = nil
	return t
}

func (t *Btree) String() string {
	return fmt.Sprint(t.Values())
}

func (t *Btree) Empty() bool {
	return t.root == nil
}

func (t *Btree) Insert(value interface{}) *Btree {
	var newNode *Node = (&Node{Value: value}).Init()
	if t.Empty() {
		t.root = newNode
	} else {
		t.root.insert(newNode)
	}
	return t
}

func (t *Btree) InsertAll(values []interface{}) *Btree {
	for _, v := range values {
		t.Insert(v)
	}
	return t
}

func (t *Btree) Contains(value interface{}) bool {
	return t.Get(value) != nil
}

func (t *Btree) ContainsAny(values []interface{}) bool {
	for _, v := range values {
		if t.Contains(v) { return true }
	}
	return false
}

func (t *Btree) ContainsAll(values []interface{}) bool {
	for _, v := range values {
		if !t.Contains(v) { return false }
	}
	return true
}

func (t *Btree) Get(value interface{}) interface{} {
	var node *Node = nil
	if t.root != nil {
		node = t.root.get(value)
	}
	if node != nil {
		return node.Value
	}
	return nil
}

func (t *Btree) Len() int {
	nodes := 0
	if !t.Empty() {
		t.count(t.root, &nodes)
	}
	return nodes
}

func (t *Btree) Head() *Node {
	var beginning *Node = t.root
	for beginning.left != nil {
		beginning = beginning.left
	}
	if beginning == nil {
		for beginning.right != nil {
			beginning = beginning.right
		}
	}
	return beginning
}

func (t *Btree) Tail() *Node {
	var beginning *Node = t.root
	for beginning.right != nil {
		beginning = beginning.right
	}
	if beginning == nil {
		for beginning.left != nil {
			beginning = beginning.left
		}
	}
	return beginning
}

func (t *Btree) Values() []interface{} {
	size := t.Len()
	slice := make([]interface{}, size)
	t.Ascend(func(n *Node, i int) bool {
		slice[i] = n.Value
		return true
	})
	return slice
}

func (t *Btree) Delete(value interface{}) *Btree {
	t.root.get(value).deleteNode()
	return t
}

func (t *Btree) DeleteAll(values []interface{}) *Btree {
	for _, k := range values {
		t.Delete(k)
	}
	return t
}

func (t *Btree) Pop() interface{} {
	var node *Node = t.Tail()

	return node.Value
}

func (t *Btree) Pull()  interface{} {
	var node *Node = t.Head()

	return node.Value
}

type NodeIterator func(n *Node, i int) bool

func (t *Btree) Ascend(iterator NodeIterator) {
	var i int = 0
	t.root.iterate(iterator, &i, true)
}

func (t *Btree) Descend(iterator NodeIterator) {
	var i int = 0
	t.root.r_iterate(iterator, &i, true)
}

func (t *Btree) Debug() {
	fmt.Println("----------------------------------------------------------------------------------------------")
	if t.Empty() {
		fmt.Println("tree is empty")
	} else { fmt.Println(t.Len(), "elements") }
	t.Ascend(func(n *Node, i int) bool {
		fmt.Print(i, "-")
		n.Debug()
		return true
	})
	fmt.Println("----------------------------------------------------------------------------------------------")
}

func (t *Btree) count(node *Node, nodes *int) {
	if node != nil {
		t.count(node.left, nodes)
		*nodes += 1
		t.count(node.right, nodes)
	}
}

func (n *Node) Init() *Node {
	n.balance = 0
	n.left = nil
	n.right = nil
	return n
}

func (n *Node) String() string {
	return fmt.Sprint(n.Value)
}

func (n *Node) Debug() {
	var children string = ""
	if n.left == nil && n.right == nil {
		children = "no children |"
	} else if n.left != nil && n.right != nil {
		children = fmt.Sprint("left child:", n.left.String(), " right child:", n.right.String())
	} else if n.right != nil {
		children = fmt.Sprint("right child:", n.right.String())
	} else {
		children = fmt.Sprint("left child:", n.left.String())
	}

	fmt.Println(n.String(), "|", "weight", n.balance, "|", children)
}

func Comp(v1, v2 interface{}) int  {
	c := 0

	switch v1.(type) {
	default:
		t1, t2 := v1.(string), v2.(string)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case int:
		t1, t2 := v1.(int), v2.(int)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case uint:
		t1, t2 := v1.(uint), v2.(uint)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case float32:
		t1, t2 := v1.(float32), v2.(float32)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case float64:
		t1, t2 := v1.(float64), v2.(float64)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case uintptr:
		t1, t2 := v1.(uintptr), v2.(uintptr)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case byte:
		t1, t2 := v1.(byte), v2.(byte)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case rune:
		t1, t2 := v1.(rune), v2.(rune)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case complex64:
		t1, t2 := v1.(float64), v2.(float64)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case complex128:
		t1, t2 := v1.(float64), v2.(float64)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case string:
		t1, t2 := v1.(string), v2.(string)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case Comparer:
		c = v1.(Comparer).Comp(v2)
	case fmt.Stringer:
		s1, s2 := v1.(fmt.Stringer).String(), v2.(fmt.Stringer).String()
		if s1 > s2 { c = 1 } else if s1 < s2 { c = -1 } else { c = 0 }
	}
	return c
}

func (n *Node) insert(newNode *Node) int {
	c := Comp(newNode.Value, n.Value)
	return c
}

func (n *Node) get(k interface{}) *Node {
	var node *Node = nil
	c := Comp(k, n.Value)
	if c < 0 {
		if n.left != nil { node = n.left.get(k) }
	} else if c > 0 {
		if n.right != nil { node = n.right.get(k) }
	} else {
		node = n
	}
	return node
}

func (n *Node) deleteNode() *Node {
	return n
}

func (n *Node) rotateLeft() *Node {
	var node *Node
	return node
}

func (n *Node) rotateRight() *Node {
	var node *Node
	return node
}

func (n *Node) iterate(iterator NodeIterator, i *int, cont bool)  {
	if n != nil && cont {
		n.left.iterate(iterator, i, cont)
		cont = iterator(n, *i)
		*i += 1
		n.right.iterate(iterator, i, cont)
	}
}

func (n *Node) r_iterate(iterator NodeIterator, i *int, cont bool)  {
	if n != nil && cont {
		n.right.iterate(iterator, i, cont)
		cont = iterator(n, *i)
		*i += 1
		n.left.iterate(iterator, i, cont)
	}
}