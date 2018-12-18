package DbService

import (
	"crypto/sha256"
	"fmt"
)

type MerkleTree struct {
	RootNode *MerkleNode
}
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right

	return &mNode
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	for _, dataitem := range data {
		node := NewMerkleNode(nil, nil, dataitem)
		nodes = append(nodes, *node)
	}
	for i := 0; i < len(data)/2; i++ {
		var newNodes []MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			if len(nodes)-j < 2 {
				node := NewMerkleNode(&nodes[j], &nodes[j], nil)
				newNodes = append(newNodes, *node)
			} else {
				node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
				newNodes = append(newNodes, *node)
			}

		}
		nodes = newNodes
	}
	mTree := MerkleTree{&nodes[0]}
	fmt.Println("here mTree", mTree.RootNode.Data)
	return &mTree
}
