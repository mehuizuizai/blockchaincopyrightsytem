package DbUtil

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"strconv"
)

var txhash = make([]string, 0)

func Store(data interface{}, filename string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}

}
func Load(data interface{}, filename string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		//		panic(err)
		empty := []string{}
		Store(empty, filename)
	}
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}
}
func TestPesitence() {
	for i := 0; i < 2; i++ {
		txhash = append(txhash, strconv.Itoa(i))
	}
	fmt.Println("txhash", txhash)
	Store(txhash, "test1")
	postRead := []string{}
	Load(&postRead, "test1")
	fmt.Println(postRead)
}
