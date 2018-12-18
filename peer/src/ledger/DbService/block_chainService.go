package DbService

import (
	"DbDao"
	"DbUtil"
	"fmt"
	"strconv"
)

func NewBlockChain() {
	DbDao.NewBlockchain()
}
func AddBlockToChain(blockchain *DbDao.Blockchain, data string) error {
	fmt.Println("I com here AddBlock")
	err := blockchain.AddBlock(data)
	return err

}
func BlockChainIterator(blockchain *DbDao.Blockchain) {
	blockchainIterator := blockchain.Iterator()
	for {
		block := blockchainIterator.Next()
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
func SendBlockChainHeigth(blockchain *DbDao.Blockchain) int {
	countBlockHeight := 0
	blockchainIterator := blockchain.Iterator()
	for {
		block := blockchainIterator.Next()
		countBlockHeight += 1
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return countBlockHeight
}
