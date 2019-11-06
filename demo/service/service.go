package service

import (
	"demo/entity"
	"demo/loghelper"
	"log"
)

var curuserinfo = "/src/demo/data/curuser.txt"
var errLog *log.Logger
type User entity.User

func  init()  {
	errLog = loghelp.Error
}

func GetCurUser() (entity.User, bool) {
	if cu, err := entity.GetCurUser(); err != nil {
		return cu, false
	} else {
		return cu, true
	}
}
//Regist
func UserRegister(username string, password string,  email string, phone string) (bool, error)  {
	user := entity.QueryUser(func (u *entity.User) bool {
		return u.Name == username
	})
	if len(user) == 1 {
		errLog.Println("User Register: username already exists")
		return false, nil
	}
	entity.CreateUser(&entity.User{username,password,email,phone})
	if err := entity.Sync(); err != nil {
		return true, err
	}
	return true, nil
}

//Login
func UserLogin(username, password string) bool  {
	user := entity.QueryUser(func (u *entity.User) bool {
		if u.Name == username && u.Password == password {
			return true
		}
		return false
	})
	if len(user) == 0 {
		errLog.Println("Login: User not exist")
		return false
	}
	entity.SetCurUser(&user[0])
	if err := entity.Sync(); err != nil {
		errLog.Println("Login: error happened when set curuser")
		return false
	}
	return true
}

//Logout
func UserLogout() bool {
	if err := entity.Logout(); err != nil {
		return false
	} else {
		return true
	}
}