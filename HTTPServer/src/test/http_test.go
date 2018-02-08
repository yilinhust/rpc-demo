package test

import (
	"testing"
	"time"
	"math/rand"
	"net/http"
	"log"
	"strings"
	"io/ioutil"
	"runtime"
	"net/http/cookiejar"
	"fmt"
)


func Benchmark_Login(b *testing.B) {
	runtime.GOMAXPROCS( runtime.NumCPU()*200)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			loginRequest()
		}
	})
}

func loginRequest(){
	un := getRandomString()
	//un := "test"
	//fmt.Print("     username:", un)

	req, err := http.NewRequest("POST", "http://127.0.0.1/login",
							strings.NewReader("username="+un +"&password=123456"))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar:jar}
	//client.CheckRedirect =  func(req *http.Request, via []*http.Request) error {
	//		return http.ErrUseLastResponse
	//		}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("..........httpclient.POST error:  ", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if data == nil{
		fmt.Print("     data:", string(data[:]))
	}
}



var bytes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func  getRandomString() string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	result := []byte{}
	for i := 0; i < 3; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}



