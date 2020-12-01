package login

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/audits"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/types"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"net/http"
	"time"
)

type IndexAction actions.Action

var TokenSalt = rands.HexString(32)

// 登录
func (this *IndexAction) RunGet() {
	// 检查IP限制
	if !configs.SharedAdminConfig().AllowIP(this.RequestRemoteIP()) {
		this.ResponseWriter.WriteHeader(http.StatusForbidden)
		this.WriteString(teaconst.TeaProductName + " Access Forbidden")
		return
	}

	b := Notify(this)
	if !b {
		return
	}

	this.Data["teaDemoEnabled"] = teaconst.DemoEnabled

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	this.Data["token"] = stringutil.Md5(TokenSalt+timestamp) + timestamp

	this.Show()
}

// 提交登录
func (this *IndexAction) RunPost(params struct {
	Username string
	Password string
	Token    string
	Remember bool
	Must     *actions.Must
	Auth     *helpers.UserShouldAuth
}) {
	// 检查IP限制
	if !configs.SharedAdminConfig().AllowIP(this.RequestRemoteIP()) {
		this.ResponseWriter.WriteHeader(http.StatusForbidden)
		this.WriteString(teaconst.TeaProductName + " Access Forbidden")
		return
	}

	b := Notify(this)
	if !b {
		return
	}

	params.Must.
		Field("username", params.Username).
		Require("请输入用户名").
		Field("password", params.Password).
		Require("请输入密码")

	if params.Password == stringutil.Md5("") {
		this.FailField("password", "请输入密码")
	}

	// 检查token
	if len(params.Token) <= 32 {
		this.log(params.Username, false)
		this.Fail("请通过登录页面登录")
	}
	timestampString := params.Token[32:]
	if stringutil.Md5(TokenSalt+timestampString) != params.Token[:32] {
		this.log(params.Username, false)
		this.FailField("refresh", "登录页面已过期，请刷新后重试")
	}
	timestamp := types.Int64(timestampString)
	if timestamp < time.Now().Unix()-1800 {
		this.log(params.Username, false)
		this.FailField("refresh", "登录页面已过期，请刷新后重试")
	}

	// 查找用户
	adminConfig := configs.SharedAdminConfig()
	user := adminConfig.FindActiveUser(params.Username)
	if user != nil {
		// 错误次数
		if user.CountLoginTries() >= 3 {
			this.log(params.Username, false)
			this.Fail("登录失败已超过3次，系统被锁定，需要重置服务后才能继续")
		}

		// 密码错误
		if !adminConfig.ComparePassword(params.Password, user.Password) {
			user.IncreaseLoginTries()
			this.log(params.Username, false)
			this.Fail("登录失败，请检查用户名密码")
		}

		user.ResetLoginTries()

		// Session
		params.Auth.StoreUsername(user.Username, params.Remember)

		// 记录登录IP
		user.LoggedAt = time.Now().Unix()
		user.LoggedIP = this.RequestRemoteIP()

		// 在开发环境下不保存登录IP，以便于不干扰git
		if !Tea.IsTesting() {
			err := adminConfig.Save()
			if err != nil {
				logs.Error(err)
			}
		}

		this.log(params.Username, true)
		this.Next("/", nil, "").Success()
		return
	}

	this.log(params.Username, false)
	this.Fail("登录失败，请检查用户名密码")
}

func (this *IndexAction) log(username string, success bool) {
	go func() {
		var message string
		if success {
			message = "登录成功"
		} else {
			message = "登录失败"
		}
		err := teadb.AuditLogDAO().InsertOne(audits.NewLog(username, audits.ActionLogin, message, map[string]string{
			"ip": this.RequestRemoteIP(),
		}))
		if err != nil {
			logs.Error(err)
		}
	}()
}
