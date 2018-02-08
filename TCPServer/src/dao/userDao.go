package dao

import (
	"model"
	"log"
)


func QueryUserByUsername(username string) *model.User{
	if username=="" {
		return nil
	}
	//log.Println("DB...", db)
	stmt, err := db.Prepare("SELECT id,username,password,nickname,IFNULL(profile, ''),CreateTime FROM t_user WHERE username=?")
	if err != nil {
		log.Println("QueryUserByUsername Prepare error...", err)
		return nil
	}
	defer stmt.Close()

	rows, err2 := stmt.Query(username)
	if err2 != nil {
		log.Println("Prepare error...", err2)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		user := model.User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Password, &(user.Nickname), &(user.Profile), &(user.CreateTime))
		checkErr(err)
		return &user
	}
	return nil
}

func QueryUserById(id int64) *model.User{
	if id <= 0 {
		return nil
	}
	stmt, err := db.Prepare("SELECT id,username,password,nickname,IFNULL(profile, ''),CreateTime FROM t_user WHERE id=?")
	if err != nil {
		log.Println("Prepare error...", err)
		return nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		log.Println("Query error...", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		user := &model.User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.Profile, &user.CreateTime)
		checkErr(err)
		return user
	}
	return nil
}


func InsertUser(user *model.User) int64{
	stmt, err := db.Prepare("INSERT t_user SET username=?,password=?,nickname=?,CreateTime=?")
	defer stmt.Close()
	checkErr(err)

	res, err := stmt.Exec(user.Username, user.Password, user.Nickname, user.CreateTime)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}


func UpdateUserNickname(userId int64, nickname string) int64{
	stmt, err := db.Prepare("UPDATE t_user SET nickname=? WHERE ID=?")
	defer stmt.Close()
	checkErr(err)

	res, err := stmt.Exec(nickname, userId)
	checkErr(err)

	count, err := res.RowsAffected()
	checkErr(err)

	return count
}




func checkErr(err error) {
	if err != nil {
		log.Println("init database successfully...", db)
		panic(err)
	}
}



