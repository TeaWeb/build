package tealogs

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"testing"
)

func TestStorage_FormatAccessLogString(t *testing.T) {
	{
		storage := &Storage{
			Format: StorageFormatJSON,
		}
		t.Log(storage.FormatAccessLogString(&accesslogs.AccessLog{
			ServerName:  "hello.com",
			RequestPath: "/webhook",
			Args:        "a=1&b=2",
		}))
	}

	{
		storage := &Storage{
			Format:   StorageFormatTemplate,
			Template: "${serverName} - ${requestPath} - ${args}",
		}
		t.Log(storage.FormatAccessLogString(&accesslogs.AccessLog{
			ServerName:  "hello.com",
			RequestPath: "/webhook",
			Args:        "a=1&b=2",
		}))
	}
}

func TestStorage_FormatAccessLogBytes(t *testing.T) {
	{
		storage := &Storage{
			Format: StorageFormatJSON,
		}
		data, err := storage.FormatAccessLogBytes(&accesslogs.AccessLog{
			ServerName:  "hello.com",
			RequestPath: "/webhook",
			Args:        "a=1&b=2",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
	}

	{
		storage := &Storage{
			Format:   StorageFormatTemplate,
			Template: "${serverName} - ${requestPath} - ${args}",
		}

		data, err := storage.FormatAccessLogBytes(&accesslogs.AccessLog{
			ServerName:  "hello.com",
			RequestPath: "/webhook",
			Args:        "a=1&b=2",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
	}
}

func TestStorage_FormatVariables(t *testing.T) {
	storage := &Storage{}
	t.Log(storage.FormatVariables("/var/log/teaweb-year${year}-month${month}-week${week}-day${day}-hour${hour}-minute${minute}-second${second}-date${date}"))
}
