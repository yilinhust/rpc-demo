package test

import (
	"testing"
	"dao"
	"model"
	"encrypt"
	"crypto"
	"time"
	"fmt"
	"math/rand"
)


func TestInsertUser(t *testing.T) {
	for i:=1; i<10000000;i++{
		fmt.Print(":", i)

		un := getRandomString()
		fmt.Print("     username:", un)

		algorithm, _ := encrypt.NewHMACAlgorithm(crypto.SHA256, encrypt.HmacKey)
		pwd, _ := algorithm.Encrypt("123456")
		fmt.Println("        pwd:", pwd)

		dao.InsertUser(&model.User{Username:un, Password:pwd, Nickname:"", CreateTime:time.Now()})
	}
}

func  getRandomString() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	count := rand.Intn(12)+6
	for i := 0; i < count; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}





