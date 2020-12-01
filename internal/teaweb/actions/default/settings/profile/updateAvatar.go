package profile

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/rands"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type UpdateAvatarAction actions.Action

func (this *UpdateAvatarAction) Run(params struct{}) {
	username := this.Session().GetString("username")
	user := configs.SharedAdminConfig().FindActiveUser(username)

	this.Data["user"] = map[string]interface{}{
		"avatar": user.Avatar,
	}

	this.Show()
}

func (this *UpdateAvatarAction) RunPost(params struct {
	AvatarFile *actions.File
}) {
	if params.AvatarFile == nil {
		this.Fail("请选择要上传的头像文件")
	}

	if !lists.ContainsString([]string{".jpg", ".jpeg", ".png", ".gif"}, params.AvatarFile.Ext) {
		this.Fail("上传的图片文件格式不正确")
	}

	reader, err := params.AvatarFile.OriginFile.Open()
	if err != nil {
		this.Fail("上传的图片文件格式不正确")
	}
	defer reader.Close()
	_, _, err = image.DecodeConfig(reader)
	if err != nil {
		this.Fail("上传的图片文件格式不正确")
	}

	dir := files.NewFile(Tea.ConfigDir() + "/avatars")
	if !dir.Exists() {
		dir.Mkdir()
	}

	rand := rands.HexString(16)
	_, err = params.AvatarFile.WriteToPath(dir.Path() + "/" + rand + params.AvatarFile.Ext)
	if err != nil {
		this.Fail("头像文件写入失败，请检查文件权限")
	}

	username := this.Session().GetString("username")
	adminConfig := configs.SharedAdminConfig()
	user := adminConfig.FindActiveUser(username)
	user.Avatar = "/avatar/" + rand + params.AvatarFile.Ext
	adminConfig.Save()

	this.Refresh().Success("上传成功")
}
