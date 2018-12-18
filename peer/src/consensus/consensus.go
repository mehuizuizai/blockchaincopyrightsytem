package consensus

import (
	"bytes"
	"chat"
	pb "chat/proto"
	"common/utils"
	"config"
	"fmt"
	"logging"
	"math"
	"strings"
	"time"
)

var logger = logging.MustGetLogger()

var conSessionMap map[string]map[string][]byte = make(map[string]map[string][]byte)

//type oneSession struct {
//	votes map[string][]byte
//	time.Timer
//}

func Initialize() {
	//register consensus request callback function.
	chat.RegisterMsg(pb.Request_CONSENSUS, consensusCallback, pb.Response_CONSENSUS)
}

func StartConsensus(selfVote []byte, sessionID string) (bool, bool, []string) {
	//put myself vote into session map.
	_, ok := conSessionMap[sessionID]
	if !ok {
		conSessionMap[sessionID] = make(map[string][]byte)
	}
	selfIP := utils.GetlocalIP()
	conSessionMap[sessionID][selfIP] = selfVote

	//broadcast my vote to other peer.
	broadcastMyVote(selfVote, selfIP, sessionID)

	//do consensus, and get result(whether equal? and majority peers list)
	isSuccessful, mostVote, mostPeers := doVote(sessionID)
	return isSuccessful, bytes.Equal(mostVote, selfVote), mostPeers

}

func consensusCallback(args interface{}) (pb.Response_Type, interface{}, error) {
	reqMsg, ok := args.(pb.ConsensusRequest)
	if !ok {
		logger.Error("assert error...")
		return pb.Response_CONSENSUS, pb.ConsensusResponse{}, fmt.Errorf("handle consensus msg error")
	}
	//put consensus session content into session map.
	_, ok = conSessionMap[reqMsg.SessionID]
	if !ok {
		conSessionMap[reqMsg.SessionID] = make(map[string][]byte)
	}
	conSessionMap[reqMsg.SessionID][reqMsg.IP] = reqMsg.Vote

	return pb.Response_CONSENSUS, pb.ConsensusResponse{}, nil
}

func broadcastMyVote(selfVote []byte, selfIP, sessionID string) {
	//get all peer.
	peers := config.GetPeers()
	var indexOfMe int
	for key, value := range peers {
		if strings.EqualFold(value, selfIP) {
			indexOfMe = key
			break
		}
	}
	peersNotMe := append(peers[:indexOfMe], peers[indexOfMe+1:]...)

	//loop send msg
	args := &pb.ConsensusRequest{
		SessionID: sessionID,
		IP:        selfIP,
		Vote:      selfVote,
	}
	for _, value := range peersNotMe {
		chat.SendMsg(pb.Request_CONSENSUS, args, value)
	}
}

func doVote(sessionID string) (isSuccessful bool, mostVote []byte, mostPeers []string) {
	var peersVotes map[string][]byte
	allPeers := config.GetPeers()
	for i := 5; i > 0; i-- {
		peersVotes = conSessionMap[sessionID]
		if len((peersVotes)) == len(allPeers) {
			break
		}
		time.Sleep(time.Second * 1)
	}

	//if mostPeers num is bigger than or equal to votingThreshold, then consensus successfully.
	var votingThreshold int
	if len(allPeers) == 1 {
		votingThreshold = 1
	} else {
		votingThreshold = int(math.Ceil(float64(len(allPeers))/float64(2))) + 1
	}

	if len(peersVotes) < votingThreshold {
		logger.Warning("peers join in voting is too few, and consensus failed")
		return
	}

	var voteSlice [][]byte
	for _, value := range peersVotes {
		voteSlice = append(voteSlice, value)
	}
	//get
	maxIndex := 0
	maxNum := 1
	for i := 0; i < len(voteSlice); i++ {
		voteNum := 1
		for j := i + 1; j < len(voteSlice); j++ {
			if bytes.Equal(voteSlice[i], voteSlice[j]) {
				voteNum++
			}
		}
		if voteNum > maxNum {
			maxNum = voteNum
			maxIndex = i
		}
	}

	if maxNum < votingThreshold {
		logger.Warning("most peers votes is too few, and consensus failed")
		return
	}

	mostVote = voteSlice[maxIndex]

	for key, value := range peersVotes {
		if bytes.Equal(value, mostVote) {
			mostPeers = append(mostPeers, key)
		}
	}

	isSuccessful = true

	return
}
