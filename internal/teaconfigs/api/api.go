package api

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"math/rand"
	"regexp"
	"time"
)

//  API定义
type API struct {
	shared.HeaderList

	Filename       string      `yaml:"filename" json:"filename"`             // 文件名
	Path           string      `yaml:"path" json:"path"`                     // 访问路径
	Address        string      `yaml:"address" json:"address"`               // 实际地址
	Methods        []string    `yaml:"methods" json:"methods"`               // 方法
	Params         []*APIParam `yaml:"params" json:"params"`                 // 参数
	Name           string      `yaml:"name" json:"name"`                     // 名称
	Description    string      `yaml:"description" json:"description"`       // 描述
	MockFiles      []string    `yaml:"mockFiles" json:"mockFiles"`           // 假数据文件（Mock）
	MockOn         bool        `yaml:"mockOn" json:"mockOn"`                 // 是否开启Mock
	Author         string      `yaml:"author" json:"author"`                 // 作者
	Company        string      `yaml:"company" json:"company"`               // 公司或团队
	IsAsynchronous bool        `yaml:"isAsynchronous" json:"isAsynchronous"` // TODO
	Timeout        float64     `yaml:"timeout" json:"timeout"`               // TODO
	MaxSize        uint        `yaml:"maxSize" json:"maxSize"`               // TODO
	TodoThings     []string    `yaml:"todo" json:"todo"`                     // 待做事宜
	DoneThings     []string    `yaml:"done" json:"done"`                     // 已完成事宜
	Response       []byte      `yaml:"response" json:"response"`             // 响应内容 TODO
	IsDeprecated   bool        `yaml:"isDeprecated" json:"isDeprecated"`     // 是否过期
	On             bool        `yaml:"on" json:"on"`                         // 是否开启
	Versions       []string    `yaml:"versions" json:"versions"`             // 版本信息
	ModifiedAt     int64       `yaml:"modifiedAt" json:"modifiedAt"`         // 最后修改时间
	Username       string      `yaml:"username" json:"username"`             // 最后修改用户名
	Groups         []string    `yaml:"groups" json:"groups"`                 // 分组
	Limit          *APILimit   `yaml:"limit" json:"limit"`                   // 限制 TODO
	AuthType       string      `yaml:"authType" json:"authType"`             // 认证方式

	TestScripts   []string `yaml:"testScripts" json:"testScripts"`     // 脚本文件
	TestCaseFiles []string `yaml:"testCaseFiles" json:"testCaseFiles"` // 单元测试存储文件

	CachePolicy string `yaml:"cachePolicy" json:"cachePolicy"` // 缓存策略
	CacheOn     bool   `yaml:"cacheOn" json:"cacheOn"`         // 缓存是否打开 TODO
	cachePolicy *shared.CachePolicy

	pathReg    *regexp.Regexp // 匹配模式
	pathParams []string
}

// 获取新API对象
func NewAPI() *API {
	return &API{
		On: true,
	}
}

// 从文件中读取API对象
func NewAPIFromFile(filename string) *API {
	if len(filename) == 0 {
		return nil
	}
	reader, err := files.NewReader(Tea.ConfigFile(filename))
	if err != nil {
		logs.Error(err)
		return nil
	}
	defer reader.Close()
	api := NewAPI()
	err = reader.ReadYAML(api)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return api
}

// 执行校验
func (this *API) Validate() error {
	// path
	this.pathParams = []string{}
	reg := regexp.MustCompile(`:\w+`)
	if reg.MatchString(this.Path) {
		newPath := reg.ReplaceAllStringFunc(this.Path, func(s string) string {
			param := s[1:]
			this.pathParams = append(this.pathParams, param)
			return "(.+)"
		})

		pathReg, err := regexp.Compile(newPath)
		if err != nil {
			return err
		}
		this.pathReg = pathReg
	}

	// limit
	if this.Limit != nil {
		err := this.Limit.Validate()
		if err != nil {
			return err
		}
	}

	// 校验缓存配置
	if len(this.CachePolicy) > 0 {
		policy := shared.NewCachePolicyFromFile(this.CachePolicy)
		if policy != nil {
			err := policy.Validate()
			if err != nil {
				return err
			}
			this.cachePolicy = policy
		}
	}

	// headers
	err := this.ValidateHeaders()
	if err != nil {
		return err
	}

	return nil
}

// 添加参数
func (this *API) AddParam(param *APIParam) {
	this.Params = append(this.Params, param)
}

// 使用正则匹配路径
func (this *API) Match(path string) (params map[string]string, matched bool) {
	if this.pathReg == nil {
		return nil, false
	}
	if !this.pathReg.MatchString(path) {
		return nil, false
	}

	params = map[string]string{}
	matched = true
	matches := this.pathReg.FindStringSubmatch(path)
	for index, match := range matches {
		if index == 0 {
			continue
		}
		params[this.pathParams[index-1]] = match
	}
	return
}

// 是否允许某个请求方法
func (this *API) AllowMethod(method string) bool {
	for _, m := range this.Methods {
		if m == method {
			return true
		}
	}
	return false
}

// 删除某个分组
func (this *API) RemoveGroup(name string) {
	result := []string{}
	for _, g := range this.Groups {
		if g != name {
			result = append(result, g)
		}
	}
	this.Groups = result
}

