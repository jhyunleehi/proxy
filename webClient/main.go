package main
 
import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
)
const (
    URL = "http://192.168.57.31:8081/"
)

func getweb(){
    log.Printf("getweb call")
    req, err := http.NewRequest("GET", URL, nil)
    if err != nil {
        panic(err)
    }
 
    //필요시 헤더 추가 가능
    req.Header.Add("User-Agent", "Crawler")
 
    // Client객체에서 Request 실행
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
 
    // 결과 출력
    bytes, _ := ioutil.ReadAll(resp.Body)
    str := string(bytes) //바이트를 문자열로
    fmt.Println(str)
}


func main() {
    for i:=0; i<1000; i++{
        getweb()
    }    
}