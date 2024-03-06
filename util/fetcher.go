package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/oliveagle/jsonpath"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

var httpClient = &http.Client{Timeout: time.Second * 30}

// Get 通过传给的url获取url对应的价格
func Get(url string) ([]byte, error) {
	response, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf(string(bs))
	}

	return bs, nil
}

func FetchPrice(url string, pricePath string) (uint64, error) {
	resp, err := Get(url)
	if err != nil {
		return 0, err
	}
	return responseParse(resp, pricePath)
}

func ParsePrice(p float64) (uint64, error) {
	return uint64(p * math.Pow10(8)), nil //把价格转换成整数再返回
}

func responseParse(response []byte, pricePath string) (uint64, error) {
	var msg interface{}
	err := json.Unmarshal(response, &msg)

	if err != nil {
		return 0, err
	}
	pat, _ := jsonpath.Compile(pricePath)
	res, err := pat.Lookup(msg)
	if err != nil {
		println(err)
	}
	switch res.(type) {
	case string:
		num, err := strconv.ParseFloat(res.(string), 64)
		if err != nil {
			return 0, err
		}
		return ParsePrice(num)
	case float64:
		return ParsePrice(res.(float64))
	default:
		return 0, nil
	}
}

func Post(url, body string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return []byte(""), err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return []byte(""), err
	}
	defer resp.Body.Close()

	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		return []byte(""), err

	}
	return responseBody.Bytes(), nil
	//defer response.Body.Close()

	//bs, err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if response.StatusCode >= 400 {
	//	return nil, fmt.Errorf(string(bs))
	//}
	//
	//return bs, nil
}
