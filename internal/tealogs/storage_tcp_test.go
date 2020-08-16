package tealogs

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teatesting"
	"net"
	"testing"
	"time"
)

func TestTCPStorage_Write(t *testing.T) {
	if !teatesting.RequirePort(9981) {
		return
	}

	go func() {
		server, err := net.Listen("tcp", "127.0.0.1:9981")
		if err != nil {
			t.Fatal(err)
		}
		for {
			conn, err := server.Accept()
			if err != nil {
				break
			}

			buf := make([]byte, 1024)
			for {
				n, err := conn.Read(buf)
				if n > 0 {
					t.Log(string(buf[:n]))
				}
				if err != nil {
					break
				}
			}
			break
		}
		_ = server.Close()
	}()

	storage := &TCPStorage{
		Storage: Storage{
		},
		Network: "tcp",
		Addr:    "127.0.0.1:9981",
	}
	err := storage.Start()
	if err != nil {
		t.Fatal(err)
	}

	{
		storage.Format = StorageFormatTemplate
		storage.Template = `${timeLocal} "${requestMethod} ${requestPath}"`
		err = storage.Write([]*accesslogs.AccessLog{
			{
				RequestMethod: "POST",
				RequestPath:   "/1",
				TimeLocal:     time.Now().Format("2/Jan/2006:15:04:05 -0700"),
			},
			{
				RequestMethod: "GET",
				RequestPath:   "/2",
				TimeLocal:     time.Now().Format("2/Jan/2006:15:04:05 -0700"),
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	err = storage.Close()
	if err != nil {
		t.Fatal(err)
	}
}
