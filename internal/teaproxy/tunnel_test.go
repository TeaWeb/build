package teaproxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestTunnel_Start(t *testing.T) {
	tunnel := NewTunnel(&teaconfigs.TunnelConfig{
		On:       true,
		Endpoint: "0.0.0.0:9001",
	})

	go func() {
		for {
			logs.Println(len(tunnel.connections))
			if len(tunnel.connections) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			req, _ := http.NewRequest(http.MethodGet, "http://hello.com/webhook", nil)
			req.Header.Set("User-Agent", "Tunnel-Call")

			resp, err := tunnel.Write(req)
			if err != nil {
				if err != io.EOF && err != io.ErrUnexpectedEOF {
					logs.Error(err)
				}
				continue
			} else {
				logs.Println(resp)
				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					logs.Error(err)
				} else {
					logs.Println(string(data))
				}
				_ = resp.Body.Close()
			}

			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		if teatesting.IsGlobal() {
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(10 * time.Second)
		}
		err := tunnel.Close()
		if err != nil {
			logs.Error(err)
		}
	}()

	t.Log(tunnel.Start())
}
