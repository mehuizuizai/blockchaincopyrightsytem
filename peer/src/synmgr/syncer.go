package synmgr

import (
	"fmt"
	"chat"
	"DbDao"
	"math/rand"
	"time"	
	"chat"
	pb "chat/proto"
	"common/utils"
	"strconv"
)

type blocks struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}
type TransRoutinePool struct {
	Queue         chan func() error
	Number        int // the number of thread
	Total         int
	result        chan error
	finshCallBack func()
}

func Initialize() {
	//register fetch blocks,blocks height,txpool request callback function.
	chat.RegisterMsg(pb.Request_FETCH_BLOCKS, blocksSyncerHandler, pb.Response_FETCH_BLOCKS)
	chat.RegisterMsg(pb.Request_FETCH_BLOCKHEIGHT, blocksHeightHandler, pb.Response_FETCH_BLOCKHEIGHT)
    chat.RegisterMsg(pb.Request_FETCH_TXPOOL, txPoolSyncer, pb.Response_FETCH_TXPOOL)
}

//区块同步
var sessionMap map[string]fetchBlocks = make(map[string]fetchBlocks)
func BlocksSyncerHandler(payload []byte) error{
	sessionID := strconv.FormatInt(time.Now().UnixNano(), 16)
	
	blocks := blocks{
		Timestamp:         Timestamp,
		Data:              Data,
		PrevBlockHash:     PrevBlockHash,
		Hash:              Hash,
	}
	sessionMap[sessionID] = blocks


	
	
	m1 = make(map[string]string)
	
	retHeight, err := chat.SendMsg(pb.FETCH_BLOCKHEIGHT,  "192.168.13.82")   //请求区块高度
	if err != nil {
		logger.Error("Send message error")
		return fmt.Errorf("Send message error")
	}
	
	resHeight, ok := args.(pb.BlockHeightRequest)      //断言
	
	var height = GetHeight()  //得到本节点区块高度
	
    var diff int =retHeight-height//得出高度之差diff
	
	args := pb.BlocksRequest{
		
		BlocksNnum:     diff,
	}
	
	retMsg, err := chat.SendMsg(pb.FETCH_BLOCKS, args, "192.168.13.82")  //请求区块
	if err != nil {
		logger.Error("Send message error")
		return fmt.Errorf("Send message error")
	}
	//TODO assert
	resMsg, ok := args.(pb.BlocksRequest)      //断言
	
	for k,v := range payload    //解析payload,读取接收的节点
	  for k2,v2 :=range v
	    m1[k2]=v2  
	
	DbDao.AddBlock(m1[json],diff)   //把这diff个区块添加到本节点的区块链中
	
	return nil

}

func blocksSyncerHandler(args interface{}) (pb.Response_Type, interface{}, error) {
	resMsg, ok := args.(pb.FetchBlocksRequest)
	if !ok {
		logger.Error("assert error...")
		return pb.Response_FETCH_BLOCKS, nil, fmt.Errorf("handle fetch blockds msg error")
	}
	
	
	sessionMap[resMsg.SessionID] = blocks{
		Timestamp:         Timestamp,
		Data:              Data,
		PrevBlockHash:     PrevBlockHash,
		Hash:              Hash,
	}

	return pb.Response_FETCH_BLOCKS, nil, nil
}


func blocksHeightHandler(args interface{}) (pb.Response_Type, interface{}, error) {
	resMsg, ok := args.(pb.BlockHeightRequest)  
	if !ok {
		logger.Error("assert error...")
		return pb.Response_FETCH_BLOCKHEIGHT, nil, fmt.Errorf("handle fetch blocks height msg error")
	}

	return pb.Response_FETCH_BLOCKHEIGHT, nil, nil
}


//同步交易
func TxPoolSyncerHandler(bytes payload) bool{
	
	nodes_tx:=make(map[string]string)
	
	retMsg, err := chat.SendMsg(pb.FETCH_TXPOOL,  "192.168.13.82")  //请求区块
	if err != nil {
		logger.Error("Send message error")
		return fmt.Errorf("Send message error")
	}
	//TODO assert
	resMsg, ok := args.(pb.TxPoolRequest)      //断言
	
	for k,v := range resMsg    //解析payload,读取接收的节点
	  for k2,v2 :=range v
	    nodes_tx[k2]=v2  
	    
    DbDao.AddTxPool(nodes_tx)   //将得到的交易池信息存储到本节点交易池中
	
}

func txPoolSyncerHandler(args interface{}) (pb.Response_Type, interface{}, error) {
	resMsg, ok := args.(pb.TxPoolRequest)  //FETCH_TXPOOL
	if !ok {
		logger.Error("assert error...")
		return pb.Response_FETCH_TXPOOL, nil, fmt.Errorf("handle fetch tx pool msg error")
	}

	return pb.Response_FETCH_TXPOOL, nil, nil
}
