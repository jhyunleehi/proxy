# Web Server Test

## Simple Web Server - We Client

### 1. WebServer-Client

##### WebServer

```go
package main
import (fmt"	"log"	"net/http"	"strconv"	"sync"	"time")
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
```

##### WebClient

```go
package main

import (	"fmt"	"io/ioutil"	"log"	"net/http")

const (
	URL = "http://192.168.57.31:8081/"
)

func getweb() {
	log.Printf("getweb call")
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("User-Agent", "Crawler")
//	req.Header.Add("Connection", "keep-Alive")
//  req.Header.Add("Keep-Alive", "timeout=5, max=1000")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

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
```



#### 1. Clinet  keep-Alive 테스트

* ​	client := &http.Client{} 이 코드에서  http 객체를 새롭게 만들기 때문에 keep alive 유지 가 안될 것 같은데... 실제 테스트 해보면 이것과는 무관하게 keep-alive 유지된다.
* Client network 상태

```
$ netstat -na | grep  8081
tcp        0      0 192.168.57.6:50706      192.168.57.31:8081      ESTABLISHED
```

* server network 

```
$ netstat -na |  grep  8081| sort
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50706      ESTABLISHED
tcp6       0      0 :::8081                 :::*                    LISTEN
```



#### 2. Client Header에서 keep-alive 속성 빼고 테스트 

* Client Header에서 두개 속성을 빼도 동일하게 keep-alive 유지된다. 

```
//	req.Header.Add("Connection", "keep-Alive")
//   req.Header.Add("Keep-Alive", "timeout=5, max=1000")
```



#### 3. WebServer에서 server.SetKeepAlivesEnabled(false) 설정

* `	server.SetKeepAlivesEnabled(false)` 설정
* server에서 netwrok 세션이 재활용 안된다. 

```sh
$ netstat -na |  grep  8081| sort
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50712      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50714      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50716      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50718      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50720      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50722      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50724      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50726      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50728      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50730      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50732      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50734      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50736      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50738      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50740      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50742      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50744      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50746      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50748      TIME_WAIT  
tcp6       0      0 192.168.57.31:8081      192.168.57.6:50750      TIME_WAIT  
tcp6       0      0 :::8081                 :::*                    LISTEN  
```



```
C:\Gocode\src>curl -v http://localhost:8081/
*   Trying 127.0.0.1:8081...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> GET / HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/7.83.1
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Sat, 06 Aug 2022 03:35:20 GMT
< Content-Length: 5
< Content-Type: text/plain; charset=utf-8
<
hello* Connection #0 to host localhost left intact

C:\Gocode\src>
```

## 

## Client->ReverseProxy->WebServer

[WebClient] -->[ReverseProxy] -->[WebServer]

* 각 구간에서 Keep-Alive 동작 상태를 확인하려고 함. 
* 잘 동작하네.. !!

##### web Client

```go
package main

import ( "fmt"	"io/ioutil"	"log"	"net/http")

const (
	URL = "http://192.168.57.31:8080/"  //proxy url 
)

func getweb() {
	log.Printf("getweb call")
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("User-Agent", "Crawler")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

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
```



##### reserse Procy

```go
package main
 
import ( "errors" "fmt"     "log"    "net/http"    "net/http/httputil"    "net/url")
 
// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
    url, err := url.Parse(targetHost)
    if err != nil {
        return nil, err
    }
 
    proxy := httputil.NewSingleHostReverseProxy(url)
 
    originalDirector := proxy.Director
    proxy.Director = func(req *http.Request) {
        originalDirector(req)
        modifyRequest(req)
    }
 
    //proxy.ModifyResponse = modifyResponse()   //<<-- return 전달위해서 
    //proxy.ErrorHandler = errorHandler()
    return proxy, nil
}
 
func modifyRequest(req *http.Request) {
    req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}
 
func errorHandler() func(http.ResponseWriter, *http.Request, error) {
    return func(w http.ResponseWriter, req *http.Request, err error) {
        fmt.Printf("Got error while modifying response: %v \n", err)
        return
    }
}
 
func modifyResponse() func(*http.Response) error {
    return func(resp *http.Response) error {
        return errors.New("response body is invalid")
    }
}
 
// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        proxy.ServeHTTP(w, r)
    }
}
 
func main() {
    // initialize a reverse proxy and pass the actual backend server url here
    proxy, err := NewProxy("http://192.168.57.31:8081")    //<<-- web server url 
    if err != nil {
        panic(err)
    }
 
    // handle all requests to your server using the proxy
    http.HandleFunc("/", ProxyRequestHandler(proxy))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```



##### web server

```go
package main
import (fmt"	"log"	"net/http"	"strconv"	"sync"	"time")
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
```



#### 1. WebClient - Proxy 구간

* keep alive  동작함.

```
$ netstat -na | grep  8080
tcp        0      0 192.168.57.6:35390      192.168.57.31:8080      ESTABLISHED
```



#### 2. proxy-Webserver 구간

* keep-alive 동작하지 않음.

