package tree

import (
	"math"
	"fmt"
)

type Btree struct {
	Root     *Node
}

type Traverser interface {
	Next() *Node
	Previous() *Node
}

type Comparer interface {
	Comp(val interface{}) int
}

type Node struct {
	Key, Value interface{}
	left, right, parent *Node
	balance int
	nodeType NodeType
	traverser Traverser
}

type NodeType uint8
const (
	ROOT NodeType = 0
	LEFT = 1
	RIGHT = 2
)

func New() *Btree { return new(Btree).Init() }

// btree methods
func (t *Btree) Init() *Btree {
	t.Root = nil
	return t
}

func (t *Btree) String() string {
	return fmt.Sprint(t.Values())
}

func (t *Btree) Empty() bool {
	return t.Root == nil
}

func (t *Btree) Put(key interface{}, value interface{}) *Btree {
	var newNode *Node = (&Node{Key: key, Value: value}).Init()
	newNode.nodeType = ROOT
	if t.Empty() {
		t.Root = newNode
	} else {
		t.Root.Insert(newNode)
		for t.Root.nodeType != ROOT {
			t.Root = t.Root.parent
		}
	}
	return t
}

func (t *Btree) PutAll(entries map[interface{}]interface{}) *Btree {
	for k, v := range entries {
		t.Put(k, v)
	}
	return t
}

func (t *Btree) Insert(value interface{}) *Btree {
	return t.Put(value, value)
}

func (t *Btree) InsertAll(values []interface{}) *Btree {
	for _, v := range values {
		t.Insert(v)
	}
	return t
}

func (t *Btree) Contains(key interface{}) bool {
	return t.Get(key) != nil
}

func (t *Btree) ContainsAny(keys []interface{}) bool {
	for _, k := range keys {
		if t.Contains(k) { return true }
	}
	return false
}

func (t *Btree) ContainsAll(keys []interface{}) bool {
	for _, k := range keys {
		if !t.Contains(k) { return false }
	}
	return true
}

func (t *Btree) GetNode(key interface{}) *Node {
	var node *Node = nil
	if !t.Empty() { node = t.Root.get(key) }
	return node
}

func (t *Btree) Get(key interface{}) interface{} {
	var node *Node = t.GetNode(key)
	if node != nil {
		return node.Value
	}
	return nil
}

func (t *Btree) Len() int {
	nodes := 0
	if !t.Empty() {
		t.count(t.Root, &nodes)
	}
	return nodes
}

func (t *Btree) Head() *Node {
	if t.Empty() { return nil }
	return t.Root.beginning()
}

func (t *Btree) Tail() *Node {
	if t.Empty() { return nil }
	return t.Root.end()
}

func (t *Btree) Keys() []interface{} {
	size := t.Len()
	slice := make([]interface{}, size)
	for i, n := 0, t.Head(); i < size && n != nil; i, n = i + 1, n.Next() {
		slice[i] = n.Key
	}
	return slice
}

func (t *Btree) Values() []interface{} {
	size := t.Len()
	slice := make([]interface{}, size)
	for i, n := 0, t.Head(); i < size && n != nil; i, n = i + 1, n.Next() {
		slice[i] = n.Value
	}
	return slice
}

func (t *Btree) Map() map[interface{}]interface{} {
	size := t.Len()
	btreeMap := make(map[interface{}]interface{}, size)
	for i, n := 0, t.Head(); i < size && n != nil; i, n = i + 1, n.Next() {
		btreeMap[n.Key] = n.Value
	}
	return btreeMap
}

func (t *Btree) RemoveNode(node *Node) *Btree {
	if node != nil {
		var newNode *Node = node.remove()
		if (newNode != nil && newNode.isRoot()) || node.isRoot() {
			t.Root = newNode
		}
	}
	return t
}

func (t *Btree) Remove(key interface{}) *Btree {
	return t.RemoveNode(t.GetNode(key))
}

func (t *Btree) Pop() *Node {
	var node *Node = t.Tail()
	t.RemoveNode(node)
	return node
}

