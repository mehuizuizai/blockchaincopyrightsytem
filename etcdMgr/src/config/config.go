package config

import (
	"fmt"
	"os"
	"path/filepath"

	"strings"

	"github.com/spf13/viper"
)

const cmdRoot = "etcdMgr" //config.yaml
var BasePath string

func Initialize() error {

	// For environment variables.
	viper.SetEnvPrefix(cmdRoot)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	//设置配置文件的搜索路径
	executable, _ := os.Executable()
	BasePath = filepath.Dir(executable) + "/"

	//设置配置文件的搜索路径
	viper.AddConfigPath(BasePath) // Path to look for the config file in
	//设置配置文件的名称
	viper.SetConfigName(cmdRoot) // Name of config file (without extension)
	//设置配置文件类型
	viper.SetConfigType("yaml")
	//读取配置文件
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		fmt.Println(err)
		return fmt.Errorf("Fatal error when reading %s config file: %s\n", cmdRoot, err)

	}

	return nil
}

func GetLocalHostIP() string {
	ip := viper.GetString("network.localhost")
	return ip
}

func GetEtcdCluterMembers() string {
	members := viper.GetString("etcd.clustermembers")
	if members == "" {
		panic(fmt.Errorf("Fatal error:no etcd clustermembers"))
	}
	return members
}

func GetEtcdClusterCapacity() string {
	port := viper.GetString("etcd.capacity")
	if port == "" {
		panic(fmt.Errorf("Fatal error:no ETCD cluster capacity"))
	}
	return port
}

func GetEtcdIdentity() string {
	identity := viper.GetString("etcd.identity")
	if identity == "" {
		panic(fmt.Errorf("Fatal error:no ETCD identity"))
	}
	return identity
}

func GetEtcdClientPort() string {
	port := viper.GetString("etcd.clientport")
	if port == "" {
		panic(fmt.Errorf("Fatal error:no ETCD client port"))
	}
	return port
}

func GetEtcdPeerPort() string {
	port := viper.GetString("etcd.peerport")
	if port == "" {
		panic(fmt.Errorf("Fatal error:no ETCD peer port"))
	}
	return port
}

func GetWebServPort() string {
	port := viper.GetString("webservice.port")
	if port == "" {
		panic(fmt.Errorf("Fatal error:no web service port"))
	}
	return port

}

const DEFAULT_WEBSERVICE_CONNS int = 1000

func GetWebServConns() int {
	conns := viper.GetInt("webservice.conns")
	if conns == 0 {
		return DEFAULT_WEBSERVICE_CONNS
	}
	return conns
}

func GetLogLevel() string {
	level := viper.GetString("log.level")
	if level == "" {
		level = "debug,info,warn,error,critical"
	}
	return level
}
func GetLogType() string {
	type_ := viper.GetString("log.type")
	if type_ == "" {
		type_ = "size"
	}
	return type_
}
func GetLogMaxdays() string {
	maxdays := viper.GetString("log.maxdays")
	if maxdays == "" {
		maxdays = "7" //one week
	}
	return maxdays
}
func GetLogMaxfiles() string {
	maxfiles := viper.GetString("log.maxfiles")
	if maxfiles == "" {
		maxfiles = "5"
	}
	return maxfiles
}
func GetLogMaxsize() string {
	maxsize := viper.GetString("log.maxsize")
	if maxsize == "" {
		maxsize = "1048576" //1M
	}
	return maxsize
}
func GetLogFileName() string {
	filename := viper.GetString("log.filename")
	if filename == "" {
		filename = "log.txt"
	}
	return filename
}

func GetEtcdMgrPath() string {
	path := viper.GetString("etcdmgr.path")
	if path == "" {
		path = "/tmp/unix_socket"
	}
	return path
}

func GetVerifyFlag() bool {
	flag := viper.GetBool("webservice.verifyflag")
	return flag
}
