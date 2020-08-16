package tealogs

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teatesting"
	"testing"
	"time"
)

func TestESStorage_Write(t *testing.T) {
	if !teatesting.RequireElasticSearch() {
		return
	}

	storage := &ESStorage{
		Storage: Storage{
		},
		Endpoint:    "http://127.0.0.1:9200",
		Index:       "logs",
		MappingType: "accessLogs",
		Username:    "hello",
		Password:    "world",
	}
	err := storage.Start()
	if err != nil {
		t.Fatal(err)
	}

	{
		storage.Format = StorageFormatJSON
		storage.Template = `${timeLocal} "${requestMethod} ${requestPath}"`
		err = storage.Write([]*accesslogs.AccessLog{
			{
				RequestMethod: "POST",
				RequestPath:   "/1",
				TimeLocal:     time.Now().Format("2/Jan/2006:15:04:05 -0700"),
				TimeISO8601:   "2018-07-23T22:23:35+08:00",
				Header: map[string][]string{
					"Content-Type": {"text/html"},
				},
			},
			{
				RequestMethod: "GET",
				RequestPath:   "/2",
				TimeLocal:     time.Now().Format("2/Jan/2006:15:04:05 -0700"),
				TimeISO8601:   "2018-07-23T22:23:35+08:00",
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
