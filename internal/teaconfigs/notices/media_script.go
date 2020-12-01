package notices

import (
	"bytes"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"os/exec"
	"runtime"
	"strings"
)

// 脚本媒介
type NoticeScriptMedia struct {
	Path       string             `yaml:"path" json:"path"`
	ScriptType string             `yaml:"scriptType" json:"scriptType"` // 脚本类型，可以为path, code
	ScriptLang string             `yaml:"scriptLang" json:"scriptLang"` // 脚本语言
	Script     string             `yaml:"script" json:"script"`         // 脚本代码
	Cwd        string             `yaml:"cwd" json:"cwd"`
	Env        []*shared.Variable `yaml:"env" json:"env"`
}

// 获取新对象
func NewNoticeScriptMedia() *NoticeScriptMedia {
	return &NoticeScriptMedia{}
}

// 添加环境变量
func (this *NoticeScriptMedia) AddEnv(name, value string) {
	this.Env = append(this.Env, &shared.Variable{
		Name:  name,
		Value: value,
	})
}

// 格式化脚本
func (this *NoticeScriptMedia) FormattedScript() string {
	script := this.Script
	script = strings.Replace(script, "\r", "", -1)
	return script
}

// 保存到本地
func (this *NoticeScriptMedia) Generate(id string) (path string, err error) {
	if runtime.GOOS == "windows" {
		path = teautils.TmpFile("notice." + id + ".bat")
	} else {
		path = teautils.TmpFile("notice." + id + ".script")
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

// 发送
func (this *NoticeScriptMedia) Send(user string, subject string, body string) (resp []byte, err error) {
	// 脚本
	if this.ScriptType == "code" {
		path, err := this.Generate(rands.HexString(16))
		if err != nil {
			return nil, err
		}
		this.Path = path

		defer func() {
			f := files.NewFile(this.Path)
			if f.Exists() {
				err := f.Delete()
				if err != nil {
					logs.Error(err)
				}
			}
		}()
	}

	if len(this.Path) == 0 {
		return nil, errors.New("'path' should be specified")
	}

	cmd := exec.Command(this.Path)

	if len(this.Env) > 0 {
		for _, env := range this.Env {
			cmd.Env = append(cmd.Env, env.Name+"="+env.Value)
		}
	}
	cmd.Env = append(cmd.Env, "NoticeUser="+user)
	cmd.Env = append(cmd.Env, "NoticeSubject="+subject)
	cmd.Env = append(cmd.Env, "NoticeBody="+body)

	if len(this.Cwd) > 0 {
		cmd.Dir = this.Cwd
	}

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		// do nothing
	}

	if stderr.Len() > 0 {
		return stdout.Bytes(), errors.New(string(stderr.Bytes()))
	}

	return stdout.Bytes(), nil
}

// 是否需要用户标识
func (this *NoticeScriptMedia) RequireUser() bool {
	return false
}
