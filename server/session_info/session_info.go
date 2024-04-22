package session_info

import (
	"github.com/gin-gonic/gin"
)


type SessionInfo interface{
	SetAuthToken(string)
	GetAuthToken() (string, error)
	SetUserId(id int64)
	GetUserId() (int64, error)
	SetRoleId(id int64)
	GetRoleId() (int64, error)
	SetUserName(name string)
	GetUserName() (string, error)
	SetSelectedMenu(string)
	GetSelectedMenu() (string, error)
}

type sessionInfo struct {
	ctx					*gin.Context
	authToken		string
	userId			int64
}

func (s sessionInfo) SetAuthToken(t string) {
	s.ctx.Set("AuthToken", t)
}

func (s sessionInfo) GetAuthToken() (string, error) {
	return s.ctx.GetString("AuthToken"), nil
}

func (s sessionInfo) SetUserId(id int64) {
	s.ctx.Set("User.Id", id)
}

func (s sessionInfo) GetUserId() (int64, error) {
	return s.ctx.GetInt64("User.Id"), nil
}

func (s sessionInfo) SetRoleId(id int64) {
	s.ctx.Set("Role.Id", id)
}

func (s sessionInfo) GetRoleId() (int64, error) {
	return s.ctx.GetInt64("Role.Id"), nil
}

func (s sessionInfo) SetUserName(name string) {
	s.ctx.Set("User.Name", name)
}

func (s sessionInfo) GetUserName() (string, error) {
	return s.ctx.GetString("User.Name"), nil
}

func (s sessionInfo) SetSelectedMenu(menu string) {
	s.ctx.Set("SelectedMenu", menu)
}

func (s sessionInfo) GetSelectedMenu() (string, error) {
	return s.ctx.GetString("SelectedMenu"), nil
}


func NewSessionInfo(c *gin.Context) SessionInfo {
	return sessionInfo{
		ctx: 				c, 
		authToken: 	"", 
		userId: 		0, 
	}
}
