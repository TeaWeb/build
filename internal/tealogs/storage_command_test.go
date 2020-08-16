package tealogs

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestCommandStorage_Write(t *testing.T) {
	php, err := exec.LookPath("php")
	if err != nil { // not found php, so we can not test
		t.Log("php:", err)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	before := time.Now()

	storage := &CommandStorage{
		Storage: Storage{
			Format:   StorageFormatTemplate,
			Template: "${requestMethod} ${requestPath}",
		},
		Command: php,
		Args:    []string{cwd + "/tests/command_storage.php"},
	}
	err = storage.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = storage.Write([]*accesslogs.AccessLog{
		{
			RequestMethod: "GET",
			RequestPath:   "/hello",
		},
		{
			RequestMethod: "GET",
			RequestPath:   "/world",
		},
		{
			RequestMethod: "GET",
			RequestPath:   "/lu",
		},
		{
			RequestMethod: "GET",
			RequestPath:   "/ping",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = storage.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(time.Since(before).Seconds(), "seconds")
}
