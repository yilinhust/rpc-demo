package main

import (
	"log"
	"net/http"
	"controller"
)

func main(){
	startHttpServer()
}

func startHttpServer() {
	mux := http.NewServeMux()

	//从net/http包中调用了一个http.HandleFucn函数来注册一个 handler
	//仿佛能看到Spring Boot的RequestMapping注解的go语言实现版
	//跟php里面的控制层（controller）函数类似
	//如果你以前是Python程序员，那么你一定听说过tornado，这个代码和他是不是很像，Go就是拥有类似Python这样动态语言的特性，写Web应用很方便。
	mux.HandleFunc("/", controller.Index)
	mux.HandleFunc("/login", controller.Login)
	mux.HandleFunc("/update/nickname", controller.UpdateNickname)
	mux.HandleFunc("/upload", controller.Upload)
	mux.HandleFunc("/image", controller.GetImage)

	log.Println("Start HTTP server......")
	//不需要一个内嵌的应用服务器,而且这个Web服务内部有支持高并发的特性
	err := http.ListenAndServe(":80", mux)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Start http server error: %s\n", err)
	}
}


