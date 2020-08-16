package agent

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestUpgradeAction_Run(t *testing.T) {
	if !teatesting.RequireTeaWebServer() {
		return
	}

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:7777/api/agent/upgrade", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Tea-Agent-Id", "4iloHPyb9hNvi1cR")
	req.Header.Set("Tea-Agent-Key", "v348p0bIE9R1V1wGZ84wVmhCH9hSXQIF")
	req.Header.Set("Tea-Agent-Os", "linux")
	//req.Header.Set("Tea-Agent-Os", "windows")
	req.Header.Set("Tea-Agent-Arch", "amd64")
	req.Header.Set("Tea-Agent-Version", "0.0.11")

	t.Log(req.Header)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatal("status code wrong:", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) < 1024 {
		t.Log(string(data))
		return
	}

	t.Log("data length", len(data))
}
