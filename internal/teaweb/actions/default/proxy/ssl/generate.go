package ssl

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"time"
)

type GenerateAction actions.Action

// 证书生成程序
// 参考lego：https://go-acme.github.io/lego/usage/library/
func (this *GenerateAction) RunGet(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}
	this.Data["selectedTab"] = "https"
	this.Data["server"] = server

	this.Data["tasks"] = []interface{}{}
	if server.SSL != nil && len(server.SSL.CertTasks) > 0 {
		this.Data["tasks"] = lists.Map(server.SSL.CertTasks, func(k int, v interface{}) interface{} {
			task := v.(*teaconfigs.SSLCertTask)
			date := task.Request.CertDate()
			return maps.Map{
				"id":        task.Id,
				"email":     task.Request.User.Email,
				"domains":   task.Request.Domains,
				"runTime":   timeutil.Format("Y-m-d H:i:s", time.Unix(task.RunAt, 0)),
				"dayFrom":   date[0],
				"dayTo":     date[1],
				"runError":  task.RunError,
				"isExpired": len(date[1]) > 0 && timeutil.Format("Y-m-d") > date[1],
			}
		})
	}

	users := teaconfigs.SharedACMELocalUserList().Users
	if len(users) > 0 {
		this.Data["users"] = users
	} else {
		this.Data["users"] = []*teaconfigs.ACMELocalUser{}
	}

	this.Show()
}
