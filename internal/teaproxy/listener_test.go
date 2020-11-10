package teaproxy

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestListener_AddServer(t *testing.T) {
	listener := NewListener()
	listener.Scheme = SchemeHTTP
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web001"
		server.AddName("a.com")
		listener.ApplyServer(server)
	}
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web002"
		server.AddName("b.com")
		listener.ApplyServer(server)
	}
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web003"
		server.AddName("b.com")
		server.On = true
		listener.ApplyServer(server)
	}
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web004"
		server.AddName("b.com")
		server.On = true
		server.Http = true
		listener.ApplyServer(server)
	}
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web005"
		server.AddName("b.com")
		server.On = true
		server.Http = false
		server.SSL = &teaconfigs.SSLConfig{
			On: true,
		}
		listener.ApplyServer(server)
	}
	printListener(listener, t)
}

func TestListener_RemoveServer(t *testing.T) {
	listener := NewListener()
	listener.Scheme = SchemeHTTP
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web004"
		server.AddName("b.com")
		server.On = true
		server.Http = true
		listener.ApplyServer(server)
	}
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web001"
		server.AddName("a.com")
		server.On = true
		server.Http = true
		listener.ApplyServer(server)
	}
	listener.RemoveServer("web001")
	printListener(listener, t)
}

func TestListener_Start(t *testing.T) {
	listener := NewListener()
	listener.Scheme = SchemeHTTP
	listener.Address = "127.0.0.1:8881"
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web001"
		server.AddName("a.com")
		server.Http = true
		server.AddName("wx.teaos.cn")
		server.AddBackend(&teaconfigs.BackendConfig{
			On:      true,
			Id:      "backend001",
			Address: "127.0.0.1:9991",
			Weight:  10,
		})
		err := server.Validate()
		if err != nil {
			t.Fatal(err)
		}
		listener.ApplyServer(server)
	}
	go func() {
		time.Sleep(5 * time.Second)
		err := listener.Shutdown()
		if err != nil {
			t.Log(err)
		}
		logs.Println("shutdown")
	}()
	err := listener.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListener_Reload_RemoveServer(t *testing.T) {
	listener := NewListener()
	listener.Scheme = SchemeHTTP
	listener.Address = "127.0.0.1:8881"
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web001"
		server.AddName("a.com")
		server.Http = true
		server.AddName("wx.teaos.cn")
		server.AddBackend(&teaconfigs.BackendConfig{
			On:      true,
			Id:      "backend001",
			Address: "127.0.0.1:9991",
			Weight:  10,
		})
		server.AddResponseHeader(&shared.HeaderConfig{
			On:     true,
			Name:   "Backend",
			Value:  "${backend.id}",
			Always: true,
		})
		err := server.Validate()
		if err != nil {
			t.Fatal(err)
		}
		listener.ApplyServer(server)
	}
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web002"
		server.AddName("a.com")
		server.Http = true
		server.AddName("wx2.teaos.cn")
		server.AddBackend(&teaconfigs.BackendConfig{
			On:      true,
			Id:      "backend002",
			Address: "127.0.0.1:9991",
			Weight:  10,
		})
		server.AddResponseHeader(&shared.HeaderConfig{
			On:     true,
			Name:   "Backend",
			Value:  "${backend.id}",
			Always: true,
		})
		err := server.Validate()
		if err != nil {
			t.Fatal(err)
		}
		listener.ApplyServer(server)
	}
	go func() {
		time.Sleep(5 * time.Second)

		listener.RemoveServer("web001")

		err := listener.Reload()
		if err != nil {
			t.Log(err)
		}
		logs.Println("reload")

		if teatesting.IsGlobal() {
			_ = listener.Shutdown()
		}
	}()
	err := listener.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListener_Reload_ChangeServer(t *testing.T) {
	listener := NewListener()
	listener.Scheme = SchemeHTTP
	listener.Address = "127.0.0.1:8881"
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web001"
		server.AddName("a.com")
		server.Http = true
		server.AddName("wx.teaos.cn")
		server.AddBackend(&teaconfigs.BackendConfig{
			On:      true,
			Id:      "backend001",
			Address: "127.0.0.1:9991",
			Weight:  10,
		})
		server.AddResponseHeader(&shared.HeaderConfig{
			On:     true,
			Name:   "Backend",
			Value:  "${backend.id}",
			Always: true,
		})
		err := server.Validate()
		if err != nil {
			t.Fatal(err)
		}
		listener.ApplyServer(server)
	}
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web002"
		server.AddName("a.com")
		server.Http = true
		server.AddName("wx2.teaos.cn")
		server.AddBackend(&teaconfigs.BackendConfig{
			On:      true,
			Id:      "backend002",
			Address: "127.0.0.1:9991",
			Weight:  10,
		})
		server.AddResponseHeader(&shared.HeaderConfig{
			On:     true,
			Name:   "Backend",
			Value:  "${backend.id}",
			Always: true,
		})
		err := server.Validate()
		if err != nil {
			t.Fatal(err)
		}
		listener.ApplyServer(server)
	}
	go func() {
		time.Sleep(5 * time.Second)

		{
			server := teaconfigs.NewServerConfig()
			server.Id = "web001"
			server.Http = true
			server.AddBackend(&teaconfigs.BackendConfig{
				On:      true,
				Id:      "backend003",
				Address: "wx2.teaos.cn",
				Weight:  10,
			})
			server.AddResponseHeader(&shared.HeaderConfig{
				On:     true,
				Name:   "Backend",
				Value:  "${backend.id}",
				Always: true,
			})
			err := server.Validate()
			if err != nil {
				t.Fatal(err)
			}
			listener.ApplyServer(server)
		}

		err := listener.Reload()
		if err != nil {
			t.Log(err)
		}
		logs.Println("reload")

		// shutdown
		if teatesting.IsGlobal() {
			err = listener.Shutdown()
			if err != nil {
				t.Log(err)
			}
		}
	}()
	err := listener.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkListener_FindNamedServer(b *testing.B) {
	listener := NewListener()
	listener.Scheme = SchemeHTTP
	listener.Address = "127.0.0.1:8881"
	for i := 0; i < 10; i++ {
		server := teaconfigs.NewServerConfig()
		server.Id = "web00" + fmt.Sprintf("%d", i)
		server.Http = true
		server.AddName("a.com")
		server.AddName("teaos" + fmt.Sprintf("%d", i) + ".cn")
		err := server.Validate()
		if err != nil {
			b.Fatal(err)
		}
		listener.ApplyServer(server)
	}
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web00" + fmt.Sprintf("%d", 20000)
		server.Http = true
		server.AddName("*.teaos.cn")
		err := server.Validate()
		if err != nil {
			b.Fatal(err)
		}
		listener.ApplyServer(server)
	}
	listener.currentServers = listener.servers

	for i := 0; i < b.N; i++ {
		listener.findNamedServer("teaos5.cn")
		listener.findNamedServer("wx.teaos.cn")
	}
}

