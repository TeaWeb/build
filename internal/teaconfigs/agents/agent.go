package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"gopkg.in/yaml.v3"
)

// Agent定义
type AgentConfig struct {
	Id                  string       `yaml:"id" json:"id"`                                   // ID
	On                  bool         `yaml:"on" json:"on"`                                   // 是否启用
	Name                string       `yaml:"name" json:"name"`                               // 名称
	Host                string       `yaml:"host" json:"host"`                               // 主机地址
	Key                 string       `yaml:"key" json:"key"`                                 // 密钥
	AllowAll            bool         `yaml:"allowAll" json:"allowAll"`                       // 是否允许所有的IP
	Allow               []string     `yaml:"allow" json:"allow"`                             // 允许的IP地址
	Apps                []*AppConfig `yaml:"apps" json:"apps"`                               // Apps
	TeaVersion          string       `yaml:"teaVersion" json:"teaVersion"`                   // TeaWeb版本
	Version             uint         `yaml:"version" json:"version"`                         // 版本
	CheckDisconnections bool         `yaml:"checkDisconnections" json:"checkDisconnections"` // 是否检查离线
	CountDisconnections int          `yaml:"countDisconnections" json:"countDisconnections"` // 错误次数
	GroupIds            []string     `yaml:"groupIds" json:"groupIds"`                       // 分组IDs
	AutoUpdates         bool         `yaml:"autoUpdates" json:"autoUpdates"`                 // 是否开启自动更新
	AppsIsInitialized   bool         `yaml:"appsIsInitialized" json:"appsIsInitialized"`     // 是否已经初始化App
	GroupKey            string       `yaml:"groupKey" json:"groupKey"`                       // 注册使用的密钥

	NoticeSetting map[notices.NoticeLevel][]*notices.NoticeReceiver `yaml:"noticeSetting" json:"noticeSetting"`
}

// 获取新对象
func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		On:                  true,
		Id:                  rands.HexString(16),
		CheckDisconnections: true,
	}
}

// 本地Agent
var localAgentConfig *AgentConfig = nil

func LocalAgentConfig() *AgentConfig {
	if localAgentConfig == nil {
		localAgentConfig = &AgentConfig{
			On:       true,
			Id:       "local",
			Name:     "本地",
			Key:      rands.HexString(16),
			AllowAll: false,
			Allow:    []string{"127.0.0.1"},
			Host:     "127.0.0.1",
		}
	}
	return localAgentConfig
}

// 从文件中获取对象
func NewAgentConfigFromFile(filename string) *AgentConfig {
	reader, err := files.NewReader(Tea.ConfigFile("agents/" + filename))
	if err != nil {
		return nil
	}
	defer func() {
		err = reader.Close()
		if err != nil {
			logs.Error(err)
		}
	}()
	agent := &AgentConfig{}
	err = reader.ReadYAML(agent)
	if err != nil {
		return nil
	}
	return agent
}

// 根据ID获取对象
func NewAgentConfigFromId(agentId string) *AgentConfig {
	if len(agentId) == 0 {
		return nil
	}
	agent := NewAgentConfigFromFile("agent." + agentId + ".conf")
	if agent != nil {
		if agent.Id == "local" && len(agent.Name) == 0 {
			agent.Name = "本地"
		}

		return agent
	}

	if agentId == "local" {
		return LocalAgentConfig()
	}

	return nil
}

// 判断是否为Local Agent
func (this *AgentConfig) IsLocal() bool {
	return this.Id == "local"
}

