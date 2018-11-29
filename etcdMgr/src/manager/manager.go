package manager

import (
	"bytes"
	"config"
	"db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logging"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"utils"

	"github.com/coreos/etcd/etcdmain"
)

var clusterMembers []db.MemberInfo //store the info of members in etcd cluster

var logger = logging.MustGetLogger()

var selfIP string
var selfClientPort string
var selfPeerPort string
var slowheartbeatCycle = 10
var broadcastSlowHeartbeatController = true
var identity string

func Initialize() error {
	//check the local ip is as same as the ip from cert
	selfIP = config.GetLocalHostIP()
	var err error
	if selfIP == "" {
		selfIP, err = utils.GetlocalIP()
		if err != nil {
			logger.Error(err.Error())
		}
	}
	selfClientPort = config.GetEtcdClientPort()
	selfPeerPort = config.GetEtcdPeerPort()
	identity = config.GetEtcdIdentity()
	//Load data related to ETCD
	if strings.EqualFold("creator", identity) {
		ok := loadDataForAdmin()
		if !ok {
			logger.Error("loadData failed!")
		}
		//start ETCD
		err := operationOfETCDForAdmin()
		if err == nil {
			startSlowHeartbeat(identity)
		} else {
			logger.Error(err.Error())
		}
	} else {
		ok := loadDataForPeer()
		if !ok {
			logger.Error("loadData failed!")
		}
		err := operationOfETCDForPeer()
		if err == nil {
			startSlowHeartbeat(identity)
		} else {
			logger.Error(err.Error())
		}
	}

	updateMembersInfo()

	communicationWithUcc()

	startMonitorClusterApp()

	startMonitorNetApp()
	return nil
}

func loadDataForAdmin() bool {

	//Get nodes info from db
	members, err := db.MembersInfoQuery()
	if err != nil {
		logger.Error("Get members info failed from db")
		return false
	}
	if len(members) == 0 { //db is null,then get creator's info
		var member db.MemberInfo
		member.IP = selfIP
		member.ClientPort = selfClientPort
		member.PeerPort = selfPeerPort

		ok := db.MemberInfoInsert(member)
		if !ok {
			logger.Error("Insert creator's info into db failed!")
			return false
		}
		clusterMembers = append(clusterMembers, member)
	} else {
		clusterMembers = members
	}
	return true

}

func loadDataForPeer() bool {
	//Get nodes info from db
	members, err := db.MembersInfoQuery()
	if err != nil {
		logger.Error("Get members info failed from db")
		return false
	}
	if len(members) == 0 { //db is null,then get creator's info
		//Get cluster members from config file.
		membersStr := config.GetEtcdCluterMembers()
		membersInfo := strings.Split(membersStr, ",")
		var peer db.MemberInfo
		for _, value := range membersInfo {
			host := strings.Split(value, ":")
			peer.IP = host[0]
			peer.ClientPort = host[1]
			clusterMembers = append(clusterMembers, peer)
		}
	} else {
		clusterMembers = members
	}
	return true

}

func operationOfETCDForAdmin() error {
	if len(clusterMembers) == 1 && clusterMembers[0].IP == selfIP && clusterMembers[0].ClientPort == selfClientPort {
		ok := activateForAdmin("new")
		if !ok {
			return fmt.Errorf("Create cluster failed!")
		}
	} else {
		if isClusterAlive() {
			if isInCluster() {
				ttl := getSlowHeartBeat()
				if ttl < int64(2*slowheartbeatCycle) {
					time.Sleep(time.Second * time.Duration(ttl+5))
					for getSlowHeartBeat() <= 0 {
						if !isInCluster() {
							ok := isJoinCluster()
							if ok {
								err := register()
								if err != nil {
									return fmt.Errorf("join cluster failed")
								}

								ok = activateForAdmin("join")
								if !ok {
									return fmt.Errorf("start ETCD failed")
								}
							} else {
								return fmt.Errorf("can not join cluster")
							}
						}
						break
					}
				}
				ok := activateForAdmin("existing")
				if !ok {
					return fmt.Errorf("Start ETCD failed")
				}
			} else {
				ok := isJoinCluster()
				if ok {
					err := register()
					if err != nil {
						return fmt.Errorf("join cluster failed")
					}
					ok = activateForAdmin("join")
					if !ok {
						return fmt.Errorf("Start ETCD failed")
					}
				} else {
					return fmt.Errorf("can not join cluster")
				}
			}
		} else {
			ok := activateForAdmin("new")
			if !ok {
				return fmt.Errorf("Start ETCD failed")
			}
		}
	}

	return nil
}

