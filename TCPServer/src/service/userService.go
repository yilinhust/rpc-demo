package service

import (
	"model"
	"dao"
	"encrypt"
	"crypto"
	"log"
	"fmt"
	"io/ioutil"
)

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
	User model.User
}

type ImageResponse struct {
	Code int
	Msg string
	Data []byte
}

type LoginParams struct {
	Username, Password string
}


type UserService struct{}



func (u *UserService) InsertUser(user *model.User, reply *Response) error{
	log.Println("InsertUser:user:", user)
	algorithm, _ := encrypt.NewHMACAlgorithm(crypto.SHA256, encrypt.HmacKey)
	pwd, _ := algorithm.Encrypt(user.Password)
	user.Password = pwd
	log.Println("InsertUser:Password:", user.Password)

	id := dao.InsertUser(user)
	if id <= 0 {
		return fmt.Errorf("InsertUser Error")
	}
	*reply = Response{Code:0, Msg:"Success"}
	log.Println("InsertUser:reply:", reply)
	return nil
}


func (u *UserService) UpdateUserNickname(user model.User, reply *Response) error{
	log.Println("UpdateUserNickname:user=", user)
	count := dao.UpdateUserNickname(user.Id, user.Nickname)
	if count != 1 {
		return fmt.Errorf("UpdateUserNickname Error")
	}
	*reply = Response{Code:0, Msg:"Success"}
	log.Println("UpdateUserNickname:reply:", *reply)
	return nil
}

func (u *UserService) Login(param LoginParams, reply *LoginResponse) error{
	//log.Println("Login.param=", param)
	user := dao.QueryUserByUsername(param.Username)
	if user == nil {
		*reply = LoginResponse{Code:10, Msg:"ERROR: Incorrect username or password"}
		//log.Println("Login:reply=", *reply)
		return nil
	}

	algorithm, _ := encrypt.NewHMACAlgorithm(crypto.SHA256, encrypt.HmacKey)
	err := algorithm.Verify(param.Password, user.Password)
	if err != nil {
		*reply = LoginResponse{Code:10, Msg:"ERROR: Incorrect username or password"}
		//log.Println("Login:reply=", *reply)
		return nil
	}
	*reply = LoginResponse{Code:0, Msg:"Success", UserId:user.Id}
	//log.Println("Login Success:param=", param)
	return nil
}


func (u *UserService) QueryUserById(id int64, reply *UserResponse) error{
	//log.Println("QueryUserById:id=", id)
	user := dao.QueryUserById(id)
	if user == nil {
		*reply = UserResponse{Code:20, Msg:"System error"}
		return nil
	}
	*reply = UserResponse{Code:0, Msg:"Success", User:*user}
	//log.Println("QueryUserById:reply=", *reply)
	return nil
}



func (u *UserService) UploadUserProfile(data []byte, reply *Response) error{
	log.Println("UploadUserProfile:file:", len(data))
	err := ioutil.WriteFile("./upload/output.jpg", data, 0666)
	if err != nil {
		*reply = Response{Code:30, Msg:"UploadUserProfile fail"}
		log.Println("UploadUserProfile:reply=", *reply)
		return nil
	}
	*reply = Response{Code:0, Msg:"Success"}
	log.Println("UploadUserProfile:reply:", *reply)
	return nil
}



func (u *UserService) GetUserProfile(userId int64, reply *ImageResponse) error{
	log.Println("GetUserProfile:userId:", userId)
	data, err := ioutil.ReadFile("./upload/output.jpg")
	if err != nil {
		*reply = ImageResponse{Code:40, Msg:"GetUserProfile fail"}
		log.Println("GetUserProfile:reply=", *reply)
		return nil
	}
	*reply = ImageResponse{Code:0, Msg:"Success", Data:data}
	log.Println("GetUserProfile:reply:", *reply)
	return nil
}