// 校验
func (this *AgentConfig) Validate() error {
	for _, a := range this.Apps {
		err := a.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// 文件名
func (this *AgentConfig) Filename() string {
	return "agent." + this.Id + ".conf"
}

// 保存
func (this *AgentConfig) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	defer func() {
		NotifyAgentsChange() // 标记列表改变
	}()

	dirFile := files.NewFile(Tea.ConfigFile("agents"))
	if !dirFile.Exists() {
		err := dirFile.Mkdir()
		if err != nil {
			logs.Error(err)
		}
	}

	writer, err := files.NewWriter(Tea.ConfigFile("agents/" + this.Filename()))
	if err != nil {
		return err
	}
	defer func() {
		err := writer.Close()
		if err != nil {
			logs.Error(err)
		}
	}()
	this.Version++
	this.TeaVersion = teaconst.TeaVersion
	_, err = writer.WriteYAML(this)
	return err
}

// 删除
func (this *AgentConfig) Delete() error {
	defer func() {
		NotifyAgentsChange() // 标记列表改变
	}()

	// 删除board
	{
		f := files.NewFile(Tea.ConfigFile("agents/board." + this.Id + ".conf"))
		if f.Exists() {
			err := f.Delete()
			if err != nil {
				return err
			}
		}
	}

	f := files.NewFile(Tea.ConfigFile("agents/" + this.Filename()))
	return f.Delete()
}

// 添加App
func (this *AgentConfig) AddApp(app *AppConfig) {
	this.Apps = append(this.Apps, app)
}

// 替换App，如果不存在则增加
func (this *AgentConfig) ReplaceApp(app *AppConfig) {
	found := false
	for index, a := range this.Apps {
		if a.Id == app.Id {
			this.Apps[index] = app
			found = true
			break
		}
	}
	if !found {
		this.Apps = append(this.Apps, app)
	}
}

// 添加一组App
func (this *AgentConfig) AddApps(apps []*AppConfig) {
	this.Apps = append(this.Apps, apps...)
}

// 删除App
func (this *AgentConfig) RemoveApp(appId string) {
	result := []*AppConfig{}
	for _, a := range this.Apps {
		if a.Id == appId {
			continue
		}
		result = append(result, a)
	}
	this.Apps = result
}

// 移动App位置
func (this *AgentConfig) MoveApp(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(this.Apps) {
		return
	}
	if toIndex < 0 || toIndex >= len(this.Apps) {
		return
	}
	if fromIndex == toIndex {
		return
	}

	location := this.Apps[fromIndex]
	newList := []*AppConfig{}
	for i := 0; i < len(this.Apps); i++ {
		if i == fromIndex {
			continue
		}
		if fromIndex > toIndex && i == toIndex {
			newList = append(newList, location)
		}
		newList = append(newList, this.Apps[i])
		if fromIndex < toIndex && i == toIndex {
			newList = append(newList, location)
		}
	}

	this.Apps = newList
}

// 查找App
func (this *AgentConfig) FindApp(appId string) *AppConfig {
	for _, a := range this.Apps {
		if a.Id == appId {
			return a
		}
	}
	return nil
}

// 判断是否有某个App
func (this *AgentConfig) HasApp(appId string) bool {
	for _, a := range this.Apps {
		if a.Id == appId {
			return true
		}
	}
	return false
}

// YAML编码
func (this *AgentConfig) EncodeYAML() ([]byte, error) {
	return yaml.Marshal(this)
}

// 查找任务
func (this *AgentConfig) FindTask(taskId string) (appConfig *AppConfig, taskConfig *TaskConfig) {
	for _, app := range this.Apps {
		for _, task := range app.Tasks {
			if task.Id == taskId {
				return app, task
			}
		}
	}
	return nil, nil
}

// 查找监控项
func (this *AgentConfig) FindItem(itemId string) (appConfig *AppConfig, item *Item) {
	for _, app := range this.Apps {
		for _, item := range app.Items {
			if item.Id == itemId {
				return app, item
			}
		}
	}
	return nil, nil
}

// 添加分组
func (this *AgentConfig) AddGroup(groupId string) {
	if lists.ContainsString(this.GroupIds, groupId) {
		return
	}
	this.GroupIds = append(this.GroupIds, groupId)
}

// 删除分组
func (this *AgentConfig) RemoveGroup(groupId string) {
	result := []string{}
	for _, g := range this.GroupIds {
		if g == groupId {
			continue
		}
		result = append(result, g)
	}
	this.GroupIds = result
}

// 判断是否有某些分组
func (this *AgentConfig) BelongsToGroups(groupIds []string) bool {
	if len(this.GroupIds) == 0 {
		this.GroupIds = []string{"default"}
	}
	if len(groupIds) == 0 {
		groupIds = []string{"default"}
	}
	for _, groupId := range groupIds {
		b := lists.ContainsString(this.GroupIds, groupId)
		if b {
			return true
		}
	}
	return false
}

// 判断是否有某个分组
func (this *AgentConfig) BelongsToGroup(groupId string) bool {
	if len(this.GroupIds) == 0 {
		this.GroupIds = []string{"default"}
	}
	if len(groupId) == 0 {
		groupId = "default"
	}
	return lists.ContainsString(this.GroupIds, groupId)
}

// 添加内置的App
func (this *AgentConfig) AddDefaultApps() {
	this.AppsIsInitialized = true
	{
		app := NewAppConfig()
		app.Id = "system"
		app.Name = "系统"
		this.AddApp(app)

		board := NewAgentBoard(this.Id)

		// 添加到看板
		defer func() {
			err := board.Save()
			if err != nil {
				logs.Error(err)
			}
		}()

		// cpu
		{
			// item
			item := NewItem()
			item.Id = "cpu.usage"
			item.Name = "CPU使用量（%）"
			item.Interval = "60s"

			source := NewCPUSource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)
			app.AddItem(item)

			// 阈值
			item.AddThreshold(source.Thresholds()...)

			// chart
			item.AddFilterCharts(source.Charts(), "cpu.chart1")
			board.AddChart(app.Id, item.Id, "cpu.chart1")
		}

		// load
		{
			// item
			item := NewItem()
			item.Id = "cpu.load"
			item.Name = "负载（Load）"
			item.Interval = "60s"

			source := NewLoadSource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)

			app.AddItem(item)

			// 阈值
			item.AddThreshold(source.Thresholds()...)

			// chart
			item.AddFilterCharts(source.Charts(), "cpu.load.chart1")
			board.AddChart(app.Id, item.Id, "cpu.load.chart1")
		}

		// memory usage
		{
			//item
			item := NewItem()
			item.Id = "memory.usage"
			item.Name = "内存使用量"
			item.Interval = "60s"

			source := NewMemorySource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)

			app.AddItem(item)

			// 阈值
			item.AddThreshold(source.Thresholds()...)

			// chart
			item.AddFilterCharts(source.Charts(), "memory.usage.chart1", "memory.usage.chart2")
			board.AddChart(app.Id, item.Id, "memory.usage.chart1")
			board.AddChart(app.Id, item.Id, "memory.usage.chart2")
		}

		// clock
		{
			// item
			item := NewItem()
			item.Id = "clock"
			item.Name = "时钟"
			item.Interval = "60s"

			source := NewDateSource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)

			app.AddItem(item)

			// 时钟
			{
				chart := widgets.NewChart()
				chart.Id = "clock"
				chart.Name = "时钟"
				chart.Columns = 1
				chart.Type = "javascript"
				chart.Options = maps.Map{
					"code": `
var chart = new charts.Clock();
var latest = NewQuery().latest(1);
if (latest.length > 0) {
	chart.timestamp = parseInt(new Date().getTime() / 1000) - (latest[0].createdAt - latest[0].value.timestamp);
}
chart.render();
`,
				}
				item.AddChart(chart)
				board.AddChart(app.Id, item.Id, chart.Id)
			}
		}

		// network out && network in
		{
			// item
			item := NewItem()
			item.Id = "network.usage"
			item.Name = "网络相关"
			item.Interval = "60s"

			source := NewNetworkSource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)

			app.AddItem(item)

			// 阈值
			item.AddThreshold(source.Thresholds()...)

			// 图表
			item.AddFilterCharts(source.Charts(), "network.usage.sent", "network.usage.received")
			board.AddChart(app.Id, item.Id, "network.usage.sent")
			board.AddChart(app.Id, item.Id, "network.usage.received")
		}

		// disk
		{
			// item
			item := NewItem()
			item.Id = "disk.usage"
			item.Name = "文件系统"
			item.Interval = "120s"

			source := NewDiskSource()
			source.DataFormat = SourceDataFormatJSON
			item.SourceCode = source.Code()
			item.SourceOptions = ConvertSourceToMap(source)

			app.AddItem(item)

			// 阈值
			item.AddThreshold(source.Thresholds()...)

			// 图表
			item.AddFilterCharts(source.Charts(), "disk.usage.chart1")
			board.AddChart(app.Id, item.Id, "disk.usage.chart1")
		}
	}
}