func operationOfETCDForPeer() error {
	//加入集群
	if isInCluster() {
		ttl := getSlowHeartBeat()
		if ttl < int64(2*slowheartbeatCycle) {
			time.Sleep(time.Second * time.Duration(ttl+5))
			for getSlowHeartBeat() <= 0 {
				if !isInCluster() {
					for {
						ok := isJoinCluster()
						if ok {
							er := register()
							if er != nil {
								logger.Warning(er.Error())
								time.Sleep(time.Second * 2)
								continue
							}
							ok = activateForPeer()
							if !ok {
								return fmt.Errorf("start ETCD failed")
							}
							break
						} else {
							return fmt.Errorf("can not join cluster")
						}
					}
				} else {
					ok := activateForPeer()
					if !ok {
						return fmt.Errorf("start ETCD failed")
					}
				}
				break
			}
		} else {
			ok := activateForPeer()
			if !ok {
				return fmt.Errorf("start ETCD failed")
			}
		}
	} else {
		for {
			ok := isJoinCluster()
			if ok {
				er := register()
				if er != nil {
					logger.Warning(er.Error())
					time.Sleep(time.Second * 2)
					continue
				}
				ok = activateForPeer()
				if !ok {
					return fmt.Errorf("start ETCD failed")
				}
				break
			} else {
				return fmt.Errorf("can not join cluster")
			}
		}

	}
	return nil
}

func isJoinCluster() bool {
	url := ""
	clusterCapacity, _ := strconv.Atoi(config.GetEtcdClusterCapacity())
	/*轮询clusterMembers，直到可以正常访问ETCD为止*/
	for _, value := range clusterMembers {
		url = "members"
		body, err := HttpGet(value.IP, value.ClientPort, url)
		if err != nil {
			if strings.EqualFold("key not found", err.Error()) {
				logger.Error("key not found")
				return false
			}
			continue
		}

		d := struct {
			Members []Member
		}{}
		if err = json.Unmarshal(body, &d); err != nil {
			logger.Error(err.Error())
			continue
		}

		if len(d.Members) >= clusterCapacity {
			logger.Warning("cluster is full")
			return false
		}

		return true
	}
	logger.Error("cluster is not accessable")
	return false
}

