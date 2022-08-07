package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var counter int
var mutex = &sync.Mutex{}

func echoString(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s][%s][%s][%v]", r.Method, r.Proto, r.URL.Path, r.Header)
	fmt.Fprintf(w, "hello  : ")
	mutex.Lock()
	counter++
	time.Sleep(time.Millisecond * 100)
	mutex.Unlock()
	fmt.Fprintf(w, strconv.Itoa(counter))
	//time.Sleep(time.Second * 1)
	log.Printf("[%s]", strconv.Itoa(counter))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", echoString)
	mux.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	server := &http.Server{Addr: ":8081", Handler: mux}
	//server.SetKeepAlivesEnabled(false)
	server.ListenAndServe()
	//http.ListenAndServe(":8081", nil)
}
