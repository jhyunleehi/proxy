package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	URL = "http://192.168.57.31:8080/"
)

func getweb() {
	log.Printf("getweb call")
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}

	//필요시 헤더 추가 가능
	req.Header.Add("User-Agent", "Crawler")
	//req.Header.Add("Connection", "keep-Alive")
    //req.Header.Add("Keep-Alive", "timeout=5, max=1000")

	// Client객체에서 Request 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 결과 출력
    log.Printf("%v",resp.Header)
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로
	fmt.Println(str)
}

func main() {
	for i := 0; i < 1000; i++ {
		getweb()
	}
}
