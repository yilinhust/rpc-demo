package model

import (
	"time"
)

type User struct {
	Id    int64
	Username     string
	Password  string
	Nickname  string
	Profile 	string
	CreateTime   time.Time
	LoginTime  time.Time       // 登录时间
}















