package certs

import (
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/certs/certutils"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
	"time"
)

func init() {
	// 路由设置
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(Helper)).
			Prefix("/proxy/certs").
			Get("", new(IndexAction)).
			GetPost("/upload", new(UploadAction)).
			Post("/delete", new(DeleteAction)).
			Get("/detail", new(DetailAction)).
			GetPost("/update", new(UpdateAction)).
			Get("/download", new(DownloadAction)).
			Get("/acme", new(AcmeAction)).
			Get("/acmeCreateTask", new(AcmeCreateTaskAction)).
			GetPost("/acmeCreateUser", new(AcmeCreateUserAction)).
			Get("/acmeUsers", new(AcmeUsersAction)).
			Post("/acmeUserDelete", new(AcmeUserDeleteAction)).
			Post("/acmeRecords", new(AcmeRecordsAction)).
			Post("/acmeDnsChecking", new(AcmeDnsCheckingAction)).
			Post("/acmeDeleteTask", new(AcmeDeleteTaskAction)).
			Post("/acmeRenewTask", new(AcmeRenewTaskAction)).
			Get("/acmeTask", new(AcmeTaskAction)).
			Get("/acmeDownload", new(AcmeDownloadAction)).
			EndAll()
	})

	// 检查ACME证书更新
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		teautils.Every(24*time.Hour, func(ticker *teautils.Ticker) { // TODO
			certutils.RenewACMECerts()
		})
	})
}
