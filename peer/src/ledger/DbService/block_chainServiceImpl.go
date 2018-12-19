package DbService

import (
	"fmt"
	"ledger/DbDao"
	"ledger/DbUtil"
	"strconv"
)

func BlockChainIterators() {
	blockChain := DbDao.NewBlockchain()
	blockChainIterator := blockChain.Iterator()
	for {
		block := blockChainIterator.Next()
		fmt.Printf("prev hash :%x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("hash: %x\n", block.Hash)
		pow := DbUtil.NewProofOfWork(block)
		fmt.Printf("Pow : %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
func GetBlockChainHeight() int {
	blockChain := DbDao.NewBlockchain()
	height := SendBlockChainHeigth(blockChain)
	fmt.Println(height)
	return height
}
func GetMerkleRoot() []byte {

	hashForMerkleRoot := [][]byte{}
	blockchain := DbDao.NewBlockchain()
	blockchainIterator := blockchain.Iterator()
	for {
		block := blockchainIterator.Next()
		hashForMerkleRoot = append(hashForMerkleRoot, block.Hash)
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	merkleNode := NewMerkleTree(hashForMerkleRoot)
	return merkleNode.RootNode.Data
}
