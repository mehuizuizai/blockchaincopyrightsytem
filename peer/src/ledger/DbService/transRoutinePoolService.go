package DbService

import (
	"DbDao"
	"DbUtil"
	"fmt"

	"time"
)

//work_entry ,trancopyrightx
var flag_pool_Init = false
var flag_create_blockchain = false
var TransList chan func() error
var pool *DbUtil.TransRoutinePool = new(DbUtil.TransRoutinePool)
var blockchain *DbDao.Blockchain = new(DbDao.Blockchain)
var txHashPool = make([]string, 0)

func Tx_Enter_routinePool(getTxHash string) {
	if flag_pool_Init == false {
		pool.Init(1, 0)
		fmt.Println("here , flag_pool_Init is False")
	}
	flag_pool_Init = true
	pool.Add_transaction()
	fmt.Println("pool total:", pool.Total)
	if flag_create_blockchain == false {
		blockchain = DbDao.NewBlockchain()
	}
	flag_create_blockchain = true
	fmt.Println("blockchain", blockchain)
	pool.AddTask(func() error { // make  the task ->quene
		return AddBlockToChain(blockchain, getTxHash) // the service method)
	})
	/*time.Sleep(time.Millisecond * 10000)*/ //time to wait for AddBlock
	fmt.Println("sleep over")
	beginTaskForAddBlock()
}

func beginTaskForAddBlock() {
	isFinish := false
	pool.SetFinishCallback(func() {
		func(isFinish *bool) {
			*isFinish = true
		}(&isFinish)
	})
	pool.Start()
	index := 0
	//	fmt.Println("here ,txHashPool is :", txHashPool)
	//---------------------------
	postRead := []string{}
	DbUtil.Load(&postRead, "txhash")
	fmt.Println("txhash", postRead)
	postRead = append(postRead[:index], postRead[index+1:]...)
	DbUtil.Store(postRead, "txhash")

	//--------------------
	for !isFinish {
		time.Sleep(time.Millisecond * 100)
	}
	//	pool.Stop()
	fmt.Println("AddBlock success")
}