// 添加通知接收者
func (this *AgentConfig) AddNoticeReceiver(level notices.NoticeLevel, receiver *notices.NoticeReceiver) {
	if this.NoticeSetting == nil {
		this.NoticeSetting = map[notices.NoticeLevel][]*notices.NoticeReceiver{}
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		receivers = []*notices.NoticeReceiver{}
	}
	receivers = append(receivers, receiver)
	this.NoticeSetting[level] = receivers
}

// 删除通知接收者
func (this *AgentConfig) RemoveNoticeReceiver(level notices.NoticeLevel, receiverId string) {
	if this.NoticeSetting == nil {
		return
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		return
	}

	result := []*notices.NoticeReceiver{}
	for _, r := range receivers {
		if r.Id == receiverId {
			continue
		}
		result = append(result, r)
	}
	this.NoticeSetting[level] = result
}

// 获取通知接收者数量
func (this *AgentConfig) CountNoticeReceivers() int {
	count := 0
	for _, receivers := range this.NoticeSetting {
		count += len(receivers)
	}
	return count
}

// 删除媒介
func (this *AgentConfig) RemoveMedia(mediaId string) (found bool) {
	for level, receivers := range this.NoticeSetting {
		result := []*notices.NoticeReceiver{}
		for _, receiver := range receivers {
			if receiver.MediaId == mediaId {
				found = true
				continue
			}
			result = append(result, receiver)
		}
		this.NoticeSetting[level] = result
	}
	return
}

