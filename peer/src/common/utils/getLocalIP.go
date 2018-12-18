package utils

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func GetlocalIP() string {
	inDocker := false

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var IPs []string

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil { //ipv4
				IPs = append(IPs, ipnet.IP.To4().String())
			}
		}
	}

	if len(IPs) == 1 {
		return IPs[0]
	} else if len(IPs) == 0 {
		return ""
	} else if len(IPs) == 2 { //commonly in docker
		for _, value := range IPs {
			if !strings.Contains(value, "172.17") {
				return value
			}
		}
		return ""
	} else if inDocker == false { //commomly in real machine
		//===========temp============
		for _, value := range IPs {
			if !strings.Contains(value, "172") {
				return value
			}
		}
		return ""
		//===========temp============
	} else {
		return ""
	}
}
