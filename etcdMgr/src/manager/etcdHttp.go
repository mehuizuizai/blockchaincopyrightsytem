package manager

import (
	"config"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//make kv
func MakeKV(key string, value string) error {
	url := key + "?value=" + value + "&prevExist=false"
	//fmt.Println("debug:make kv url:" + url)
	return HttpPut(url)
}

//set kv
func SetKV(key string, value string) error {
	url := key + "?value=" + value
	return HttpPut(url)
}

//make kv with ttl
func TTLKV(key string, value string) error {
	url := key + "?ttl=&value=" + value + "&prevExist=false"
	//fmt.Println("debug:make kv url:" + url)
	return HttpPut(url)
}

//
func ElectionTTLKV(key string, value string, ttl int) error {
	url := key + "?ttl=" + strconv.Itoa(ttl) + "&value=" + value + "&prevExist=false"
	//fmt.Println("debug:make kv url:" + url)
	return HttpPut(url)
}

//remove kv
func RemoveKV(ip string, port string, key string) error {
	url := key
	return HttpDelete(ip, port, url)
}

//update kv
func UpdateKV(key string, value string) error {
	url := key + "?value=" + value + "&prevExist=true"
	return HttpPut(url)
}

func HttpGet(ip string, port string, key string) ([]byte, error) {
	url := "http://" + ip + ":" + port + "/v2/" + key
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(20 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*5)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			//ResponseHeaderTimeout: time.Second * 5,
			MaxIdleConnsPerHost: 10,
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		//fmt.Println(err.Error())
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("key not found")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println(err.Error())
		return nil, err
	}

	return body, nil
}

func HttpPut(parameters string) error {
	client := &http.Client{}

	url := "http://127.0.0.1:" + config.GetEtcdClientPort() + "/v2/keys" + parameters
	req, err := http.NewRequest("PUT", url, strings.NewReader("a"))
	check(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	response := string(body)
	if strings.Contains(response, "Key already exists") {
		return fmt.Errorf("ERROR:Key already exists:", response)
	} else if strings.Contains(response, "errorCode") {
		return fmt.Errorf("ERROR:response error:", response)
	}

	return nil
}

func HttpDelete(ip string, port string, parameters string) error {
	client := &http.Client{}

	url := "http://" + ip + ":" + port + "/v2/" + parameters
	req, err := http.NewRequest("DELETE", url, strings.NewReader("a"))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error: DELETE do: ", err.Error())
		return err
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error: DELETE ioutil: ", err.Error())
		fmt.Println("debug: DELETE resp: ", string(body))

	}

	return nil
}
