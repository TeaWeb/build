package ssl

import (
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/ssl/sslutils"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
	"time"
)

func init() {
	// 路由定义
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(proxy.Helper)).
			Module("").
			Prefix("/proxy/ssl").
			Get("", new(IndexAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/startHttps", new(StartHttpsAction)).
			Post("/shutdownHttps", new(ShutdownHttpsAction)).
			Get("/downloadFile", new(DownloadFileAction)).
			Get("/generate", new(GenerateAction)).
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
			Post("/makeShared", new(MakeSharedAction)).
			EndAll()
	})

	// 检查ACME证书更新
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		teautils.Every(24*time.Hour, func(ticker *teautils.Ticker) {
			sslutils.RenewACMECerts()
		})
	})
}
