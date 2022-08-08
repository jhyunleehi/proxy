package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var serverCount = 0

// These constant is used to define all backend servers
const (
	SERVER = "https://localhost:9443"
)

// Given a request send it to the appropriate url
func loadBalacer(res http.ResponseWriter, req *http.Request) {
	// Get address of one backend server on which we forward request
	url := getProxyURL()
	// Log the request
	logRequestPayload(url)
	// Forward request to original request
	serveReverseProxy(url, res, req)
}

// Balance returns one of the servers using round-robin algorithm
func getProxyURL() string {
	var servers = []string{SERVER}
	server := servers[serverCount]
	serverCount++
	// reset the counter and start from the beginning
	if serverCount >= len(servers) {
		serverCount = 0
	}
	return server
}

// Log the redirect url
func logRequestPayload(proxyURL string) {
	log.Printf("proxy_url: %s\n", proxyURL)
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)
	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &http.Transport{
		MaxIdleConns:          1000,
		IdleConnTimeout:       time.Duration(60) * time.Second,
		ResponseHeaderTimeout: time.Duration(10) * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	proxy.ServeHTTP(res, req)
}

func main() {
	http.Handle("/", http.HandlerFunc(loadBalacer))
	log.Fatal(http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil))
}
