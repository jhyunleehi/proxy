package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	URL = "https://localhost:10443"
)

func getweb() {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transCfg}

	//response, err := client.Get("https://localhost:10443")

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// 결과 출력
	log.Printf("%v", response.Header)
	htmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(htmlData))
}


func main() {
	for i := 0; i < 1000; i++ {
		getweb()
	}
}

