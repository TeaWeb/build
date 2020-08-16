package agents

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/iwind/TeaGo/files"
)

// 数据文件
type FileSource struct {
	Source `yaml:",inline"`

	Path string `yaml:"path" json:"path"`
}

// 获取新对象
func NewFileSource() *FileSource {
	return &FileSource{}
}

// 校验
func (this *FileSource) Validate() error {
	if len(this.Path) == 0 {
		return errors.New("path should not be empty")
	}

	return nil
}

// 名称
func (this *FileSource) Name() string {
	return "数据文件"
}

// 代号
func (this *FileSource) Code() string {
	return "file"
}

// 描述
func (this *FileSource) Description() string {
	return "通过读取本地文件获取数据"
}

// 执行
func (this *FileSource) Execute(params map[string]string) (value interface{}, err error) {
	if len(this.Path) == 0 {
		return nil, errors.New("path should not be empty")
	}

	file := files.NewFile(this.Path)
	if !file.Exists() {
		return nil, errors.New("file does not exist")
	}

	data, err := file.ReadAll()
	if err != nil {
		return nil, err
	}
	return DecodeSource(data, this.DataFormat)
}

// 选项表单
func (this *FileSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	group := form.NewGroup()

	{
		field := forms.NewTextField("数据文件路径", "")
		field.IsRequired = true
		field.Code = "path"
		field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入数据文件路径");
}

return value;
`

		group.Add(field)
	}

	return form
}

// 显示信息
func (this *FileSource) Presentation() *forms.Presentation {
	return &forms.Presentation{
		HTML: `
<tr>
	<td>数据文件路径</td>
	<td>{{source.path}}</td>
</tr>`,
	}
}
