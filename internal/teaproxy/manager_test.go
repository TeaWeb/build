package teaproxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"testing"
	"time"
)

func TestManager_Start(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	manager := NewManager()
	err := manager.Start()
	if err != nil {
		t.Fatal(err)
	}
	manager.Wait()
}

func TestManager_AddServer(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	manager := NewManager()
	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web001"

		// http
		server.Http = true
		server.AddListen("127.0.0.1")
		server.AddListen("8081")
		server.AddListen(":8080")

		// https
		server.SSL = &teaconfigs.SSLConfig{
			On:     true,
			Listen: []string{"192.168.1.100", "192.168.2.100:587"},
		}
		manager.ApplyServer(server)
	}

	{
		server := teaconfigs.NewServerConfig()
		server.Id = "web001"

		// http
		server.Http = true
		server.AddListen("127.0.0.1:8881")

		// backends
		server.AddBackend(&teaconfigs.BackendConfig{
			On:      true,
			Id:      "backend001",
			Address: "127.0.0.1:9991",
			Weight:  10,
		})
		server.Validate()
		manager.ApplyServer(server)

		// 修改端口测试
		go func() {
			time.Sleep(10 * time.Second)

			server := teaconfigs.NewServerConfig()
			server.Id = "web001"

			// http
			server.Http = true
			server.AddListen("127.0.0.1:8882")

			// backends
			server.AddBackend(&teaconfigs.BackendConfig{
				On:      true,
				Id:      "backend001",
				Address: "127.0.0.1:9991",
				Weight:  10,
			})
			server.Validate()
			manager.ApplyServer(server)
			manager.Reload()
		}()
	}

	// 关闭测试
	go func() {
		time.Sleep(20 * time.Second)
		logs.Println("shutdown")
		manager.Shutdown()
	}()

	for _, listener := range manager.listeners {
		t.Log("===" + listener.Address + "===")
		if len(listener.servers) == 0 {
			t.Log("no servers")
		} else {
			for _, s := range listener.servers {
				t.Log(s.Id)
			}
		}
	}

	manager.Reload()

	t.Log("=======================================")
	for _, listener := range manager.listeners {
		t.Log("===" + listener.Address + "===")
		if len(listener.servers) == 0 {
			t.Log("no servers")
		} else {
			for _, s := range listener.servers {
				t.Log(s.Id)
			}
		}
	}

	manager.Wait()
}
