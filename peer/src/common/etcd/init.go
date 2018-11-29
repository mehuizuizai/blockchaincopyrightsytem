package etcd

import (
	"config"
	"logging"
	"net"
	"os"
	"strings"
	"time"
)

var logger = logging.MustGetLogger()

func Initialize() bool {
	var flag bool
	flag = false
	for i := 0; i < 3; i++ {
		ok := communicationWithEtcdMgr()
		if !ok {
			logger.Warning("Get members list from etcdMgr failed!")
			time.Sleep(time.Second * 3)
			continue
		}
		flag = true
		break
	}
	return flag
}

func communicationWithEtcdMgr() bool {
	var socketFilePath = config.GetEtcdMgrPath()
	unixAddr, er := net.ResolveUnixAddr("unix", socketFilePath)
	if er != nil {
		logger.Warning(os.Stderr, "Fatal error: %s", er.Error())
		return false
	}
	unixConn, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		logger.Warning(os.Stderr, "Fatal error: %s", err.Error())
		return false
	}

	ok := sender(unixConn)
	if !ok {
		return false
	}

	return true
}

func sender(conn *net.UnixConn) bool {
	words := "GET"
	n, err := conn.Write([]byte(words))
	//logger.Info("The size of sent message is: ", n)
	if err != nil {
		logger.Error(err.Error())
		return false
	}
	buf := make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	defer conn.Close()

	changeStringToStructs(string(buf[:n]))
	return true

}

func changeStringToStructs(str string) {
	hosts := strings.Split(str, "-")
	var members []MemberInfo
	for i := 0; i < len(hosts)-1; i++ {
		host := strings.Split(hosts[i], ":")
		var member MemberInfo
		member.IP = host[0]
		member.ClientPort = host[1]
		members = append(members, member)
	}

	ClusterMembers = members

}
