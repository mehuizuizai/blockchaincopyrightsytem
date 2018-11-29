package etcd

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetKV(key string) ([]byte, error) {
	for _, value := range ClusterMembers {
		body, statusCode, err := httpGet(value.IP, value.ClientPort, key)
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") { //there exist node to access in cluster and have not the key
				return nil, fmt.Errorf("key not found")
			}
			continue
		}

		if statusCode == 404 {
			return nil, fmt.Errorf("key not found")
		}

		return body, nil
	}

	return nil, fmt.Errorf("cluster is unaccessable!")
}

//make kv
func MakeKV(key string, value string) error {
	url := key + "?value=" + value + "&prevExist=false"
	for _, value := range ClusterMembers {
		err := httpPut(value.IP, value.ClientPort, url)
		if err != nil {
			//logger.Error(err.Error())
			if strings.Contains(err.Error(), "Key already exists") || strings.Contains(err.Error(), "errorCode") {
				return err
			}
			continue
		}
		return nil
	}
	return fmt.Errorf("cluster is unaccessable!")
}

//set kv
func SetKV(key string, value string) error {
	url := key + "?value=" + value
	for _, value := range ClusterMembers {
		err := httpPut(value.IP, value.ClientPort, url)
		if err != nil {
			logger.Warning(err.Error())
			continue
		}
		return nil
	}
	return fmt.Errorf("cluster is unaccessable!")
}

func SetTTLKV(parameters string) error {
	for _, value := range ClusterMembers {
		err := httpPut(value.IP, value.ClientPort, parameters)
		if err != nil {
			//logger.Warning(err.Error())
			continue
		}
		return nil
	}
	return fmt.Errorf("cluster is unaccessable!")
}

//make kv with ttl
func TTLKV(key string, value string) error {
	url := key + "?ttl=&value=" + value + "&prevExist=false"
	for _, value := range ClusterMembers {
		err := httpPut(value.IP, value.ClientPort, url)
		if err != nil {
			logger.Warning(err.Error())
			if strings.Contains(err.Error(), "Key already exists") || strings.Contains(err.Error(), "errorCode") {
				return err
			}
			continue
		}
		return nil
	}
	return fmt.Errorf("cluster is unaccessable!")
}

func ElectionTTLKV(key string, value string, ttl int) error {
	url := key + "?ttl=" + strconv.Itoa(ttl) + "&value=" + value + "&prevExist=false"
	for _, value := range ClusterMembers {
		err := httpPut(value.IP, value.ClientPort, url)
		if err != nil {
			//logger.Warning(err.Error())
			if strings.Contains(err.Error(), "Key already exists") || strings.Contains(err.Error(), "errorCode") {
				return err
			}
			continue
		}
		return nil
	}
	return fmt.Errorf("cluster is unaccessable!")
}

//remove kv
func RemoveKV(key string) error {
	url := key
	for _, value := range ClusterMembers {
		err := httpDelete(value.IP, value.ClientPort, url)
		if err != nil {
			logger.Warning(err.Error())
			continue
		}
		return nil
	}
	return fmt.Errorf("cluster is unaccessable!")
}

//update kv
func UpdateKV(key string, value string) error {
	url := key + "?value=" + value
	//+ "&prevExist=true"
	for _, value := range ClusterMembers {
		err := httpPut(value.IP, value.ClientPort, url)
		if err != nil {
			logger.Warning(err.Error())
			continue
		}
		return nil
	}
	return fmt.Errorf("cluster is unaccessable!")
}

func httpGet(ip string, port string, key string) ([]byte, int, error) {
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
			ResponseHeaderTimeout: time.Second * 5,
			MaxIdleConnsPerHost:   10,
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	statusCode := resp.StatusCode
	if statusCode == 404 {
		return nil, statusCode, nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, statusCode, err
	}

	return body, statusCode, nil
}

func httpPut(ip, port, parameters string) error {
	client := &http.Client{}

	url := "http://" + ip + ":" + port + "/v2/keys" + parameters
	req, err := http.NewRequest("PUT", url, strings.NewReader("a"))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := string(body)
	if strings.Contains(response, "Key already exists") {
		return fmt.Errorf("ERROR:Key already exists:", response)
	} else if strings.Contains(response, "errorCode") {
		return fmt.Errorf("ERROR:response error:", response)
	}

	return nil
}

func httpDelete(ip string, port string, parameters string) error {
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
