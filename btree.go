package tree

import (
	"fmt"
)

type Btree struct {
	root *Node
	values []interface{}
	len int
}

type Comparer interface {
	Comp(val interface{}) int8
}

type Node struct {
	Value interface{}
	left, right *Node
	height int8
}

func New() *Btree { return new(Btree).Init() }

func (t *Btree) Init() *Btree {
	t.root = nil
	t.values = nil
	t.len = 0
	return t
}

func (t *Btree) String() string {
	return fmt.Sprint(t.Values())
}

func (t *Btree) Empty() bool {
	return t.root == nil
}

func (t *Btree) NotEmpty() bool {
	return t.root != nil
}

func (t *Btree) balance() int8 {
	if t.root != nil {
		return balance(t.root)
	}
	return 0
}

func (t *Btree) Insert(value interface{}) *Btree {
	added := false
	t.root = insert(t.root, value, &added)
	if added { t.len += 1 }
	t.values = nil
	return t
}

func insert(n *Node, value interface{}, added *bool) *Node {
	if n == nil {
		*added = true
		return (&Node{Value: value}).Init()
	} else {
		c := Comp(value, n.Value)
		if c > 0 {
			n.right = insert(n.right, value, added)
		} else if c < 0 {
			n.left = insert(n.left, value, added)
		} else {
			n.Value = value
			*added = false
			return n
		}

		n.height = n.maxHeight() + 1
		c = balance(n)

		if c > 1 {
			c = Comp(value, n.left.Value)
			if c < 0 {
				return n.rotateRight()
			} else if c > 0 {
				n.left = n.left.rotateLeft()
				return n.rotateRight()
			}
		} else if c < -1 {
			c = Comp(value, n.right.Value)
			if c > 0 {
				return n.rotateLeft()
			} else if c < 0 {
				n.right = n.right.rotateRight()
				return n.rotateLeft()
			}
		}
	}
	return n
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
	return t.len
}

func (t *Btree) Head() interface{} {
	if t.root == nil { return nil }
	var beginning *Node = t.root
	for beginning.left != nil {
		beginning = beginning.left
	}
	if beginning == nil {
		for beginning.right != nil {
			beginning = beginning.right
		}
	}
	if beginning != nil { return beginning.Value }
	return nil
}

func (t *Btree) Tail() interface{} {
	if t.root == nil { return nil }
	var beginning *Node = t.root
	for beginning.right != nil {
		beginning = beginning.right
	}
	if beginning == nil {
		for beginning.left != nil {
			beginning = beginning.left
		}
	}
	if beginning != nil { return beginning.Value }
	return nil
}

func (t *Btree) Values() []interface{} {
	if t.values == nil {
		t.Ascend(func(n *Node, i int) bool {
			t.values = append(t.values, n.Value)
			return true
		})
	}
	return t.values
}

func (t *Btree) Delete(value interface{}) *Btree {
	deleted := false
	t.root = deleteNode(t.root, value, &deleted)
	if deleted { t.len -= 1 }
	t.values = nil
	return t
}

func (t *Btree) DeleteAll(values []interface{}) *Btree {
	for _, v := range values {
		t.Delete(v)
	}
	return t
}

func deleteNode(n *Node, value interface{}, deleted *bool) *Node {
	if n == nil { return n }

	c := Comp(value, n.Value)

	if c < 0 {
		n.left = deleteNode(n.left, value, deleted)
	} else if c > 0 {
		n.right = deleteNode(n.right, value, deleted)
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
		n.right = deleteNode(n.right, t.Value, deleted)
		*deleted = true
	}

	//re-balance
	if n == nil { return n }
	n.height = n.maxHeight() + 1
	bal := balance(n)
	if bal > 1 {
		if balance(n.left) >= 0 {
			return n.rotateRight()
		} else {
			n.left =  n.left.rotateLeft()
			return n.rotateRight()
		}
	} else if bal < -1 {
		if balance(n.right) <= 0 {
			return n.rotateLeft()
		} else {
			n.right = n.right.rotateRight()
			return n.rotateLeft()
		}
	}

	return n
}

func (t *Btree) Pop() interface{} {
	value := t.Tail()
	if value != nil { t.Delete(value) }
	return value
}

func (t *Btree) Pull()  interface{} {
	value := t.Head()
	if value != nil { t.Delete(value) }
	return value
}