// 查找一个或多个级别对应的接收者，并合并相同的接收者
func (this *AgentConfig) FindAllNoticeReceivers(level ...notices.NoticeLevel) []*notices.NoticeReceiver {
	if len(level) == 0 {
		return []*notices.NoticeReceiver{}
	}

	m := maps.Map{} // mediaId_user => bool
	result := []*notices.NoticeReceiver{}
	for _, l := range level {
		receivers, ok := this.NoticeSetting[l]
		if !ok {
			continue
		}
		for _, receiver := range receivers {
			if !receiver.On {
				continue
			}
			key := receiver.Key()
			if m.Has(key) {
				continue
			}
			m[key] = true
			result = append(result, receiver)
		}
	}
	return result
}

// 获取分组名
func (this *AgentConfig) GroupName() string {
	if len(this.GroupIds) == 0 {
		return "默认分组"
	}
	groupId := this.GroupIds[0]
	if len(groupId) == 0 {
		return "默认分组"
	}

	group := SharedGroupList().FindGroup(groupId)
	if group == nil {
		return "默认分组"
	}
	return group.Name
}

// 取得第一个分组
func (this *AgentConfig) FirstGroup() *Group {
	if len(this.GroupIds) == 0 {
		return SharedGroupList().FindDefaultGroup()
	}
	group := SharedGroupList().FindGroup(this.GroupIds[0])
	if group != nil {
		return group
	}
	return SharedGroupList().FindDefaultGroup()
}

// 判断是否匹配关键词
func (this *AgentConfig) MatchKeyword(keyword string) (matched bool, name string, tags []string) {
	if teautils.MatchKeyword(this.Name, keyword) ||
		teautils.MatchKeyword(this.Host, keyword) ||
		this.Id == keyword {
		matched = true
		name = this.Name
		if len(this.Host) > 0 {
			tags = []string{this.Host}
		}
		return
	}
	return
}
