package tealogs

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/Tea"
	"testing"
	"time"
)

func TestFileStorage_Write(t *testing.T) {
	storage := &FileStorage{
		Storage: Storage{
		},
		Path: Tea.Root + "/logs/access-${date}.log",
	}
	err := storage.Start()
	if err != nil {
		t.Fatal(err)
	}

	{
		err = storage.Write([]*accesslogs.AccessLog{
			{
				RequestPath: "/hello",
			},
			{
				RequestPath: "/world",
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err = storage.Write([]*accesslogs.AccessLog{
			{
				RequestPath: "/1",
			},
			{
				RequestPath: "/2",
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		storage.Format = StorageFormatTemplate
		storage.Template = `${timeLocal} "${requestMethod} ${requestPath}" ${log}`
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
