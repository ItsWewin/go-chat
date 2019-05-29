package process

import (
	"encoding/json"
	"fmt"
	commen "go-chat/commen/message"
	"go-chat/server/model"
	"go-chat/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn
}

func login(userName, passWord string) (user model.User, err error) {
	// 判断用户名和密码
	user, err = model.CurrentUserDao.Login(userName, passWord)
	return
}

// 响应客户端
func (this *UserProcess) responseClient(code int, err error) {
	var responseMessage commen.ResponseMessage
	responseMessage.Code = code
	if err != nil {
		responseMessage.Error = fmt.Sprintf("login error: %v", err)
	}

	responseData, err := json.Marshal(responseMessage)
	if err != nil {
		fmt.Printf("some error when generate response message, error: %v", err)
	}

	dispatcher := utils.Dispatcher{Conn: this.Conn}
	err = dispatcher.WirteData(responseData)
}

func (this *UserProcess) UserLogin(message string) (err error) {
	var info commen.LoginMessage
	var code int
	err = json.Unmarshal([]byte(message), &info)
	if err != nil {
		code = commen.ServerError
	}

	_, err = login(info.UserName, info.Password)
	switch err {
	case nil:
		code = commen.LoginSucceed
	case model.ERROR_USER_NOT_EXISTS:
		code = 404
	case model.ERROR_USER_PWD:
		code = 403
	default:
		code = 500
	}
	this.responseClient(code, err)
	return
}