func printListener(listener *Listener, t *testing.T) {
	for _, s := range listener.servers {
		logs.PrintAsJSON(maps.Map{
			"id":          s.Id,
			"name":        s.Name,
			"http":        s.Http,
			"ssl":         s.SSL,
			"description": s.Description,
		}, t)
	}
}

func TestDetectIPOrDomain(t *testing.T) {
	a := assert.NewAssertion(t)
	a.IsTrue(testIsIP("192.168.1.101"))
	a.IsTrue(testIsIP("192.168.1.102:1000"))
	a.IsTrue(testIsIP("[1:2:3:4]:1000"))
	a.IsTrue(!testIsIP("192.168.1.com"))
	a.IsTrue(!testIsIP("www.teaos.cn"))
	a.IsTrue(!testIsIP("192.168.1.com:12345"))
	a.IsTrue(!testIsIP("local345:12345"))
}

func BenchmarkIPOrDomain(b *testing.B) {
	runtime.GOMAXPROCS(1)

	for i := 0; i < b.N; i++ {
		_ = testIsIP("192.168.1.101")
		_ = testIsIP("192.168.1.101:1000")
		_ = testIsIP("www.example.com")
		_ = testIsIP("www.example.com:12345")
	}
}

func testIsIP(host string) bool {
	// IPv6
	if strings.Index(host, "[") > -1 {
		return true
	}

	for _, b := range host {
		if b >= 'a' && b <= 'z' {
			return false
		}
	}

	return true
}