func updateMembersInfo() {
	go func() {
		for {
			//update cache
			peers := getMembers()
			if len(peers) != 0 {
				//check if has new members is in cluster
				for _, value1 := range peers {
					var flag = false
					for _, value2 := range clusterMembers {
						if value1.IP == value2.IP && value1.ClientPort == value2.ClientPort {
							flag = true
						}
					}
					if !flag {
						ok := db.MemberInfoInsert(value1)
						if !ok {
							logger.Error("Insert new member info into db failed!")
						}
					}
				}
				//check if has old members is not in cluster
				for _, value1 := range clusterMembers {
					var flag = false
					for _, value2 := range peers {
						if value1.IP == value2.IP && value1.ClientPort == value2.ClientPort {
							flag = true
						}
					}
					if !flag {
						ok := db.MemberInfoDelete(value1)
						if !ok {
							logger.Error("Insert new member info into db failed!")
						}
					}
				}
				clusterMembers = peers
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

func communicationWithUcc() {
	go func() {
		//if ucc folder that store socket file is existing?
		//Not: create.
		//		isExsting, err := isSocketFileExisting(socketFolderPath)
		//		if err != nil {
		//			logger.Error(err.Error())
		//			return
		//		} else if !isExsting {
		//			err = os.Mkdir(socketFolderPath, os.ModePerm)
		//			if err != nil {
		//				logger.Error(err.Error())
		//				return
		//			}
		//		}
		var socketFilePath = config.GetEtcdMgrPath()
		isExsting, err := isSocketFileExisting(socketFilePath)
		if err != nil {
			logger.Error(err.Error())
			return
		} else if isExsting {
			err = os.Remove(socketFilePath)
			if err != nil {
				logger.Error(err.Error())
				return
			}
		}

		//建立socket，监听端口
		unixAddr, er := net.ResolveUnixAddr("unix", socketFilePath)
		if er != nil {
			logger.Error(er.Error())
		}
		unixListener, err := net.ListenUnix("unix", unixAddr)
		if err != nil {
			logger.Error(err.Error())
		}
		defer unixListener.Close()

		//fmt.Println("Waiting for clients")
		for {
			unixConn, err := unixListener.AcceptUnix()
			if err != nil {
				continue
			}

			//logger.Info(conn.RemoteAddr().String(), " tcp connect success")
			go handleConnection(unixConn)
		}
	}()
}

func isSocketFileExisting(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//处理连接
func handleConnection(conn *net.UnixConn) {

	buffer := make([]byte, 2048)

	n, err := conn.Read(buffer)

	if err != nil {
		if !strings.EqualFold(err.Error(), "EOF") {
			logger.Error(conn.RemoteAddr().String(), " connection error: ", err)
		}
		return
	}
	if strings.EqualFold(string(buffer[:n]), "GET") {
		//send members-list of cluster to creator
		conn.Write([]byte(parseMembersListToString()))

	}

}

func parseMembersListToString() string {
	var str string
	for _, value := range clusterMembers {
		str = str + value.IP + ":"
		str = str + value.ClientPort + "-"
	}
	return str
}

func isInCluster() bool {
	peers := getMembers()
	for _, value := range peers {
		if value.IP == selfIP && value.ClientPort == selfClientPort {
			return true
		}
	}

	os.RemoveAll("Creator" + selfIP + ".etcd")
	return false
}

func getMembers() []db.MemberInfo {
	url := "members"
	var peers []db.MemberInfo
	/*轮询clusterMembers，直到可以正常访问ETCD为止*/
	for _, value := range clusterMembers {
		body, err := HttpGet(value.IP, value.ClientPort, url)
		if err != nil {
			//logger.Error(err.Error())
			continue
		}

		d := struct {
			Members []Member
		}{}
		if err = json.Unmarshal(body, &d); err != nil {
			logger.Error(err.Error())
			continue
		}
		var peer db.MemberInfo
		for _, value := range d.Members {
			if len(value.ClientURLs) > 0 {
				//length := len(value.ClientURLs)
				IP := value.ClientURLs[0]
				IP2 := value.PeerURLs[0]
				host := IP[7:len(value.ClientURLs[0])] //ip+port
				host2 := IP2[7:len(value.PeerURLs[0])]
				idx := strings.Index(host, ":")
				idx2 := strings.Index(host2, ":")
				peer.IP = host[:idx]
				peer.ClientPort = host[idx+1:]
				peer.PeerPort = host2[idx2+1:]
				peer.Name = value.Name
				peer.ID = value.ID
				peers = append(peers, peer)
			} else {
				peer.IP = selfIP
				peer.ClientPort = selfClientPort
				peer.PeerPort = selfPeerPort
				peer.ID = value.ID
				peer.Name = identity + selfIP
				peers = append(peers, peer)
			}

		}

		break
	}
	return peers
}

func startMonitorNetApp() {
	var preStatus = 1
	go func() {
		url := "members"
		for {
			_, err := HttpGet(selfIP, selfClientPort, url)
			if err != nil {
				if strings.Contains(err.Error(), "network is unreachable") { //the net is unaccessble,put 0 into queue.
					preStatus = 0
					logger.Error("The net is not reachable")
				} else { //the net is accessble,but my ETCD is not working.
					if preStatus == 0 {
						operationForNetStatusChange()

					}
					preStatus = 1
				}
			} else { //the net is unaccessble,put 1 into queue.
				if preStatus == 0 {
					operationForNetStatusChange()

				}
				preStatus = 1
			}

			time.Sleep(time.Second * 5)
		}

	}()
}

func operationForNetStatusChange() {
	ttl := getSlowHeartBeat()
	if ttl < int64(2*slowheartbeatCycle) {
		broadcastSlowHeartbeatController = false
		time.Sleep(time.Second * time.Duration(ttl+5))
		for getSlowHeartBeat() <= 0 {
			if !isInCluster() {
				err := register()
				if err != nil {
					logger.Error("Join cluster falied!")
					return
				}
				var ok bool
				if strings.EqualFold("creator", identity) {
					ok = activateForAdmin("join")
				} else {
					ok = activateForPeer()
				}

				if !ok {
					logger.Error("Start ETCD failed!")
					return
				}
			}
			break
		}
		broadcastSlowHeartbeatController = true
	}
}

func startMonitorClusterApp() {
	go func() {
		for {
			IdSet, err := getIdSetOfDiedPeer()
			if err != nil {
				//logger.Error(err.Error())
				continue
			}
			if IdSet != nil {
				for _, value := range IdSet {
					is := false
					is, err = isLeader()
					if err != nil {
						logger.Error(err.Error())
					}
					if is {
						ok := cleanPeerFromCluster(value)
						if !ok {
							logger.Error("Clean peer failed")
						}
					}
				}
			}

			time.Sleep(time.Second * 10)
		}
	}()
}

func getIdSetOfDiedPeer() ([]string, error) {
	url := "keys/nodeStatusOfSlowheartbeat?recursive=true"
	var m Event
	var heartbeatSet []string
	for _, value := range clusterMembers {
		body, err := HttpGet(value.IP, value.ClientPort, url)
		if err != nil {
			//logger.Error(err.Error())
			continue
		}
		err = json.Unmarshal(body, &m)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		break
	}
	//fmt.Println("cluster---", clusterMembers)
	//fmt.Println("m.Node.Nodes------", m.Node.Nodes)
	if m.Node != nil {
		for _, value := range clusterMembers {
			flag := false
			for _, value2 := range m.Node.Nodes {
				name := value2.Key[27:]
				if strings.EqualFold(value.Name, name) {
					flag = true
					break
				}
			}
			if flag == false {
				heartbeatSet = append(heartbeatSet, value.ID)
			}
		}

	} else {
		return nil, fmt.Errorf("acluster is unaccessable!")
	}

	return heartbeatSet, nil

}

func peerWatcher() (string, error) {
	url := "keys/nodeStatusOfSlowheartbeat?recursive=true&wait=true"
	var m Event
	for _, value := range clusterMembers {
		body, err := HttpGet(value.IP, value.ClientPort, url)
		if err != nil {
			//logger.Error(err.Error())
			continue
		}
		err = json.Unmarshal(body, &m)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		break
	}
	if m.Action == "" {
		return "", fmt.Errorf("cluster is not accessable!")
	} else if m.Action == "expire" {
		IP := m.PrevNode.Key[28:]
		return IP, nil
	} else {
		return "", nil
	}

}

func isLeader() (bool, error) {
	leaderIP, leaderPeerPort, err := getLeaderIP()
	if err != nil {
		return false, err
	}
	if leaderIP == selfIP && leaderPeerPort == selfPeerPort {
		return true, nil
	}

	return false, nil

}

func cleanPeerFromCluster(id string) bool {

	url := "members/" + id
	err := HttpDelete(selfIP, selfClientPort, url)
	if err != nil {
		logger.Error("Clean peer failed")
		return false
	}

	return true
}

func register() error {

	/*for循环遍历clusterMembers进行注册，一旦注册成功，就执行清理ETCD相关文件的操作——cleanETCDFile()，然后返回true，若遍历完clusterMembers中所有的节点还是未注册成功的话则返回false*/
	b, err := json.Marshal(request{PeerURLs: []string{"http://" + selfIP + ":" + selfPeerPort}})
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	var numberOfPeers = 0 //记录已经向多少个节点进行了注册操作
	peers := getMembers()
	for _, value := range peers {
		numberOfPeers++
		logger.Info("Send register namerequest...")
		resp, err := http.Post("http://"+value.IP+":"+value.ClientPort+"/v2/members", "application/json", bytes.NewReader(b))
		if err != nil {
			//logger.Error("registert--------" + err.Error())
			continue
		}
		logger.Info("Receive register response")
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error("registert--------" + err.Error())
			resp.Body.Close()
			continue
		}
		logger.Info("register response", string(body))
		if strings.Contains(string(body), "Error") {
			//logger.Error(string(body))
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		numberOfPeers--
		break

	}
	if numberOfPeers == len(peers) {
		return fmt.Errorf("cluster is not accessble!")
	}

	cleanETCDFile()

	return nil

}

func cleanETCDFile() {
	//启动ETCD前，检查目录中是否前一次启动ETCD产生的相关文件，如果有则进行清理操作
	exist_etcd := checkFileIsExist("./" + identity + selfIP + ".etcd")
	if exist_etcd {
		os.RemoveAll("./" + identity + selfIP + ".etcd")
	}
}

func activateForPeer() bool {

	peers := getMembers()

	//str := "#!/bin/bash\n\n"
	var str string = ""
	str = str + config.BasePath + "/etcd --name " + "peer" + selfIP + " --data-dir " + config.BasePath + "/" + "peer" + selfIP + ".etcd" + " --initial-advertise-peer-urls http://" + selfIP + ":" + selfPeerPort
	str = str + " --listen-peer-urls http://" + selfIP + ":" + selfPeerPort
	str = str + " --listen-client-urls http://" + selfIP + ":" + selfClientPort + ",http://127.0.0.1:" + selfClientPort
	str = str + " --advertise-client-urls http://" + selfIP + ":" + selfClientPort
	str = str + " --initial-cluster "
	for _, value := range peers {
		str = str + value.Name + "=http://" + value.IP + ":" + value.PeerPort + ","
	}
	str = str + " --initial-cluster-state existing"

	//Split configure info (string) to []string
	var args []string = make([]string, 0)
	strs := strings.Split(str, " ")
	for i := 0; i < len(strs); i++ {
		args = append(args, strs[i])
	}
	//Put configure info ([]string) to os.Args
	var len_Args int = 0
	for key, _ := range os.Args {
		os.Args[key] = args[key]
		len_Args = len_Args + 1
	}
	left_args := args[len_Args:]
	os.Args = append(os.Args, left_args...)

	go etcdmain.Main()

	//Send msg to check etcd is running
	timeout := make(chan bool, 1)
	ch := make(chan bool, 1)
	go func() {
		time.Sleep(time.Second * 5)
		timeout <- true
	}()

	go func() {
		for {
			if etcdDone() {
				ch <- true
				break
			}
		}
	}()

	select {
	case <-ch:
		fmt.Println("start ETCD..........")
	case <-timeout:
		logger.Warning("ETCD never response")
		return false
	}
	return true
}

func activateForAdmin(flag string) bool {

	var file myfile
	var str string = ""
	str = str + config.BasePath + "etcd --name " + "creator" + selfIP + " --data-dir " + config.BasePath + "/" + "creator" + selfIP + ".etcd" + " --initial-advertise-peer-urls http://" + selfIP + ":" + selfPeerPort
	str = str + " --listen-peer-urls http://" + selfIP + ":" + selfPeerPort
	str = str + " --listen-client-urls http://" + selfIP + ":" + selfClientPort + ",http://127.0.0.1:" + selfClientPort
	str = str + " --advertise-client-urls http://" + selfIP + ":" + selfClientPort
	if strings.EqualFold(flag, "join") || strings.EqualFold(flag, "existing") {
		str = str + " --initial-cluster "
		peers := getMembers()
		for _, value := range peers {
			str = str + value.Name + "=http://" + value.IP + ":" + value.PeerPort + ","
		}
	} else {
		str = str + " --initial-cluster-token etcd-basic-cluster"
		str = str + " --initial-cluster " + "creator" + selfIP + "=http://" + selfIP + ":" + selfPeerPort + ","
	}

	if strings.EqualFold(flag, "existing") || strings.EqualFold(flag, "join") {
		str = str + " --initial-cluster-state existing"
	} else {
		file.remove(config.BasePath + "/" + "creator" + selfIP + ".etcd")
		str = str + " --initial-cluster-state new"
	}

	//Split configure info (string) to []string
	var args []string = make([]string, 0)
	strs := strings.Split(str, " ")
	for i := 0; i < len(strs); i++ {
		args = append(args, strs[i])
	}
	//Put configure info ([]string) to os.Args
	var len_Args int = 0
	for key, _ := range os.Args {
		os.Args[key] = args[key]
		len_Args = len_Args + 1
	}
	left_args := args[len_Args:]
	os.Args = append(os.Args, left_args...)

	go etcdmain.Main()

	//Send msg to check etcd is running
	timeout := make(chan bool, 1)
	ch := make(chan bool, 1)
	go func() {
		time.Sleep(time.Second * 5)
		timeout <- true
	}()

	go func() {
		for {
			if etcdDone() {
				ch <- true
				break
			}
		}
	}()

	select {
	case <-ch:
		fmt.Println("start ETCD..........")
	case <-timeout:
		logger.Warning("ETCD never response")
		return false
	}

	return true
}

func isClusterAlive() bool {
	url := ""
	/*轮询clusterMembers，直到可以正常访问ETCD为止*/
	for _, value := range clusterMembers {
		url = "members"
		body, err := HttpGet(value.IP, value.ClientPort, url)
		if err != nil {
			//logger.Error(err.Error())
			continue
		}

		d := struct {
			Members []Member
		}{}
		if err = json.Unmarshal(body, &d); err != nil {
			logger.Error(err.Error())
			continue
		}
		if len(d.Members) > 0 {
			return true
		}
	}

	return false
}

func getSlowHeartBeat() int64 {
	/*
	  check if ETCD has my heartbeat record,if exist ,get the rest TTL
	*/
	//peers, _ := getMembers()
	for _, value := range clusterMembers {
		url := "keys/nodeStatusOfSlowheartbeat?quorum=false&recursive=true"
		//fmt.Println("debug: url:", url)
		body, err := HttpGet(value.IP, value.ClientPort, url)
		if err != nil {
			continue
		}

		if strings.Contains(string(body), "key not found") {
			logger.Error(fmt.Errorf("Key not found"))
			return 0
		}

		var m *Event
		err = json.Unmarshal(body, &m)
		if err != nil {
			continue
		}

		if m.Node != nil {
			for _, value := range m.Node.Nodes {
				if selfIP == value.Key[28:] {
					return value.TTL
				}
			}
		}
	}

	return 0
}

func startSlowHeartbeat(identity string) {
	//Broadcast slow heartBeat
	go func() {
		for {
			if broadcastSlowHeartbeatController {
				parameters := "/nodeStatusOfSlowheartbeat/" + identity + selfIP + "?ttl=120&value=ok"
				HttpPut(parameters)
				time.Sleep(time.Second * time.Duration(slowheartbeatCycle))
			}
		}
	}()
}

func getLeaderIP() (string, string, error) {
	url := "members/leader"
	var member Member
	for _, value := range clusterMembers {
		body, err := HttpGet(value.IP, value.ClientPort, url)
		if err != nil {
			//logger.Error(err.Error())
			continue
		}
		err = json.Unmarshal(body, &member)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		break
	}
	if member.ID == "" {
		return "", "", fmt.Errorf("cluster is not accessable!")
	}
	IP := member.PeerURLs[0]
	host := IP[7:len(member.PeerURLs[0])]
	idx := strings.Index(host, ":")
	peerPort := host[idx+1:]
	host = host[:idx]
	return host, peerPort, nil

}
