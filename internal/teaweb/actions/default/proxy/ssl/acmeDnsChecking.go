package ssl

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/rands"
	"net"
	"strings"
	"time"
)

type AcmeDnsCheckingAction actions.Action

// ACME DNS检查
func (this *AcmeDnsCheckingAction) RunPost(params struct {
	ServerId string
	UserId   string
	Domains  string
	Records  string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	user := teaconfigs.SharedACMELocalUserList().FindUser(params.UserId)
	if user == nil {
		this.Fail("找不到用户信息")
	}

	records := []*teaconfigs.ACMEDNSRecord{}
	err := json.Unmarshal([]byte(params.Records), &records)
	if err != nil {
		this.Fail("无法解析DNS记录：" + err.Error())
	}

	if len(records) == 0 {
		this.Fail("DNS记录为空")
	}

	for _, record := range records {
		values, err := net.LookupTXT(record.FQDN)
		if err != nil {
			this.Fail("DNS解析失败：" + err.Error())
		}

		if !lists.ContainsString(values, record.Value) {
			this.Fail("DNS解析不匹配，域名" + record.FQDN + "，现在解析值：" + strings.Join(values, ", ") + " 期望解析值：" + record.Value + "。如果确认已经设置正确，请试着清除本地DNS缓存再重试。")
		}
	}

	// 获取证书
	req := teaconfigs.NewACMERequest()
	req.User = user
	req.Domains = strings.Split(params.Domains, ",")
	req.DNSRecords = records

	client, err := req.Client()
	if err != nil {
		this.Fail("证书获取失败：" + err.Error())
	}

	err = req.Retrieve(client)
	if err != nil {
		this.Fail("证书获取失败：" + err.Error())
	}

	// write cert & key
	certFile := "ssl." + rands.HexString(16) + ".pem"
	keyFile := "ssl." + rands.HexString(16) + ".key"
	err = req.WriteCertFile(Tea.ConfigFile(certFile))
	if err != nil {
		this.Fail("证书保存失败：" + err.Error())
	}

	err = req.WriteKeyFile(Tea.ConfigFile(keyFile))
	if err != nil {
		this.Fail("密钥保存失败：" + err.Error())
	}

	if server.SSL == nil {
		server.SSL = &teaconfigs.SSLConfig{
			On: false,
		}
	}

	task := teaconfigs.NewSSLCertTask()
	task.RunAt = time.Now().Unix()
	task.Request = req

	cert := teaconfigs.NewSSLCertConfig(certFile, keyFile)
	cert.TaskId = task.Id
	cert.Description = "通过ACME生成的证书"
	server.SSL.AddCert(cert)
	server.SSL.AddCertTask(task)

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	proxyutils.NotifyChange()

	this.Success()
}
