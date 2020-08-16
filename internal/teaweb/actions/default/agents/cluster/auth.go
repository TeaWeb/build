package cluster

import "github.com/iwind/TeaGo/actions"

type AuthAction actions.Action

// 认证
func (this *AuthAction) Run(params struct {
	Master   string
	Dir      string
	AuthType string
	Username string
	Password string
	Key      *actions.File
	Must     *actions.Must
}) {
	params.Must.
		Field("master", params.Master).
		Require("请输入TeaWeb访问地址").
		Field("dir", params.Dir).
		Require("请输入安装目录").
		Field("username", params.Username).
		Require("请输入登录主机的用户名")

	this.Data["key"] = ""

	if params.AuthType == "password" {
		params.Must.Field("password", params.Password).
			Require("请输入登录主机的密码")
	} else {
		if params.Key == nil {
			this.FailField("key", "请选择密钥文件")
		}

		data, err := params.Key.Read()
		if err != nil {
			this.FailField("key", "密钥读取失败："+err.Error())
		}
		this.Data["key"] = string(data)
	}

	this.Success()
}
