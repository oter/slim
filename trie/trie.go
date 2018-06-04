package trie

type Node struct {
	Children map[int]*Node
	Branches []int
	Step     uint16
	Value    interface{}
}

type Mode byte

const (
	LT Mode = 4 // less than    0100
	EQ Mode = 2 // equal        0010
	GT Mode = 1 // greater than 0001
)

const leafBranch = -1

func New(keys [][]byte, values []interface{}) (root *Node) {

	root = &Node{Children: make(map[int]*Node), Step: 1}

	for i, key := range keys {

		var node = root
		var j int
		var b byte

		for j, b = range key {
			br := int(b)
			if node.Children[br] == nil {
				break
			}
			node = node.Children[br]
		}

		for _, b = range key[j:] {
			br := int(b)
			n := &Node{Children: make(map[int]*Node), Step: 1}

			node.Children[br] = n
			node.Branches = append(node.Branches, br)
			node = n
		}

		leaf := &Node{Value: values[i]}

		node.Children[leafBranch] = leaf
		node.Branches = append(node.Branches, leafBranch)
	}
	return
}

func (r *Node) Squash() {

	for k, n := range r.Children {
		n.Squash()

		if len(n.Branches) == 1 {
			if n.Branches[0] == leafBranch {
				continue
			}
			child := n.Children[n.Branches[0]]
			child.Step += 1
			r.Children[k] = child
		}
	}
}

func (r *Node) Search(key []byte, mode Mode) (value interface{}) {

	var eqNode = r
	var ltNode *Node
	var gtNode *Node

	for i := 0; ; {

		var br int
		if len(key) == i {
			br = leafBranch
		} else {
			br = int(key[i])
		}

		li, ri := neighborBranches(eqNode.Branches, br)
		if li >= 0 {
			ltNode = eqNode.Children[eqNode.Branches[li]]
		}
		if ri >= 0 {
			gtNode = eqNode.Children[eqNode.Branches[ri]]
		}

		eqNode = eqNode.Children[br]

		if eqNode == nil {
			break
		}

		if br == leafBranch {
			break
		}

		i += int(eqNode.Step)

		if i > len(key) {
			gtNode = eqNode
			eqNode = nil
			break
		}
	}

	if mode&LT == LT && ltNode != nil {
		value = ltNode.rightMost().Value
	}
	if mode&GT == GT && gtNode != nil {
		value = gtNode.leftMost().Value
	}
	if mode&EQ == EQ && eqNode != nil {
		value = eqNode.Value
	}

	return
}

func neighborBranches(branches []int, br int) (ltIndex, rtIndex int) {

	if len(branches) == 0 {
		return
	}

	var i int
	var b int

	for i, b = range branches {
		if b >= br {
			break
		}
	}

	if b == br {
		rtIndex = i + 1
		ltIndex = i - 1

		if rtIndex == len(branches) {
			rtIndex = -1
		}
		return
	}

	if b > br {
		rtIndex = i
		ltIndex = i - 1
		return
	}

	rtIndex = -1
	ltIndex = i

	return
}

func (r *Node) leftMost() *Node {

	node := r
	for {
		if len(node.Branches) == 0 {
			return node
		}

		firstBr := node.Branches[0]
		node = node.Children[firstBr]
	}
}

func (r *Node) rightMost() *Node {

	node := r
	for {
		if len(node.Branches) == 0 {
			return node
		}

		lastBr := node.Branches[len(node.Branches)-1]
		node = node.Children[lastBr]
	}
}