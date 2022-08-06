package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var serverCount = 0

// These constant is used to define all backend servers
const (
	SERVER1 = "https://www.google.com"
	SERVER2 = "https://www.facebook.com"
	SERVER3 = "https://www.yahoo.com"
	PORT    = "1338"
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
	var servers = []string{SERVER1, SERVER2, SERVER3}
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
	// Note that ServeHttp is non blocking & uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func main() {
	// start server
	http.HandleFunc("/", loadBalacer)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
