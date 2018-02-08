package test

import (
	"testing"
	"fmt"
	"service"
	"model"
	"time"
	"math/rand"
	"runtime"
	"log"
)


func Benchmark_Login(b *testing.B) {
	runtime.GOMAXPROCS( runtime.NumCPU()*100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			login4Service()
		}
	})
}

func login4Service() {
	un := GetRandomString()
	//fmt.Println("     username:", un)
	//un := "test"

	us := &service.UserService{}
	user := service.LoginParams{Username:un, Password:"123456"}
	var reply service.LoginResponse

	err := us.Login(user, &reply)
	if err != nil {
		log.Println("Login error: ", err)
	}
}

func insertUser4Service() {
	un := GetRandomString()
	fmt.Print("     username:", un)
	us := &service.UserService{}
	us.InsertUser(&model.User{Username:un, Password:"123456", Nickname:"", CreateTime:time.Now()}, nil)
}



var bytes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func  GetRandomString() string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	result := []byte{}
	for i := 0; i < 4; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}











