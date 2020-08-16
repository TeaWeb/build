package teaproxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teatesting"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFTPClient_Do(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	backend := &teaconfigs.BackendConfig{
		Address: "192.168.2.30:21",
		FTP: &teaconfigs.FTPBackendConfig{
			Username: "www",
			Password: "123456",
			Dir:      "",
		},
	}
	client := SharedFTPClientPool.client(nil, backend, nil)

	for _, file := range []string{"/index.html", "index.a", "/dir1/dir2/hello.txt"} {
		func() {
			req, err := http.NewRequest(http.MethodGet, file, nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Log(file+":", err.Error())
				return
			}
			defer func() {
				_ = resp.Body.Close()
			}()

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Log(file+":", err.Error())
				return
			}

			t.Log(file+":", string(data))
		}()
	}
}

func TestFTPClient_Do_ChangeDir(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	backend := &teaconfigs.BackendConfig{
		Address: "192.168.2.30:21",
		FTP: &teaconfigs.FTPBackendConfig{
			Username: "www",
			Password: "123456",
			Dir:      "/dir1/dir2",
		},
	}
	client := SharedFTPClientPool.client(nil, backend, nil)

	for _, file := range []string{"hello.txt", "hello1.txt", "hello2.txt"} {
		func() {
			req, err := http.NewRequest(http.MethodGet, file, nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Log(file+":", err.Error())
				return
			}
			defer func() {
				_ = resp.Body.Close()
			}()

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Log(file+":", err.Error())
				return
			}

			t.Log(file+":", string(data))
		}()
	}
}
