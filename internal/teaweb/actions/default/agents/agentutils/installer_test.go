package agentutils

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"testing"
	"time"
)

func TestInstaller_Start(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	installer := NewInstaller()
	installer.Host = "192.168.2.33"
	installer.Port = 22
	installer.Timeout = 5 * time.Second
	installer.AuthUsername = "root"
	installer.AuthPassword = "123456"
	installer.Master = "http://192.168.2.43:7777"
	installer.Dir = "/opt/teaweb/"

	err := installer.Start()
	if err != nil {
		// log
		for _, l := range installer.Logs {
			t.Log(l)
		}

		t.Fatal(err)
	}

	t.Log("hostName:", installer.HostName)
	t.Log("os:", installer.OS)
	t.Log("arch", installer.Arch)

	// log
	for _, l := range installer.Logs {
		t.Log(l)
	}
}

func TestInstaller_Start_KnownHost(t *testing.T) {

}
