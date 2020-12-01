package agents

import (
	"bytes"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"os/exec"
	"runtime"
	"strings"
)

// Script文件数据源
type ScriptAction struct {
	Path       string             `yaml:"path" json:"path"`
	ScriptType string             `yaml:"scriptType" json:"scriptType"` // 脚本类型，可以为path, code
	ScriptLang string             `yaml:"scriptLang" json:"scriptLang"` // 脚本语言
	Script     string             `yaml:"script" json:"script"`         // 脚本代码
	Env        []*shared.Variable `yaml:"env" json:"env"`               // 环境变量设置
	Cwd        string             `yaml:"cwd" json:"cwd"`
}

// 获取新对象
func NewScriptAction() *ScriptAction {
	return &ScriptAction{
		Env: []*shared.Variable{},
	}
}

func (this *ScriptAction) Validate() error {
	return nil
}

func (this *ScriptAction) Name() string {
	return "Shell脚本"
}

// 代号
func (this *ScriptAction) Code() string {
	return "script"
}

// 描述
func (this *ScriptAction) Description() string {
	return "通过执行本地的Shell脚本文件"
}

// 格式化脚本
func (this *ScriptAction) FormattedScript() string {
	script := this.Script
	script = strings.Replace(script, "\r", "", -1)
	return script
}

// 保存到本地
func (this *ScriptAction) Generate(id string) (path string, err error) {
	if runtime.GOOS == "windows" {
		path = Tea.ConfigFile("agents/action." + id + ".bat")
	} else {
		path = Tea.ConfigFile("agents/action." + id + ".script")
	}
	shFile := files.NewFile(path)
	if !shFile.Exists() {
		err = shFile.WriteString(this.FormattedScript())
		if err != nil {
			return
		}
		err = shFile.Chmod(0777)
		if err != nil {
			return
		}
	}
	return
}

// 执行
func (this *ScriptAction) Run(params map[string]string) (result string, err error) {
	currentPath := this.Path

	// 脚本
	if this.ScriptType == "code" {
		path, err := this.Generate(rands.HexString(16))
		if err != nil {
			return "", err
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
		return "", errors.New("path or script should not be empty")
	}

	cmd := exec.Command(currentPath)

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
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		// do nothing
	}

	if stderr.Len() > 0 {
		result += string(stderr.Bytes())
	}

	result += string(stdout.Bytes())
	return result, nil
}

// 获取简要信息
func (this *ScriptAction) Summary() maps.Map {
	return maps.Map{
		"name":        this.Name(),
		"code":        this.Code(),
		"description": this.Description(),
	}
}

// 添加环境变量
func (this *ScriptAction) AddEnv(name, value string) {
	this.Env = append(this.Env, &shared.Variable{
		Name:  name,
		Value: value,
	})
}
