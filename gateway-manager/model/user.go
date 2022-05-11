package model

import (
	. "github.com/gohutool/boot4go-fastjson"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : user.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/24 20:24
* 修改历史 : 1. [2022/4/24 20:24] 创建文件 by LongYong
*/

type AdminUser struct {
	UserId   string
	UserName string
	Password string
	Salt     string
}

func (d *AdminUser) Unmarshal(value *Value) error {
	d.UserId = value.GetString("UserId")
	d.UserName = value.GetString("UserName")
	d.Password = value.GetString("Password")
	d.Salt = value.GetString("Salt")
	return nil
}
