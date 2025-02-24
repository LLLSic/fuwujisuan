package entity

import (
	"os"
	"io"
	"bufio"
	"path/filepath"
	"errors"
	"log"
	"encoding/json"
	"demo/loghelper"
	"demo/deepcopy"
)

type UserFilter func (*User) bool

var userinfoPath = "/src/demo/data/userinfo"
var curUserPath = "/src/demo/data/curuser.txt"

var curUserName *string
var dirty bool
var uData []User
var errLog *log.Logger

func init()  {
	errLog = loghelp.Error
	dirty = false
	userinfoPath  = filepath.Join(loghelp.GoPath, userinfoPath)
	curUserPath = filepath.Join(loghelp.GoPath, curUserPath)
	if err := readFromFile(); err != nil {
		errLog.Println("readFromFile fail: ",err)
	}
}

func Sync() error {
	if err := writeToFile(); err != nil {
		errLog.Println("writeToFile fail: ",err)
		return err
	}
	return nil
}
func CreateUser(v *User)  {
	uData = append(uData, deepcopy.Copy(*v).(User))
	dirty = true
}
func QueryUser(filter UserFilter)  []User{
	var user []User
	for _, v := range uData {
		if filter(&v) {
			user = append(user,v)
		}
	}
	return user
}

func SetCurUser(u *User)  {
	if u == nil {
		curUserName = nil
		return
	}
	if curUserName == nil {
		p := u.Name
		curUserName = &p
	} else {
		*curUserName = u.Name
	}
}

//return current user
func GetCurUser() (User, error) {
	if curUserName == nil {
		return User{}, errors.New("Current user doesn't exist")
	}
	for _, v := range uData {
		if v.Name == *curUserName {
			return v, nil
		}
	}
	return User{}, errors.New("Current user doesn't exist")
}

func Logout() error {
	curUserName = nil
	return Sync()
}
func readFromFile() error {
	var e []error
	str, err1 := readString(curUserPath)
	if err1 != nil {
		e = append(e, err1)
	}
	curUserName = str
	if err := readUser(); err != nil {
		e = append(e, err)
	}
	//if err := readMet(); err != nil {
	//	e = append(e, err)
	//}
	if len(e) == 0 {
		return nil
	}
	er := e[0]
	for i := 1; i < len(e); i++ {
		er = errors.New(er.Error() + e[i].Error())
	}
	return er
}

// writeToFile : write file content from memory
// @return if fail, error will be returned
func writeToFile() error {
	var e []error
	if err := writeString(curUserPath, curUserName); err != nil {
		e = append(e, err)
	}
	if dirty {
		if err := writeJSON(userinfoPath, uData); err != nil {
			e = append(e, err)
		}
		//if err := writeJSON(metinfoPath, mData); err != nil {
		//	e = append(e, err)
		//}
	}
	if len(e) == 0 {
		return nil
	}
	er := e[0]
	for i := 1; i < len(e); i++ {
		er = errors.New(er.Error() + e[i].Error())
	}
	return er
}

func readUser() error {
	file, err := os.Open(userinfoPath);
	if err != nil {
		errLog.Println("Open File Fail:", userinfoPath, err)
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	switch err := dec.Decode(&uData); err {
	case nil, io.EOF:
		return nil
	default:
		errLog.Println("Decode User Fail:", err)
		return err
	}
}

/*
func readMet() error {
	file, err := os.Open(metinfoPath);
	if err != nil {
		errLog.Println("Open File Fail:", metinfoPath, err)
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	switch err := dec.Decode(&mData); err {
	case nil, io.EOF:
		return nil
	default:
		errLog.Println("Decode Met Fail:", err)
		return err
	}
}
*/

func writeJSON(fpath string, data interface{}) error {
	file, err := os.Create(fpath);
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)

	if err := enc.Encode(&data); err != nil {
		errLog.Println("writeJSON:", err)
		return err
	}
	return nil
}

func writeString(path string, data *string) error {
	file, err := os.Create(path)
	if err != nil {
		loghelp.Error.Println("Create file error:", path)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if data != nil {
		if _, err := writer.WriteString(*data); err != nil {
			loghelp.Error.Println("Write file fail:", path)
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		loghelp.Error.Println("Flush file fail:", path)
		return err
	}
	return nil
}

func readString(path string) (*string, error) {
	file, err := os.Open(path)
	if err != nil {
		loghelp.Error.Println("Open file error:", path)
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	str, err := reader.ReadString('\n');
	if err != nil && err != io.EOF {
		loghelp.Error.Println("Read file fail:", path)
		return nil, err
	}
	return &str, nil
}
