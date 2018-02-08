package controller

import (
	"net/http"
	"fmt"
	"io"
	"log"
	"time"
	"crypto/md5"
	"strconv"
	"html/template"
	"session"
	"client"
	"io/ioutil"
)


var SessionManager *session.Manager   //全局的session管理器

var indexTmpl *template.Template


func init() {
	var err error

	indexTmpl, err = template.ParseFiles("static/html/index.html")
	if err != nil {
		log.Println("init error: \n", err)
		return
	}

	SessionManager, _ = session.NewManager("memory","gosessionid",3600)
	go SessionManager.GC()
}



//http://127.0.0.1:80
func Index(writer http.ResponseWriter, request *http.Request) {
	session := SessionManager.SessionStart(writer, request)
	//log.Println("   Index sessionId:", session.SessionID())
	userId := session.Get("userId")
	if userId == nil {
		log.Println("Redirect........")
		http.Redirect(writer, request, "/login", 302)
		return
	}

	var reply UserResponse
	err := client.Call("UserService.QueryUserById", userId, &reply)
	if err != nil {
		log.Println("Index error: \n", err)
		io.WriteString(writer,  "System error......")
		return
	}
	log.Println("QueryUserById.reply:Username", reply.User.Username)

	indexTmpl.Execute(writer, reply.User)
}

//http://127.0.0.1:80/login
func Login(response http.ResponseWriter, request *http.Request) {
	session := SessionManager.SessionStart(response, request)
	//log.Println("   Login sessionId:", session.SessionID())
	userId := session.Get("userId")
	if userId != nil {
		log.Println("   Login Redirect:", session.SessionID())
		http.Redirect(response, request, "/", 302)
		return
	}

	if request.Method == "GET" {
		tmpl, err := template.ParseFiles("static/html/login.html")
		if err != nil {
			log.Println("Login error: \n", err)
			io.WriteString(response,  "System error......")
			return
		}

		token := getToken()
		session := SessionManager.SessionStart(response, request)
		session.Set("token", token)

		tmpl.Execute(response, token)
	} else {
		request.ParseForm()  //解析参数，默认是不会解析的
		token := request.Form.Get("token")
		if token != "" {
			//验证token的合法性
		} else {
			//不存在token报错
		}

		username := template.HTMLEscapeString(request.Form.Get("username"))
		password := template.HTMLEscapeString(request.Form.Get("password"))

		var reply LoginResponse
		err := client.Call("UserService.Login", LoginParams{username, password}, &reply)
		if err != nil {
			//log.Println("Login error: \n", err)
			io.WriteString(response,  "System error......")
			return
		}
		//log.Println("Login reply:", reply)

		if reply.Code == 0 {
			session := SessionManager.SessionStart(response, request)
			session.Set("userId", reply.UserId)

			//log.Println("   Login success:", session.SessionID())
			http.Redirect(response, request, "/", 302)
		}else{
			io.WriteString(response, "<script>alert('"+reply.Msg+"');window.location.href='/login'</script>")
		}
	}
}


//http://127.0.0.1:80/update/nickname
func UpdateNickname(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("method:", request.Method) //获取请求的方法

	session := SessionManager.SessionStart(writer, request)
	fmt.Println("    sessionId:", session.SessionID())
	userId := session.Get("userId")
	if userId == nil {
		http.Redirect(writer, request, "/login", 302)
		return
	}

	if request.Method == "GET" {
		tmpl, err := template.ParseFiles("static/html/update.html")
		if err != nil {
			log.Println("UpdateNickname error: \n", err)
			io.WriteString(writer,  "System error......")
			return
		}
		tmpl.Execute(writer, nil)
	} else {
		request.ParseForm()  //解析参数，默认是不会解析的
		nickname := request.Form.Get("nickname")
		fmt.Println("nickname:", nickname)
		if nickname == "" {
			io.WriteString(writer, "<script>alert('The nickname is required');</script>")
			return
		}

		user := UserParams{Id:userId.(int64), Nickname:nickname}
		var reply Response
		err := client.Call("UserService.UpdateUserNickname", user, &reply)
		if err != nil {
			log.Println("UpdateNickname error: \n", err)
			io.WriteString(writer,  "System error......")
			return
		}
		//fmt.Println("reply:", reply)

		http.Redirect(writer, request, "/", 302)
	}
}


//http://127.0.0.1:80/upload
func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	session := SessionManager.SessionStart(w, r)
	fmt.Println("    sessionId:", session.SessionID())
	userId := session.Get("userId")
	if userId == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("static/html/upload.html")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(2*1024*1024)//把上传的文件存储在内存和临时文件中
		file, handler, err := r.FormFile("uploadfile")//获取文件句柄
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		fmt.Println("    handler.Header=", handler.Header)

		p, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("Upload error: \n", err)
			io.WriteString(w,  "System error......")
			return
		}
		fmt.Println("    fileLength:", len(p))

		var reply Response
		err = client.Call("UserService.UploadUserProfile", p, &reply)
		if err != nil {
			log.Println("Upload error: \n", err)
			io.WriteString(w,  "System error......")
			return
		}
		fmt.Println("reply:", reply)

		http.Redirect(w, r, "/", 302)
	}
}


//http://127.0.0.1:80/image
func GetImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetImage.method:", r.Method) //获取请求的方法
	session := SessionManager.SessionStart(w, r)
	fmt.Println("    sessionId:", session.SessionID())
	userId := session.Get("userId")
	if userId == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	var reply ImageResponse
	err := client.Call("UserService.GetUserProfile", userId.(int64), &reply)
	if err != nil {
		log.Println("GetImage error: \n", err)
		io.WriteString(w,  "System error......")
		return
	}
	fmt.Println("reply:", reply)
	w.Write(reply.Data)
}

func getToken() string {
	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}




type UserParams struct {
	Id    int64
	Nickname  string
}


type Response struct {
	Code int
	Msg string
}

type LoginResponse struct {
	Code int
	Msg string
	UserId int64
}

type UserResponse struct {
	Code int
	Msg string
	User User
}

type ImageResponse struct {
	Code int
	Msg string
	Data []byte
}

type LoginParams struct {
	Username, Password string
}


type User struct {
	Id    int64
	Username     string
	Password  string
	Nickname  string
	Profile 	string
	CreateTime   time.Time
	LoginTime  time.Time       // 登录时间
}














