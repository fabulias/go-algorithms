package huffman

import (
	"container/heap"
)

type TreeNode struct {
	Char  rune
	Freq  int
	Left  *TreeNode
	Right *TreeNode
}

type PriorityQueue []*TreeNode

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Freq < pq[j].Freq }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x any)        { *pq = append(*pq, x.(*TreeNode)) }
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func compress(text string) error {
	if text == "" {
		return ErrTextEmpty
	}

	frequencies := make(map[rune]int)
	for _, char := range text {
		frequencies[char]++
	}

	pq := make(PriorityQueue, 0)
	for char, freq := range frequencies {
		heap.Push(&pq, &TreeNode{
			Char: char,
			Freq: freq,
		})
	}

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*TreeNode)
		right := heap.Pop(&pq).(*TreeNode)

		root := &TreeNode{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}
		heap.Push(&pq, root)
	}

	root := heap.Pop(&pq).(*TreeNode)
	var huffmanCodes = make(map[rune]string)
	generateCodes(root, "", huffmanCodes)
	writeBits(text, huffmanCodes)
	return nil
}

func generateCodes(node *TreeNode, currentPath string, codes map[rune]string) {
	if node == nil {
		return
	}
	if node.Left == nil && node.Right == nil {
		codes[node.Char] = currentPath
		return
	}

	generateCodes(node.Left, currentPath+"0", codes)
	generateCodes(node.Right, currentPath+"1", codes)
}

func writeBits(text string, codes map[rune]string) []byte {
	var compressed []byte
	var currentByte byte
	var bitCount uint8

	for _, char := range text {
		code := codes[char]

		for _, bit := range code {
			currentByte = (currentByte << 1) | byte(bit-'0')
			bitCount++

			if bitCount == 8 {
				compressed = append(compressed, currentByte)
				currentByte, bitCount = 0, 0
			}
		}
	}
	if bitCount > 0 {
		currentByte <<= (8 - bitCount)
		compressed = append(compressed, currentByte)
	}
	return compressed
}
