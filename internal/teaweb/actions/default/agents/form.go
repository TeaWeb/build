package agents

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type FormAction actions.Action

// 表单测试
func (this *FormAction) Run(params struct{}) {
	form := forms.NewForm("user")

	group := form.NewGroup()

	{
		field := &forms.TextField{}
		field.Title = "主机"
		field.Subtitle = "Host"
		field.Placeholder = "请输入主机名"
		field.Comment = "比如127.0.0.1"
		field.Code = "host"
		field.Attr("id", "host-field")
		//field.Attr("v-model", "host")
		field.Attr("@input", "changeHost()")
		field.CSS = `
#host-field {
	background: red!important;
}
`
		field.Javascript = `
this.host = "";

this.$delay(function () {
	this.$find("#host-field").focus();
});

this.changeHost = function () {
	console.log("change:", this.host);
};
`
		group.Add(field)
	}

	{
		field := &forms.TextField{}
		field.Title = "端口"
		field.Subtitle = "Port"
		field.Placeholder = "请输入端口"
		field.Comment = "比如8080"
		field.Code = "port"
		field.Value = "8080"
		group.Add(field)
	}

	{
		field := forms.NewOptions("请求方法", "")
		field.Code = "method"
		field.AddOption("GET", "GET")
		field.AddOption("POST", "POST")
		field.Attr("style", "width:10em")
		group.Add(field)
	}

	{
		field := forms.NewTextField("刷新间隔", "")
		field.Code = "interval"
		field.RightLabel = "分钟"
		field.Attr("style", "width:5em")
		field.Value = "30"
		group.Add(field)
	}

	{
		field := forms.NewScriptBox("脚本", "")
		field.Code = "script"
		field.IsRequired = true
		group.Add(field)
	}

	{
		field := forms.NewCheckBox("是否启用", "")
		field.IsChecked = true
		field.Label = "这是标签"
		field.Code = "on"
		field.Value = 1
		group.Add(field)
	}

	{
		field := forms.NewTextBox("描述", "")
		field.Code = "description"
		field.Value = "描述文字..."
		field.Placeholder = "请输入..."
		field.Rows = 5
		field.Cols = 60
		group.Add(field)
	}

	{
		field := forms.NewEnvBox("环境变量", "")
		field.Code = "env"
		group.Add(field)
	}

	data := []byte(`{
   "description": "描述文字...",
   "env": [
      {
         "name": "a",
         "value": "b"
      },
      {
         "name": "name",
         "value": "YinLu"
      }
   ],
   "host": "127.0.0.1",
   "interval": "55",
   "method": "POST",
   "on": "0",
   "port": "8881",
   "script": {
      "scriptCode": "hello world",
      "scriptLang": "php",
      "scriptPath": "/home/www.sh",
      "scriptType": "code"
   }
}`)
	m := map[string]interface{}{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		this.Fail(err.Error())
	}
	form.Init(m)

	form.Compose()
	this.Data["form"] = form

	this.Show()
}

func (this *FormAction) RunPost(params struct{}) {
	form := forms.NewForm("user")
	group := form.NewGroup()

	{
		field := &forms.TextField{}
		field.Title = "主机"
		field.Subtitle = "Host"
		field.Placeholder = "请输入主机名"
		field.Comment = "比如127.0.0.1"
		field.Code = "host"
		group.Add(field)
	}

	{
		field := &forms.TextField{}
		field.Title = "端口"
		field.Subtitle = "Port"
		field.Placeholder = "请输入端口"
		field.Comment = "比如8080"
		field.Code = "port"
		field.ValidateCode = `if (value.length < 4) {
	throw new Error("端口长度不能小于4");
}
return value;
`
		group.Add(field)
	}

	{
		field := forms.NewOptions("请求方法", "")
		field.Code = "method"
		field.AddOption("GET", "GET")
		field.AddOption("POST", "POST")
		field.Attr("style", "width:10em")
		group.Add(field)
	}

	{
		field := forms.NewTextField("刷新间隔", "")
		field.Code = "interval"
		field.RightLabel = "分钟"
		field.Attr("style", "width:5em")
		field.Value = "30"
		group.Add(field)
	}

	{
		field := forms.NewScriptBox("脚本", "")
		field.IsRequired = true
		field.Code = "script"
		group.Add(field)
	}

	{
		field := forms.NewCheckBox("是否启用", "")
		field.IsChecked = true
		field.Code = "on"
		field.Label = "这是标签"
		group.Add(field)
	}

	{
		field := forms.NewTextBox("描述", "")
		field.Code = "description"
		field.Value = "描述文字..."
		field.Placeholder = "请输入..."
		field.Rows = 5
		field.Cols = 60
		group.Add(field)
	}

	{
		field := forms.NewEnvBox("环境变量", "")
		field.Code = "env"
		group.Add(field)
	}

	values, field, err := form.ApplyRequest(this.Request)
	if err != nil {
		this.FailField(form.Namespace+"_"+field, err.Error())
	}

	logs.PrintAsJSON(values)
}
