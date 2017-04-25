package tree

import "math"

type Btree struct {
	root *Node
	comparer Comparer
}

type Traverser interface {
	Next() *Node
	Previous() *Node
}

type Comparer interface {
	Comp(v1, v2 interface{}) int
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

func (n *Node) Init() *Node {
	n.balance = 0
	n.left = nil
	n.right = nil
	n.parent = nil
	n.nodeType = ROOT
	return n
}

func New() *Btree { return new(Btree).Init() }

func (t *Btree) Init() *Btree {
	t.root = nil
	return t
}

func (t *Btree) isEmpty() bool {
	return t.root == nil
}

func (t *Btree) Put(key, value interface{}) *Btree {
	var newNode *Node = (&Node{Key: key, Value: value}).Init()
	if t.isEmpty() {
		t.root = newNode
	} else {
		t.root.Insert(newNode)
		for t.root.nodeType != ROOT {
			t.root = t.root.parent
		}
	}
	return t
}

func (t *Btree) Insert(value interface{}) *Btree {
	return t.Put(value, value)
}

func (t *Btree) GetNode(key interface{}) *Node {
	var node *Node = nil
	if !t.isEmpty() { node = t.root.get(key) }
	return node
}

func (t *Btree) Get(key interface{}) interface{} {
	var node *Node = t.GetNode(key)
	if node != nil {
		return node.Value
	}
	return nil
}

func (t *Btree) Size() int {
	nodes := 0
	if !t.isEmpty() {
		t.count(t.root, &nodes)
	}
	return nodes
}

func (t *Btree) count(node *Node, nodes *int) {
	if node != nil {
		t.count(node.left, nodes)
		*nodes += 1
		t.count(node.right, nodes)
	}
}

func (t *Btree) beginning() *Node {
	if t.isEmpty() { return nil }
	return t.root.beginning()
}

func (t *Btree) end() *Node {
	if t.isEmpty() { return nil }
	return t.root.end()
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
		for n.parent != nil {
			node = n.parent
			if Comp(node.Value, n.Value) > 0 {
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
		prev = prev.end()
	} else if n.hasParent() {
		var node *Node = n
		for node.parent != nil {
			node = node.parent
			if Comp(node.Value, n.Value) < 0 {
				prev = node
				break
			}
		}
	}
	return prev
}

func Comp(v1, v2 interface{}) int  {
	c := 0

	switch v1.(type) {
	default:
		if v1.(int) > v2.(int) { c = 1 } else if v1.(int) < v2.(int) { c = -1 } else { c = 0 }
	case int:
		if v1.(int) > v2.(int) { c = 1 } else if v1.(int) < v2.(int) { c = -1 } else { c = 0 }
	case uint:
		if v1.(uint) > v2.(uint) { c = 1 } else if v1.(uint) < v2.(uint) { c = -1 } else { c = 0 }
	case float32:
		if v1.(float32) > v2.(float32) { c = 1 } else if v1.(float32) < v2.(float32) { c = -1 } else { c = 0 }
	case float64:
		if v1.(float64) > v2.(float64) { c = 1 } else if v1.(float64) < v2.(float64) { c = -1 } else { c = 0 }
	case uintptr:
		if v1.(uintptr) > v2.(uintptr) { c = 1 } else if v1.(uintptr) < v2.(uintptr) { c = -1 } else { c = 0 }
	case byte:
		if v1.(byte) > v2.(byte) { c = 1 } else if v1.(byte) < v2.(byte) { c = -1 } else { c = 0 }
	case rune:
		if v1.(rune) > v2.(rune) { c = 1 } else if v1.(rune) < v2.(rune) { c = -1 } else { c = 0 }
	case complex64:
		if v1.(float64) > v2.(float64) { c = 1 } else if v1.(float64) < v2.(float64) { c = -1 } else { c = 0 }
	case complex128:
		if v1.(float64) > v2.(float64) { c = 1 } else if v1.(float64) < v2.(float64) { c = -1 } else { c = 0 }
	case string:
		if v1.(string) > v2.(string) { c = 1 } else if v1.(string) < v2.(string) { c = -1 } else { c = 0 }
	case Comparer:
		c = v1.(Comparer).Comp(v1, v2)
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
	var node *Node = nil
	if n.hasRightChild() {
		node = n.right
		n.right = nil
		node.parent = nil
	}
	return node
}

func (n *Node) detachLeftNode() *Node {
	var node *Node = nil
	if n.hasLeftChild() {
		node = n.left
		n.left = nil
		node.parent = nil
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
		n.nodeType = ROOT
	} else if n.hasParent() {
		n.parent.attachNode(node, n.nodeType)
	}
	if node.hasRightChild() {
		n.attachLeftNode(node.detachRightNode(), true)
	}
	node.attachRightNode(n, true)
	if math.Abs(float64(n.balance)) < 2 { node.balance = node.balance * -1 }
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