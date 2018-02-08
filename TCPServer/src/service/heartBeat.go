package service

import "fmt"

type HeartBeat int

type HeartBeatMessage struct {
	Name string
	Ip string
}

type HeartReplyMessage struct {
	Ack bool
	Ip string
}

func (r *HeartBeat) HeartBeat(msg HeartBeatMessage, reply *HeartReplyMessage) error {
	fmt.Println("Receive HeartBeat: ", msg)
	reply.Ack = true
	return nil
}




