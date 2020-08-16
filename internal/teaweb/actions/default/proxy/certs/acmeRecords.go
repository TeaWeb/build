package certs

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"strings"
)

type AcmeRecordsAction actions.Action

// ACME记录
func (this *AcmeRecordsAction) RunPost(params struct {
	UserId  string
	Domains string
}) {
	user := teaconfigs.SharedACMELocalUserList().FindUser(params.UserId)
	if user == nil {
		this.Fail("找不到用户")
	}

	req := teaconfigs.NewACMERequest()
	req.User = user
	req.Domains = strings.Split(params.Domains, ",")
	client, err := req.Client()
	if err != nil {
		this.Fail("请求失败：" + err.Error())
	}

	records, err := req.RetrieveDNSRecords(client)
	if err != nil {
		this.Fail("请求失败：" + err.Error())
	}
	if len(records) == 0 {
		this.Fail("获取失败，请检查域名")
	}

	this.Data["records"] = records

	this.Success()
}
