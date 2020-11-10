package teaweb

import (
	"compress/gzip"
	_ "github.com/TeaWeb/build/internal/teacache"
	_ "github.com/TeaWeb/build/internal/teacluster"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/TeaWeb/build/internal/teautils"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/agents"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/agents/apps"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/agents/board"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/agents/cluster"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/agents/groups"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/agents/notices"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/agents/settings"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/api/agent"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/api/monitor"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/api/v1"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/cache"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/dashboard"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/index"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/install"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/log"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/login"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/logout"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/mongo"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/notices"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/plugins"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/backend"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/board"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/cache"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/certs"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/fastcgi"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/groups"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/headers"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/access"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/backends"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/websocket"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/log"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/log/policies"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/notices"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/rewrite"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/servers"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/settings"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/ssl"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/stat"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/tunnel"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/waf"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/search"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/settings"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/settings/backup"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/settings/cluster"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/settings/database"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/settings/login"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/settings/mongo"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/settings/mysql"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/settings/postgres"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/settings/profile"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/settings/server"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/settings/update"
	_ "github.com/TeaWeb/build/internal/teaweb/actions/default/ui"
	"github.com/TeaWeb/build/internal/teaweb/cmd"
	"github.com/TeaWeb/build/internal/teaweb/utils"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/sessions"
	"net/http"
	"time"
)

// 启动
var server *TeaGo.Server

func Start() {
	// 命令行
	shell := cmd.NewWebShell()
	shell.Start(server)
	if shell.ShouldStop {
		return
	}

	// 设置资源限制
	teautils.SetSuitableRLimit()

	// 日志
	writer := new(utils.LogWriter)
	writer.Init()
	logs.SetWriter(writer)

	// 启动代理
	go func() {
		time.Sleep(1 * time.Second)

		// 启动代理
		err := teaproxy.SharedManager.Start()
		if err != nil {
			logs.Error(err)
		}
		teaproxy.SharedManager.Wait()
	}()

	// 启动测试服务器
	if Tea.IsTesting() {
		go func() {
			time.Sleep(1 * time.Second)

			teatesting.StartTestServer()
		}()
	}

	// 设置变量
	Tea.SetPublicDir(teautils.WebRoot() + Tea.DS + "public")
	Tea.SetViewsDir(teautils.WebRoot() + Tea.DS + "views")
	Tea.SetTmpDir(teautils.WebRoot() + Tea.DS + "tmp")

	// 启动管理界面
	server = TeaGo.NewServer(false).
		ReadHeaderTimeout(3*time.Second).
		ReadTimeout(60*time.Second).
		AccessLog(false).

		Get("/", new(index.IndexAction)).
		Get("/logout", new(logout.IndexAction)).
		Get("/css/semantic.min.css", func(req *http.Request, writer http.ResponseWriter) {
			compressResource(writer, Tea.PublicDir()+"/css/semantic.min.css", "text/css; charset=utf-8")
		}).
		Get("/js/echarts.min.js", func(req *http.Request, writer http.ResponseWriter) {
			compressResource(writer, Tea.PublicDir()+"/js/echarts.min.js", "text/javascript; charset=utf-8")
		}).
		Get("/js/vue.min.js", func(req *http.Request, writer http.ResponseWriter) {
			compressResource(writer, Tea.PublicDir()+"/js/vue.min.js", "text/javascript; charset=utf-8")
		}).
		Get("/js/vue.tea.js", func(req *http.Request, writer http.ResponseWriter) {
			compressResource(writer, Tea.PublicDir()+"/js/vue.tea.js", "text/javascript; charset=utf-8")
		}).
		Get("/js/vue.components.js", func(req *http.Request, writer http.ResponseWriter) {
			compressResource(writer, Tea.PublicDir()+"/js/vue.components.js", "text/javascript; charset=utf-8")
		}).
		Get("/js/vue.js", func(req *http.Request, writer http.ResponseWriter) {
			compressResource(writer, Tea.PublicDir()+"/js/vue.js", "text/javascript; charset=utf-8")
		}).
		Get("/js/sortable.min.js", func(req *http.Request, writer http.ResponseWriter) {
			compressResource(writer, Tea.PublicDir()+"/js/sortable.min.js", "text/javascript; charset=utf-8")
		}).

		EndAll().

		Session(sessions.NewFileSessionManager(
			86400,
			"gSeDQJJ67tAVdnguDAQdGmnDVrjFd2I9",
		))
	server.Start()
}

// 压缩Javascript、CSS等静态资源
func compressResource(writer http.ResponseWriter, path string, mimeType string) {
	cssFile := files.NewFile(path)
	data, err := cssFile.ReadAll()
	if err != nil {
		return
	}

	gzipWriter, err := gzip.NewWriterLevel(writer, 5)
	if err != nil {
		_, err := writer.Write(data)
		if err != nil {
			logs.Error(err)
		}
		return
	}
	defer func() {
		err = gzipWriter.Close()
		if err != nil {
			logs.Error(err)
		}
	}()

	header := writer.Header()
	header.Set("Content-Encoding", "gzip")
	header.Set("Transfer-Encoding", "chunked")
	header.Set("Vary", "Accept-Encoding")
	header.Set("Accept-encoding", "gzip, deflate, br")
	header.Set("Content-Type", mimeType)
	header.Set("Last-Modified", "Sat, 02 Mar 2015 09:31:16 GMT")

	_, err = gzipWriter.Write(data)
	if err != nil {
		logs.Error(err)
	}
}
