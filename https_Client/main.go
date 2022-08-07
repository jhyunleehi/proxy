package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	URL = "https://localhost:10443"
)

var Client *http.Client

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	Client = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(3 * time.Second),
	}
}

func getweb() {
	// transCfg := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// client := &http.Client{Transport: transCfg}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}

	response, err := Client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	defer response.Body.Close()

	// 결과 출력
	log.Printf("%v", response.Header)
	htmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(htmlData))
}

func main() {
	for i := 0; i < 1000; i++ {
		getweb()
	}
}