func (t *Btree) Pull() *Node {
	var node *Node = t.Head()
	t.RemoveNode(node)
	return node
}

func (t *Btree) count(node *Node, nodes *int) {
	if node != nil {
		t.count(node.left, nodes)
		*nodes += 1
		t.count(node.right, nodes)
	}
}

// Node methods
func (n *Node) Init() *Node {
	n.balance = 0
	n.left = nil
	n.right = nil
	n.parent = nil
	return n
}

func (n *Node) Insert(newNode *Node) int {
	c := Comp(newNode.Key, n.Key)
	if c < 0 {
		if n.left == nil {
			n.attachLeftNode(newNode, true)
		} else {
			weight := int(math.Abs(float64(n.left.Insert(newNode))))
			if weight == 0 { return weight }
			n.balance = n.balance - weight
		}
	} else if c > 0 {
		if n.right == nil {
			n.attachRightNode(newNode, true)
		} else {
			weight := int(math.Abs(float64(n.right.Insert(newNode))))
			if weight == 0 { return weight }
			n.balance = n.balance + weight
		}
	} else {
		n.replace(newNode)
		return 0
	}
	if n.balance < -1 { n.rotateRight() }
	if n.balance > 1 { n.rotateLeft() }
	return n.balance
}

func (n *Node) get(k interface{}) *Node {
	var node *Node = nil
	c := Comp(k, n.Key)
	if c < 0 {
		if n.hasLeftChild() { node = n.left.get(k) }
	} else if c > 0 {
		if n.hasRightChild() { node = n.right.get(k) }
	} else {
		node = n
	}
	return node
}

func (n *Node) Next() *Node {
	var next *Node = nil
	if n.hasRightChild() {
		next = n.right.beginning()
	} else if n.hasParent() {
		node := n
		for node.parent != nil {
			node = node.parent
			if Comp(node.Key, n.Key) > 0 {
				next = node
				break
			}
		}
	}
	return next
}

func (n *Node) Prev() *Node {
	var prev *Node = nil
	if n.hasLeftChild() {
		prev = n.left.end()
	} else if n.hasParent() {
		var node *Node = n
		for node.parent != nil {
			node = node.parent
			if Comp(node.Key, n.Key) < 0 {
				prev = node
				break
			}
		}
	}
	return prev
}

func (n *Node) remove() *Node {
	var node *Node = n.parent
	var replace *Node = nil
	if node == nil {
		replace = n.Prev()
		if replace == nil {
			replace = n.Next()
		}
		if replace != nil {
			var replaceParent *Node = replace.parent
			replace.attachLeftNode(n.detachLeftNode(), false)
			replace.attachRightNode(n.detachRightNode(), false)
			replace.nodeType = ROOT
			if replace.right == replace { replace.right = nil }
			if replace.left == replace { replace.left = nil }
			node = replaceParent
		}
	} else if n.hasTwoChildren() {
		node.detachNode(n.nodeType)
		replace = n.Prev()
		var previousParent *Node = replace.parent
		node.attachNode(previousParent.detachNode(replace.nodeType), n.nodeType)
		replace.attachLeftNode(n.detachLeftNode(), false)
		replace.attachRightNode(n.detachRightNode(), false)
		node = previousParent
	} else if n.hasNoChildren() {
		node.detachNode(n.nodeType)
	} else if n.hasLeftChild() {
		node.detachNode(n.nodeType)
		replace = n.detachLeftNode()
		node.attachNode(replace, n.nodeType)
	} else {
		node.detachNode(n.nodeType)
		replace = n.detachRightNode()
		node.attachNode(replace, n.nodeType)
	}
	n.Init()
	var newRoot *Node = n.notifyParents(node)
	if newRoot != nil { replace = newRoot }
	return replace
}

