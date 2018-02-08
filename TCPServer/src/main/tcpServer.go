package main

import (
	"service"
	"log"
	"net/rpc"
	"net"
)



func main(){
	startTcpServer()
}


func startTcpServer() {
	//注册rpc服务
	heartBeat := new(service.HeartBeat)
	rpc.Register(heartBeat)
	rect := new(service.UserService)
	rpc.Register(rect)

	//获取tcpaddr
	tcpaddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:1234")
	checkErr(err)

	//监听端口
	tcplisten, err2 := net.ListenTCP("tcp", tcpaddr)
	checkErr(err2)
	log.Println("Start Tcp server......")

	//死循环处理连接请求
	for {
		conn, err3 := tcplisten.Accept()
		if err3 != nil {
			log.Println("Accept error ......")
			continue
		}
		log.Println("Received new connection from ", conn.RemoteAddr())
		//使用goroutine单独处理rpc连接请求
		go rpc.ServeConn(conn)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Open Tcp server error: %s\n", err)
	}
}

