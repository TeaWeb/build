package locations

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
	"gopkg.in/yaml.v3"
)

type ImportAction actions.Action

// 导入路径规则
func (this *ImportAction) RunGet(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = server
	this.Data["selectedTab"] = "location"
	this.Data["selectedSubTab"] = "detail"

	this.Show()
}

// 提交保存
func (this *ImportAction) RunPost(params struct {
	ServerId string
	File     *actions.File
	Must     *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if params.File == nil {
		this.Fail("请上传要导入的路径规则文件")
	}

	data, err := params.File.Read()
	if err != nil {
		this.Fail("文件读取失败：" + err.Error())
	}

	location := &teaconfigs.LocationConfig{}
	err = yaml.Unmarshal(data, location)
	if err != nil {
		this.Fail("文件解析失败：" + err.Error())
	}
	location.Id = rands.HexString(16)

	server.AddLocation(location)
	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	proxyutils.NotifyChange()

	this.Success()
}
