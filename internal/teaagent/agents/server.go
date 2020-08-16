package teaagents

import (
	"github.com/TeaWeb/build/internal/teaagent/agentconst"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

// Agent提供的Web Server，可以用来读取状态、进行控制等
type Server struct {
	Addr string

	server *http.Server
}

// 获取新服务对象
func NewServer() *Server {
	return &Server{
		Addr: "127.0.0.1:7778",
	}
}

// 启动
func (this *Server) Start() error {
	// 如果没启动，则启动
	handler := http.NewServeMux()
	handler.HandleFunc("/status", func(writer http.ResponseWriter, req *http.Request) {
		writer.Write([]byte(stringutil.JSONEncode(maps.Map{
			"version": agentconst.AgentVersion,
		})))
	})

	this.server = &http.Server{
		Addr:        this.Addr,
		Handler:     handler,
		IdleTimeout: 2 * time.Minute,
	}
	err := this.server.ListenAndServe()
	return err
}

// 关闭
func (this *Server) Shutdown() {
	if this.server != nil {
		this.server.Shutdown(context.Background())
		this.server = nil
	}
}
