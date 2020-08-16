package notices

import (
	"bytes"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"testing"
)

func TestNewItemFromRequest(t *testing.T) {
	b := bytes.NewBuffer([]byte("backendDownNoticeOn=1&backendDownNoticeLevel=2&backendDownNoticeSubject=subject123&backendDownNoticeBody=body456"))
	req, err := http.NewRequest(http.MethodPost, "/", b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	err = req.ParseForm()
	if err != nil {
		t.Fatal(err)
	}
	item := NewItemFromRequest(req, "backendDown")
	logs.PrintAsJSON(item, t)
}

func TestItem_FormatSubject(t *testing.T) {
	{
		item := NewItem(NoticeLevelWarning)
		item.Subject = "server '${server.address}' has error '${error}'"
		t.Log(item.FormatSubject(maps.Map{
			"server.address": "LB001",
			"error":          "not so good",
		}))
	}

	{
		item := NewItem(NoticeLevelWarning)
		item.Subject = "server '${ server.address }' has error '${error}'" // has spaces in var
		t.Log(item.FormatSubject(maps.Map{
			"server.address": "LB001",
			"error":          "not so good",
		}))
	}

	{
		item := NewItem(NoticeLevelWarning)
		item.Subject = "server '${ server.address }' has error '${error}'" // has spaces in var
		t.Log(item.FormatSubject(maps.Map{
			"server.address": "LB001",
		}))
	}
}

func TestItem_FormatBody(t *testing.T) {
	{
		item := NewItem(NoticeLevelWarning)
		item.Body = "server '${server.address}' has error '${error}'"
		t.Log(item.FormatBody(maps.Map{
			"server.address": "LB001",
			"error":          "not so good",
		}))
	}
}
