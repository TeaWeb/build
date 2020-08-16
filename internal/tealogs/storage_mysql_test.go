package tealogs

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teatesting"
	"testing"
	"time"
)

func TestMySQLStorage_Write(t *testing.T) {
	if !teatesting.RequireMySQL() {
		return
	}

	before := time.Now()
	defer func() {
		t.Log("cost:", time.Since(before).Seconds(), "seconds")
	}()

	storage := &MySQLStorage{
		Storage: Storage{
		},
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "teaweb",
		Table:    "accessLogs${date}",
		LogField: "log",
	}

	err := storage.Start()
	if err != nil {
		t.Fatal(err)
	}

	{
		storage.Format = StorageFormatJSON
		storage.Template = `${timeLocal} "${requestMethod} ${requestPath}"`
		err := storage.Write([]*accesslogs.AccessLog{
			{
				RequestMethod: "POST",
				RequestPath:   "/1",
				TimeLocal:     time.Now().Format("2/Jan/2006:15:04:05 -0700"),
				Header: map[string][]string{
					"Content-Type": {"text/html"},
				},
			},
			{
				RequestMethod: "GET",
				RequestPath:   "/2",
				TimeLocal:     time.Now().Format("2/Jan/2006:15:04:05 -0700"),
				Header: map[string][]string{
					"Content-Type": {"text/css"},
				},
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
