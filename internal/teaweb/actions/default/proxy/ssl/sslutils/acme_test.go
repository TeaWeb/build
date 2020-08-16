package sslutils

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teatesting"
	"strings"
	"testing"
)

func TestReloadACMECert(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	server := &teaconfigs.ServerConfig{
		Id: "abc",
	}
	server.SSL = teaconfigs.NewSSLConfig()
	server.SSL.Certs = []*teaconfigs.SSLCertConfig{
		{
			On:       true,
			TaskId:   "123",
			CertFile: "123.pem",
			KeyFile:  "456.key",
		},
	}
	teaproxy.SharedManager.ApplyServer(server)
	errs := ReloadACMECert("abc", "123")
	for _, err := range errs {
		if strings.Contains(err.Error(), "failed:open") { // 符合预期
			continue
		}
		t.Fatal(err)
	}
}

func TestRenewACMECerts(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	RenewACMECerts()
}
