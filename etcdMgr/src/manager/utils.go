package manager

import (
	"bufio"
	"config"
	"net/http"
	"os"
	"sync"
)

type myfile struct {
	file  *os.File
	mutex sync.RWMutex // 文件读写锁
}

var n sync.Mutex

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func etcdDone() bool {
	url := "http://" + selfIP + ":" + config.GetEtcdClientPort() + "/v2/members"
	_, err := http.Get(url)
	if err != nil {
		return false
	}
	return true
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func (mf *myfile) writeAndRemove(filename string, Info string) {
	if checkFileIsExist(filename) {
		os.RemoveAll(filename) //删除文件

		f, err := os.Create(filename) //创建文件
		check(err)
		mf := &myfile{
			file: f,
		}

		mf.writeToFile(Info)
		f.Close()

	} else {
		f, err := os.Create(filename)
		check(err)
		mf := &myfile{
			file: f,
		}

		mf.writeToFile(Info)
		f.Close()

	}

}

func (mf *myfile) remove(filename string) {
	if checkFileIsExist(filename) {
		os.RemoveAll(filename) //删除文件
	}
}

func (mf *myfile) writeToFile(Info string) {
	mf.mutex.Lock()
	defer mf.mutex.Unlock()
	w := bufio.NewWriter(mf.file)
	_, err := w.WriteString(Info)
	check(err)
	w.Flush()

}