// 修改API某个分组名
func (this *API) ChangeGroup(oldName string, newName string) {
	result := []string{}
	for _, g := range this.Groups {
		if g == oldName {
			result = append(result, newName)
		} else {
			result = append(result, g)
		}
	}
	this.Groups = result
}

// 删除某个版本
func (this *API) RemoveVersion(name string) {
	result := []string{}
	for _, g := range this.Versions {
		if g != name {
			result = append(result, g)
		}
	}
	this.Versions = result
}

// 修改API某个版本号
func (this *API) ChangeVersion(oldName string, newName string) {
	result := []string{}
	for _, g := range this.Versions {
		if g == oldName {
			result = append(result, newName)
		} else {
			result = append(result, g)
		}
	}
	this.Versions = result
}

// 保存到文件
func (this *API) Save() error {
	if len(this.Filename) == 0 {
		this.Filename = "api." + rands.HexString(16) + ".conf"
	}
	writer, err := files.NewWriter(Tea.ConfigFile(this.Filename))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}

// 开始监控
func (this *API) StartWatching() {
	SharedApiWatching.Add(this.Path)
}

// 结束监控
func (this *API) StopWatching() {
	SharedApiWatching.Remove(this.Path)
}

// 是否在监控
func (this *API) IsWatching() bool {
	return SharedApiWatching.Contains(this.Path)
}

// 添加测试脚本
func (this *API) AddScript(script *APIScript) {
	if lists.ContainsString(this.TestScripts, script.Filename) {
		return
	}
	this.TestScripts = append(this.TestScripts, script.Filename)
}

// 读取所有测试脚本
func (this *API) FindTestScripts() []*APIScript {
	result := []*APIScript{}
	for _, filename := range this.TestScripts {
		reader, err := files.NewFile(Tea.ConfigFile(filename)).Reader()
		if err != nil {
			continue
		}
		script := NewAPIScript()
		err = reader.ReadYAML(script)
		reader.Close()
		if err != nil {
			continue
		}
		result = append([]*APIScript{script}, result ...)
	}
	return result
}

// 查找单个脚本
func (this *API) FindTestScript(filename string) *APIScript {
	if !lists.ContainsString(this.TestScripts, filename) {
		return nil
	}

	reader, err := files.NewFile(Tea.ConfigFile(filename)).Reader()
	if err != nil {
		return nil
	}

	script := NewAPIScript()
	err = reader.ReadYAML(script)
	reader.Close()

	if err != nil {
		return nil
	}

	return script
}

// 删除测试脚本
func (this *API) DeleteTestScript(filename string) error {
	if lists.ContainsString(this.TestScripts, filename) {
		script := NewAPIScript()
		script.Filename = filename
		err := script.Delete()
		if err != nil {
			return err
		}
	}

	this.TestScripts = lists.Delete(this.TestScripts, filename).([]string)
	return nil
}

// 删除API
func (this *API) Delete() error {
	if len(this.Filename) == 0 {
		return nil
	}

	// 删除脚本
	for _, scriptFile := range this.TestScripts {
		files.NewFile(Tea.ConfigFile(scriptFile)).DeleteIfExists()
	}

	return files.NewFile(Tea.ConfigFile(this.Filename)).DeleteIfExists()
}

// 添加测试用例
func (this *API) AddTestCase(filename string) {
	if len(filename) == 0 {
		return
	}
	if lists.ContainsString(this.TestCaseFiles, filename) {
		return
	}
	this.TestCaseFiles = append(this.TestCaseFiles, filename)
}

// 查找所有的测试用例
func (this *API) FindTestCases() []*APITestCase {
	cases := []*APITestCase{}
	for _, filename := range this.TestCaseFiles {
		case1 := NewAPITestCaseFromFile(filename)
		if case1 != nil {
			cases = append(cases, case1)
		}
	}
	return cases
}

// 删除测试用例
func (this *API) DeleteTestCase(filename string) {
	this.TestCaseFiles = lists.Delete(this.TestCaseFiles, filename).([]string)
}

// 添加Mock文件
func (this *API) AddMock(filename string) {
	if len(filename) == 0 {
		return
	}
	if lists.ContainsString(this.MockFiles, filename) {
		return
	}
	this.MockFiles = append(this.MockFiles, filename)
}

// 获取所有Mock的文件
func (this *API) MockDataFiles() []string {
	result := []string{}
	for _, filename := range this.MockFiles {
		mock := NewAPIMockFromFile(filename)
		if mock != nil && len(mock.File) > 0 {
			result = append(result, mock.File)
		}
	}
	return result
}

// 删除Mock
func (this *API) DeleteMock(mockFile string) {
	this.MockFiles = lists.Delete(this.MockFiles, mockFile).([]string)
}

// 随机取得一个Mock
func (this *API) RandMock() *APIMock {
	if len(this.MockFiles) == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	file := this.MockFiles[rand.Int()%len(this.MockFiles)]
	if len(file) == 0 {
		return nil
	}

	return NewAPIMockFromFile(file)
}

// 缓存策略
func (this *API) CachePolicyObject() *shared.CachePolicy {
	return this.cachePolicy
}