func (n *Node) notifyParents(node *Node) *Node {
	var nodeType NodeType = n.nodeType
	for node != nil {
		if nodeType == LEFT {
			node.balance += 1
		} else if nodeType == RIGHT {
			node.balance -= 1
		}
		if node.balance < -1 {
			node = node.rotateRight()
			if node.nodeType == ROOT {
				return node
			}
		}
		if node.balance > 1 {
			node = node.rotateLeft()
			if node.nodeType == ROOT {
				return node
			}
		}
		nodeType = node.nodeType
		if node.balance == 0 {
			node = node.parent
		} else { node = nil }
	}
	return nil
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

// non-exported methods
func (n *Node) attachLeftNode(node *Node, addWeight bool) *Node {
	n.left = node
	if node != nil {
		node.nodeType = LEFT
		node.parent = n
		if addWeight { n.balance -= 1 }
	}
	return n
}

func (n *Node) attachRightNode(node *Node, addWeight bool) *Node {
	n.right = node
	if node != nil {
		node.nodeType = RIGHT
		node.parent = n
		if addWeight { n.balance += 1 }
	}
	return n
}

func (n *Node) detachRightNode() *Node {
	if n.hasRightChild() {
		node := n.right
		n.right = nil
		node.parent = nil
		return node
	}
	return nil
}

func (n *Node) detachLeftNode() *Node {
	if n.hasLeftChild() {
		node := n.left
		n.left = nil
		node.parent = nil
		return node
	}
	return nil
}

func (n *Node) detachNode(nodeType NodeType) *Node {
	var node *Node = nil
	switch nodeType {
	case LEFT:
		node = n.detachLeftNode()
	case RIGHT:
		node = n.detachRightNode()
	}
	return node
}

func (n *Node) replace(node *Node) {
	n.Value = node.Value
}

func (n *Node) attachNode(node *Node, nodeType NodeType) *Node {
	switch nodeType {
	default:
		break
	case LEFT:
		n.attachLeftNode(node, false)
	case RIGHT:
		n.attachRightNode(node, false)
	}
	return n
}

func (n *Node) rotateLeft() *Node {
	var node *Node = n.detachRightNode()
	if node.leansLeft() {
		node = node.rotateRight()
	}
	if n.isRoot(){
		node.nodeType = ROOT
	} else if n.hasParent() {
		n.parent.attachNode(node, n.nodeType)
	}
	if node.hasLeftChild() {
		n.attachRightNode(node.detachLeftNode(), true)
	}
	node.attachLeftNode(n, true)
	if math.Abs(float64(n.balance)) < 2 { node.balance = n.balance * -1 }
	n.setBalance()
	return node
}

func (n *Node) rotateRight() *Node {
	var node *Node = n.detachLeftNode()
	if node.leansRight() {
		node = node.rotateLeft()
	}
	if n.isRoot() {
		node.nodeType = ROOT
	} else if n.hasParent() {
		n.parent.attachNode(node, n.nodeType)
	}
	if node.hasRightChild() {
		n.attachLeftNode(node.detachRightNode(), true)
	}
	node.attachRightNode(n, true)
	if math.Abs(float64(n.balance)) < 2 { node.balance = n.balance * -1 }
	n.setBalance()
	return node
}

func (n *Node) setBalance() {
	if n.hasNoChildren() || n.hasTwoChildren() {
		n.balance = 0
	} else if n.hasLeftChild() {
		n.balance = -1
	} else { n.balance = 1 }
}

func (n *Node) hasNoChildren() bool {
	return n.right == nil && n.left == nil
}

func (n *Node) hasTwoChildren() bool {
	return n.right != nil && n.left != nil
}

func (n *Node) hasRightChild() bool {
	return n.right != nil
}

func (n *Node) hasLeftChild() bool {
	return n.left != nil
}

func (n *Node) hasParent() bool {
	return n.parent != nil
}

func (n *Node) leansRight() bool {
	return n.balance > 0
}

func (n *Node) leansLeft() bool {
	return n.balance < 0
}

func (n *Node) isRoot() bool {
	return n.nodeType == ROOT
}

func (n *Node) beginning() *Node {
	var beginning *Node = n
	for beginning.left != nil {
		beginning = beginning.left
	}
	return beginning
}

func (n *Node) end() *Node {
	var end *Node = n
	for end.right != nil {
		end = end.right
	}
	return end
}