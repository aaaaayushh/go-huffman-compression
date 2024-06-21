package huff

import "container/heap"

type BaseNode interface {
	IsLeaf() bool
	Weight() int
}
type LeafNode struct {
	weight  int
	element rune
}

func (node LeafNode) IsLeaf() bool {
	return true
}
func (node LeafNode) Weight() int {
	return node.weight
}
func (node LeafNode) Value() rune {
	return node.element
}
func NewHuffLeafNode(el rune, w int) *LeafNode {
	return &LeafNode{element: el, weight: w}
}

type InternalNode struct {
	weight      int
	left, right BaseNode
	LeftEdge    int
	RightEdge   int
}

func NewHuffInternalNode(l, r BaseNode, w int) *InternalNode {
	return &InternalNode{left: l, right: r, weight: w, LeftEdge: 0, RightEdge: 1}
}
func (node InternalNode) IsLeaf() bool {
	return false
}
func (node InternalNode) Weight() int {
	return node.weight
}
func (node InternalNode) Left() BaseNode {
	return node.left
}
func (node InternalNode) Right() BaseNode {
	return node.right
}

type Tree struct {
	root BaseNode
}

func NewHuffTreeFromLeaf(r BaseNode) *Tree {
	return &Tree{root: r}
}
func NewHuffTreeFromNodes(l, r BaseNode, wt int) *Tree {
	return &Tree{root: NewHuffInternalNode(l, r, wt)}
}
func (tree *Tree) Root() BaseNode {
	return tree.root
}
func (tree *Tree) Weight() int {
	return tree.root.Weight()
}
func (tree *Tree) CompareTo(other *Tree) int {
	if tree.Weight() < other.Weight() {
		return -1
	} else if tree.Weight() == other.Weight() {
		return 0
	} else {
		return 1
	}
}

// HuffmanHeap is a min-heap of HuffTree pointers
type HuffmanHeap []*Tree

func (h *HuffmanHeap) Len() int           { return len(*h) }
func (h *HuffmanHeap) Less(i, j int) bool { return (*h)[i].Weight() < (*h)[j].Weight() }
func (h *HuffmanHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *HuffmanHeap) Push(x interface{}) {
	*h = append(*h, x.(*Tree))
}

func (h *HuffmanHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// BuildHuffmanTree constructs a Huffman tree from a map of character frequencies
func BuildHuffmanTree(freqsMap map[rune]int) *Tree {
	// Create a min-heap
	h := &HuffmanHeap{}
	heap.Init(h)

	// Create a leaf node for each character and add it to the heap
	for ch, freq := range freqsMap {
		huffNode := NewHuffLeafNode(ch, freq)
		tree := NewHuffTreeFromLeaf(huffNode)
		heap.Push(h, tree)
	}

	// While there is more than one tree in the heap
	for h.Len() > 1 {
		// Remove the two trees with the lowest weight
		tree1 := heap.Pop(h).(*Tree)
		tree2 := heap.Pop(h).(*Tree)

		// Create a new internal node with these two nodes as children
		combinedWeight := tree1.Weight() + tree2.Weight()
		newTree := NewHuffTreeFromNodes(tree1.Root(), tree2.Root(), combinedWeight)

		// Add the new tree back to the heap
		heap.Push(h, newTree)
	}

	// The last remaining tree is the Huffman tree
	return heap.Pop(h).(*Tree)
}
