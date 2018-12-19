package DbDao

import (
	"fmt"
	"ledger/DbUtil"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain2.db"
const blocksBucket = "blocks2"

// tip 尾部的意思，这里是存储最后一个块的hash值 ,存储最后的tip就能推导出整条chain
// 在链的尾端可能会短暂分叉的情况，所以选择tip其实是选择那条链
// db 存储数据库连接
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func NewBlockchain() *Blockchain {
	var tip []byte
	// 打开一个BoltDB文件
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		// 如果数据库中不存在区块链就创建一个，否则直接读取最后一个块的hash值
		if bucket == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := DbUtil.NewGenesisBlock()
			bucket, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = bucket.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = bucket.Put([]byte("1"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = bucket.Get([]byte("1"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	blockchain := Blockchain{tip, db}
	return &blockchain
}

// 加入区块时，需要将区块持久化到数据库中
func (blockchain *Blockchain) AddBlock(data string) error {
	var lastHash []byte
	// 首先获取最后一个块的哈希用于生成新的哈希
	err := blockchain.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		lastHash = bucket.Get([]byte("1"))
		return nil
	})
	if err != nil {
		log.Panic(err)
		return err
	}
	newBlock := DbUtil.NewBlock(data, lastHash)
	err = blockchain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
			return err
		}
		err = bucket.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
			return err
		}
		blockchain.tip = newBlock.Hash
		return nil
	})
	return err

}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	blockchainiterator := &BlockchainIterator{blockchain.tip, blockchain.db}
	return blockchainiterator
}

// 返回链中的下一个块
func (i *BlockchainIterator) Next() *DbUtil.Block {
	var block *DbUtil.Block
	err := i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(i.currentHash)
		block = DbUtil.DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.currentHash = block.PrevBlockHash
	return block
}
