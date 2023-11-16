package pkg

import (
	"time"

	"github.com/kataras/jwt"
)

type WxAuthInfo struct {
	ID         string    `json:"id"`
	CAK        string    `json:"cak"`
	AppID      string    `json:"appID"`
	OpenID     string    `json:"openID"`
	UnionId    string    `json:"unionID"`
	NickName   string    `json:"nickName"`
	HeadImgURL string    `json:"headImgURL"`
	Sex        string    `json:"sex"`
	Country    string    `json:"country"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	Subscribe  int       `json:"subscribe"`
	CreateTime time.Time `json:"createTime" time_format:"2006-01-02 15:04:05"`
}

type WxAuthInfoClaims struct {
	jwt.Claims
	WxAuthInfo
}

type BaseUserInfo struct {
	ID              string      `json:"id"`
	NickName        string      `json:"nickName"`
	Avatar          string      `json:"avatar"`
	UserName        string      `json:"userName"`
	Sex             string      `json:"sex"`
	Birth           time.Time   `json:"birth"`
	Mobile          string      `json:"mobile"`
	WeiXin          string      `json:"weiXin"`
	QQ              string      `json:"qq"`
	EMail           string      `json:"email"`
	Province        string      `json:"province"`
	City            string      `json:"city"`
	District        string      `json:"district"`
	Street          string      `json:"street"`
	Roles           interface{} `json:"roles"`
	Description     string      `json:"description"`
	Address         []string    `json:"address"`
	UserType        int         `json:"type"`
	DepartmentID    string      `json:"departmentId"`
	DepartmentTitle string      `json:"departmentTitle"`
}

type BaseUserInfoClaims struct {
	jwt.Claims
	BaseUserInfo
}