type NodeIterator func(n *Node, i int) bool

func (t *Btree) Ascend(iterator NodeIterator) {
	var i int = 0
	if t.root != nil { t.root.iterate(iterator, &i, true) }
}

func (t *Btree) Descend(iterator NodeIterator) {
	var i int = 0
	if t.root != nil { t.root.r_iterate(iterator, &i, true) }
}

func (t *Btree) Debug() {
	fmt.Println("----------------------------------------------------------------------------------------------")
	if t.Empty() {
		fmt.Println("tree is empty")
	} else { fmt.Println(t.Len(), "elements") }

	t.Ascend(func(n *Node, i int) bool {
		if Comp(t.root.Value, n.Value) == 0 {
			fmt.Print("ROOT ** ")
		}
		n.Debug()
		return true
	})
	fmt.Println("----------------------------------------------------------------------------------------------")
}

func (n *Node) Init() *Node {
	n.height = 1
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

	fmt.Println(n.String(), "|", "height", n.height, "|", "balance", balance(n), "|", children)
}

func height(n *Node) int8 {
	if n != nil {
		return n.height
	}
	return 0
}

func balance(n *Node) int8 {
	if n == nil { return 0 }
	return height(n.left) - height(n.right)
}

func (n *Node) get(val interface{}) *Node {
	var node *Node = nil
	c := Comp(val, n.Value)
	if c < 0 {
		if n.left != nil { node = n.left.get(val) }
	} else if c > 0 {
		if n.right != nil { node = n.right.get(val) }
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
	if rh > lh { return rh }
	return lh
}

func Comp(v1, v2 interface{}) int8  {
	var c int8 = 0

	switch v1.(type) {
	default:
		t1, t2 := v1.(string), v2.(string)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case int:
		t1, t2 := v1.(int), v2.(int)
		if t1 > t2 { c = 1 } else if t1 < t2 { c = -1 } else { c = 0 }
	case string:
		t1, t2 := v1.(string), v2.(string)
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
	case *string:
		t1, t2 := v1.(*string), v2.(*string)
		if *t1 > *t2 { c = 1 } else if *t1 < *t2 { c = -1 } else { c = 0 }
	case *int:
		t1, t2 := v1.(*int), v2.(*int)
		if *t1 > *t2 { c = 1 } else if *t1 < *t2 { c = -1 } else { c = 0 }
	case *uint:
		t1, t2 := v1.(*uint), v2.(*uint)
		if *t1 > *t2 { c = 1 } else if *t1 < *t2 { c = -1 } else { c = 0 }
	case *float32:
		t1, t2 := v1.(*float32), v2.(*float32)
		if *t1 > *t2 { c = 1 } else if *t1 < *t2 { c = -1 } else { c = 0 }
	case *float64:
		t1, t2 := v1.(*float64), v2.(*float64)
		if *t1 > *t2 { c = 1 } else if *t1 < *t2 { c = -1 } else { c = 0 }
	case *uintptr:
		t1, t2 := v1.(*uintptr), v2.(*uintptr)
		if *t1 > *t2 { c = 1 } else if *t1 < *t2 { c = -1 } else { c = 0 }
	case *byte:
		t1, t2 := v1.(*byte), v2.(*byte)
		if *t1 > *t2 { c = 1 } else if *t1 < *t2 { c = -1 } else { c = 0 }
	case *rune:
		t1, t2 := v1.(*rune), v2.(*rune)
		if *t1 > *t2 { c = 1 } else if *t1 < *t2 { c = -1 } else { c = 0 }
	case *complex64:
		t1, t2 := v1.(*float64), v2.(*float64)
		if *t1 > *t2 { c = 1 } else if *t1 < *t2 { c = -1 } else { c = 0 }
	case *complex128:
		t1, t2 := v1.(*float64), v2.(*float64)
		if *t1 > *t2 { c = 1 } else if *t1 < *t2 { c = -1 } else { c = 0 }
	case Comparer:
		c = v1.(Comparer).Comp(v2)
	case fmt.Stringer:
		s1, s2 := v1.(fmt.Stringer).String(), v2.(fmt.Stringer).String()
		if s1 > s2 { c = 1 } else if s1 < s2 { c = -1 } else { c = 0 }
	}
	return c
}