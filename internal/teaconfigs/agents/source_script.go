package agents

import (
	"bytes"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/utils/string"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Script文件数据源
type ScriptSource struct {
	Source `yaml:",inline"`

	Path       string             `yaml:"path" json:"path"`
	ScriptType string             `yaml:"scriptType" json:"scriptType"` // 脚本类型，可以为path, code
	ScriptLang string             `yaml:"scriptLang" json:"scriptLang"` // 脚本语言
	Script     string             `yaml:"script" json:"script"`         // 脚本代码
	Env        []*shared.Variable `yaml:"env" json:"env"`               // 环境变量设置
	Cwd        string             `yaml:"cwd" json:"cwd"`
}

// 获取新对象
func NewScriptSource() *ScriptSource {
	return &ScriptSource{
		Env: []*shared.Variable{},
	}
}

// 校验
func (this *ScriptSource) Validate() error {
	return nil
}

// 名称
func (this *ScriptSource) Name() string {
	return "Shell脚本"
}

// 代号
func (this *ScriptSource) Code() string {
	return "script"
}

// 描述
func (this *ScriptSource) Description() string {
	return "通过执行本地的Shell脚本文件获取数据"
}

// 格式化脚本
func (this *ScriptSource) FormattedScript() string {
	script := this.Script
	script = strings.Replace(script, "\r", "", -1)

	// 是否有解释器
	if lists.ContainsString([]string{"darwin", "linux", "freebsd"}, runtime.GOOS) {
		if !strings.HasPrefix(strings.TrimSpace(script), "#!") {
			script = "#!/usr/bin/env bash\n" + script
		} else {
			script = strings.TrimLeft(script, " \r\n\t")
		}
	}

	return script
}

// 保存到本地
func (this *ScriptSource) Generate(id string) (path string, err error) {
	if runtime.GOOS == "windows" {
		path = Tea.ConfigFile("agents/source." + id + ".bat")
	} else {
		path = Tea.ConfigFile("agents/source." + id + ".script")
	}
	shFile := files.NewFile(path)
	if !shFile.Exists() {
		err = shFile.WriteString(this.FormattedScript())
		if err != nil {
			return
		}
		if runtime.GOOS != "windows" {
			err = shFile.Chmod(0777)
			if err != nil {
				return
			}
		}
	}
	return
}

// 执行
func (this *ScriptSource) Execute(params map[string]string) (value interface{}, err error) {
	currentPath := this.Path

	// 脚本
	if this.ScriptType == "code" {
		path, err := this.Generate(rands.HexString(16) + stringutil.Md5(this.Code()))
		if err != nil {
			return nil, err
		}
		currentPath = path

		defer func() {
			f := files.NewFile(currentPath)
			if f.Exists() {
				err := f.Delete()
				if err != nil {
					logs.Error(err)
				}
			}
		}()
	}

	if len(currentPath) == 0 {
		return nil, errors.New("path or script should not be empty")
	}

	// 检查shell
	var cmd *exec.Cmd = nil
	if lists.ContainsString([]string{"darwin", "linux", "freebsd"}, runtime.GOOS) {
		data, err := files.NewFile(currentPath).ReadAll()
		if err == nil {
			if !strings.HasPrefix(strings.TrimSpace(string(data)), "#!") {
				cmd = exec.Command("sh", currentPath)
			}
		}
		if cmd == nil {
			cmd = exec.Command(currentPath)
		}
	} else {
		cmd = exec.Command(currentPath)
	}

	if len(this.Env) > 0 {
		for _, env := range this.Env {
			cmd.Env = append(cmd.Env, env.Name+"="+env.Value)
		}
	}

	if len(params) > 0 {
		for key, value := range params {
			cmd.Env = append(cmd.Env, key+"="+value)
		}
	}

	if len(this.Cwd) > 0 {
		cmd.Dir = this.Cwd
	}

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Start()
	if err != nil {
		if strings.Index(err.Error(), "text file busy") > -1 {
			// try again
			time.Sleep(100 * time.Millisecond)
			err = cmd.Start()
			if err != nil {
				return nil, err
			}
		}
	}

	err = cmd.Wait()
	if err != nil {
		// do nothing
	}

	if stderr.Len() > 0 {
		logs.Println("error:", string(stderr.Bytes()))
	}

	return DecodeSource(stdout.Bytes(), this.DataFormat)
}

// 添加环境变量
func (this *ScriptSource) AddEnv(name, value string) {
	this.Env = append(this.Env, &shared.Variable{
		Name:  name,
		Value: value,
	})
}

// 选项表单
func (this *ScriptSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	{
		group := form.NewGroup()
		{
			field := forms.NewScriptBox("脚本", "")
			field.Code = this.Code()
			field.IsRequired = true
			field.InitCode = `return {
	"scriptType": values.scriptType,
	"scriptCode": values.script,
	"scriptLang": values.scriptLang,
	"scriptPath": values.path
}`
			group.Add(field)
		}
	}

	{
		group := form.NewGroup()
		{
			field := forms.NewTextField("当前工作目录", "CWD")
			field.Code = "cwd"
			field.MaxLength = 500
			group.Add(field)
		}

		{
			field := forms.NewEnvBox("环境变量", "ENV")
			field.Code = "env"
			group.Add(field)
		}
	}

	form.ValidateCode = `
if (values.script.scriptType == "path") {
	if (values.script.scriptPath.length == 0) {
		return FieldError("scriptPath", "请输入脚本路径")
	}
} else {
	if (values.script.scriptCode.length == 0) {
		return FieldError("scriptCode", "请输入脚本代码");
	}
}

return {
	"cwd": values.cwd,
	"env": values.env,
	"scriptType": values.script.scriptType,
	"script": values.script.scriptCode,
	"scriptLang": values.script.scriptLang,
	"path": values.script.scriptPath
}
`

	return form
}

// 显示界面
func (this *ScriptSource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
 <tr v-if="source.scriptType == 'path' || source.scriptType == null || source.scriptType.length == 0">
	<td>脚本路径</td>
	<td>
		{{source.path}}
	</td>
</tr>
<tr v-if="source.scriptType == 'code'">
	<td>脚本代码</td>
	<td>
		<div id="script-code-editor"></div>
	</td>
</tr>
<tr>
	<td>当前工作目录<em>（CWD）</em></td>
	<td>
		<span v-if="source.cwd.length > 0">{{source.cwd}}</span>
		<span v-if="source.cwd.length == 0" class="disabled">没有设置</span>
	</td>
</tr>
<tr>
	<td>环境变量<em>（ENV）</em></td>
	<td>
		<span v-if="source.env != null && source.env.length > 0" class="ui label small" v-for="(var1, index) in source.env">
			<em>{{var1.name}}</em>: {{var1.value}}
			<a href="" @click.prevent="removeEnv(index)"></a>
		</span>
		<span v-if="source.env == null || source.env.length == 0" class="disabled">没有设置</span>
	</td>
</tr>
`
	p.Javascript = `
if (this.source.scriptType == "code") {
	this.loadCodeEditor(this.source.scriptLang, this.source.script);
}
`
	return p
}
