package forms

import (
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestForm_Compose(t *testing.T) {
	form := NewForm("user")
	group := form.NewGroup()
	{
		field := NewTextField("主机", "Host")
		field.MaxLength = 100
		field.Placeholder = "127.0.0.1"
		field.Value = "123"
		field.Size = 8
		field.Code = "host"
		group.Add(field)
	}

	form.Compose()

	t.Log(group.HTML)
}

func TestForm_Init(t *testing.T) {
	form := NewForm("user")
	group := form.NewGroup()
	{
		field := NewTextField("主机", "Host")
		field.MaxLength = 100
		field.Placeholder = "127.0.0.1"
		field.Value = "127.0.0.2"
		field.Size = 8
		field.Code = "host"
		group.Add(field)
	}

	{
		field := NewTextField("主机", "Host")
		field.MaxLength = 100
		field.Placeholder = "127.0.0.1"
		field.Value = "127.0.0.2"
		field.Size = 8
		field.Code = "host"
		group.Add(field)
	}

	form.Init(map[string]interface{}{
		"host": []string{"192.168.1.100", "192.168.1.101"},
	})
	form.Compose()

	t.Log(group.HTML)
}

func TestForm_Values(t *testing.T) {
	form := NewForm("user")
	group := form.NewGroup()
	{
		field := NewTextField("主机", "Host")
		field.MaxLength = 100
		field.Code = "host"
		field.ValidateCode = `
if (value == null) {
	throw new Error('should not be null');
}

if (value.length > 100) {
	throw new Error('too long');
}

return value;
`
		group.Add(field)
	}

	{
		field := NewTextField("端口", "port")
		field.MaxLength = 6
		field.Code = "port"
		group.Add(field)
	}

	form.Compose()

	v := url.Values{}
	v.Add("user_host", "192.168.1.100")
	v.Add("user_port", "1234")
	req, err := http.NewRequest(http.MethodPost, "http://teaos.cn", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Form = v
	values, _, err := form.ApplyRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(values)
}

func TestForm_Encode(t *testing.T) {
	form := NewForm("user")
	group := form.NewGroup()
	{
		field := NewTextField("主机", "Host")
		field.MaxLength = 100
		field.Code = "host"
		field.ValidateCode = `
if (value == null) {
	throw new Error('should not be null');
}

if (value.length > 100) {
	throw new Error('too long');
}

return value;
`
		field.Comment = "比如 127.0.0.1"
		group.Add(field)
	}

	{
		field := NewTextField("端口", "port")
		field.MaxLength = 6
		field.Code = "port"
		field.Placeholder = "8080"
		group.Add(field)
	}

	data, err := form.Encode()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

	{
		form := new(Form)
		err := FormDecode(data, form)
		if err != nil {
			t.Fatal(err)
		}
		logs.PrintAsJSON(form)
	}
}
