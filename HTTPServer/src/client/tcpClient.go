package client

import (
	"net/rpc"
	"log"
	"time"
	"strconv"
	"fmt"
)

var client *rpc.Client

func init(){
	connect()
	go Handler()
}

func connect(){
	clientTmp, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Println("Connect tcp server error: \n", err)
		return
	}
	client = clientTmp
}

func Handler() {
	go heartBeat()
}

func heartBeat() {
	i := 0
	for {
		time.Sleep(5 * time.Second)
		var reply HeartReplyMessage
		fmt.Println("send HeartBeat "+strconv.Itoa(i))
		i++

		if client == nil {
			connect()
			continue
		}
		m := HeartBeatMessage{"http", "127.0.0.1"}
		err := client.Call("HeartBeat.HeartBeat", m, &reply)

		if err != nil {
			log.Println("heartBeat error...", err)
			connect()
			continue
		}
		if !reply.Ack{
			log.Println("heartBeat unsuccessfully...")
			connect()
		}
	}
}


func GetClient() *rpc.Client{
	if client == nil {
		log.Println("GetClient: get connection ......")
		connect()
	}
	return client
}

func Call(serviceMethod string, args interface{}, reply interface{}) error {
	err := GetClient().Call(serviceMethod, args, reply)
	if err == rpc.ErrShutdown {
		log.Println("connection is shut down ......")
		connect()
		return GetClient().Call(serviceMethod, args, reply)
	}
	return err
}



type HeartBeatMessage struct {
	Name string
	Ip string
}

type HeartReplyMessage struct {
	Ack bool
	Ip string
}

















