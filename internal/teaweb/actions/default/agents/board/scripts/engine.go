package scripts

import (
	"encoding/json"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teamemory"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/robertkrimen/otto"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var engineCache = teamemory.NewGrid(1, teamemory.NewLimitCountOpt(1000))
var dayReg = regexp.MustCompile(`^(\d+)-(\d+)-(\d+)$`)

// 脚本引擎
type Engine struct {
	vm           *otto.Otto
	chartOptions []maps.Map
	widgetCodes  map[string]maps.Map // "code" => { name, ..., definition:FUNCTION CODE }
	context      *Context
	output       []string
	dbEnabled    bool

	cache bool

	// 导出
	isExporting bool
	result      interface{}
}

// 获取新引擎
func NewEngine() *Engine {
	engine := &Engine{
		chartOptions: []maps.Map{},
		widgetCodes:  map[string]maps.Map{},
		cache:        true,
	}
	engine.init()
	return engine
}

// 设置数据库是否可用
func (this *Engine) SetDBEnabled(b bool) {
	this.dbEnabled = b
}

// 设置是否导出数值
func (this *Engine) Exporting() {
	this.isExporting = true
}

// 获取导出的的内容
func (this *Engine) Result() interface{} {
	return this.result
}

// 设置上下文信息
func (this *Engine) SetContext(context *Context) {
	this.context = context

	if context.Agent != nil {
		options := map[string]interface{}{
			"isOn":    context.Agent.On,
			"id":      context.Agent.Id,
			"isLocal": context.Agent.IsLocal(),
			"name":    context.Agent.Name,
			"host":    context.Agent.Host,
			"apps": lists.Map(context.Agent.Apps, func(k int, v interface{}) interface{} {
				app := v.(*agents.AppConfig)
				return maps.Map{
					"id":   app.Id,
					"isOn": app.On,
					"name": app.Name,
					"tasks": lists.Map(app.Tasks, func(k int, v interface{}) interface{} {
						task := v.(*agents.TaskConfig)
						return maps.Map{
							"id":           task.Id,
							"isOn":         task.On,
							"name":         task.Name,
							"isBooting":    task.IsBooting,
							"isManual":     task.IsManual,
							"isScheduling": len(task.Schedule) > 0,
						}
					}),
				}
			}),
		}

		_, err := this.vm.Run(`context.agent = new agents.Agent(` + stringutil.JSONEncode(options) + `);`)
		if err != nil {
			logs.Error(err)
		}
	}

	if context.App != nil {
		options := map[string]interface{}{
			"id":   context.App.Id,
			"isOn": context.App.On,
			"name": context.App.Name,
			"tasks": lists.Map(context.App.Tasks, func(k int, v interface{}) interface{} {
				task := v.(*agents.TaskConfig)
				return maps.Map{
					"id":           task.Id,
					"isOn":         task.On,
					"name":         task.Name,
					"isBooting":    task.IsBooting,
					"isManual":     task.IsManual,
					"isScheduling": len(task.Schedule) > 0,
				}
			}),
		}

		_, err := this.vm.Run(`context.app = new agents.App(` + stringutil.JSONEncode(options) + `);`)
		if err != nil {
			logs.Error(err)
		}
	}

	if context.Item != nil {
		options := map[string]interface{}{
			"id":       context.Item.Id,
			"isOn":     context.Item.On,
			"name":     context.Item.Name,
			"interval": context.Item.IntervalDuration().Seconds(),
		}
		_, err := this.vm.Run(`context.item = new agents.Item(` + stringutil.JSONEncode(options) + `);`)
		if err != nil {
			logs.Error(err)
		}
	}

	// 时间相关
	if context.TimeType == "past" {
		context.TimeUnit = teaconfigs.TimePastUnit(context.TimePast)
		_, err := this.vm.Run(`context.timeUnit = "` + context.TimeUnit + `";`)
		if err != nil {
			logs.Error(err)
		}
	} else if context.TimeType == "range" {
		context.TimeUnit = teaconfigs.TimeUnitDay
		_, err := this.vm.Run(`context.timeUnit = "` + context.TimeUnit + `";`)
		if err != nil {
			logs.Error(err)
		}
	}

	// 可供使用的特性
	features := []string{}
	if this.dbEnabled {
		features = append(features, "db")
	}
	features = append(features, runtime.GOOS)
	features = append(features, runtime.GOARCH)
	_, err := this.vm.Run(`context.features=` + stringutil.JSONEncode(features) + `;`)
	if err != nil {
		logs.Error(err)
	}
}

// 获取Context
func (this *Engine) Context() *Context {
	return this.context
}

// 设置是否开启缓存
func (this *Engine) SetCache(cache bool) {
	this.cache = cache
}

// 初始化
func (this *Engine) init() {
	this.vm = otto.New()

	err := this.vm.Set("callConsoleLog", this.callConsoleLog)
	if err != nil {
		logs.Error(err)
	}
	_, err = this.vm.Run("console.log = callConsoleLog;")
	if err != nil {
		logs.Error(err)
	}

	this.loadLib("libs/array.js")
	this.loadLib("libs/times.js")
	this.loadLib("libs/caches.js")
	this.loadLib("libs/mutex.js")
	this.loadLib("libs/agent.values.js")
	this.loadLib("libs/colors.js")
	this.loadLib("libs/widgets.js")
	this.loadLib("libs/charts.js")
	this.loadLib("libs/charts.gauge.js")
	this.loadLib("libs/charts.html.js")
	this.loadLib("libs/charts.line.js")
	this.loadLib("libs/charts.pie.js")
	this.loadLib("libs/charts.url.js")
	this.loadLib("libs/charts.progress.js")
	this.loadLib("libs/charts.stackbar.js")
	this.loadLib("libs/charts.clock.js")
	this.loadLib("libs/charts.table.js")
	this.loadLib("libs/context.js")
	this.loadLib("libs/agent.js")

	err = this.vm.Set("callSetCache", this.callSetCache)
	if err != nil {
		logs.Error(err)
	}
	err = this.vm.Set("callGetCache", this.callGetCache)
	if err != nil {
		logs.Error(err)
	}
	err = this.vm.Set("callChartRender", this.callRenderChart)
	if err != nil {
		logs.Error(err)
	}
	err = this.vm.Set("callExecuteQuery", this.callExecuteQuery)
	if err != nil {
		logs.Error(err)
	}
}

// 运行Widget代码
func (this *Engine) RunCode(code string) error {
	_, err := this.vm.Run(`(function () {` + code + `
	
	widget.callRun();
})();`)
	return err
}

// 获取Widget中的图表对象
func (this *Engine) Charts() []maps.Map {
	return this.chartOptions
}

// 添加Output
func (this *Engine) AddOutput(output string) {
	this.output = append(this.output, output)
}

// 获取控制台输出
func (this *Engine) Output() []string {
	if this.output == nil {
		return []string{}
	}
	return this.output
}

func (this *Engine) callConsoleLog(call otto.FunctionCall) otto.Value {
	values := []string{}
	for _, v := range call.ArgumentList {
		i, err := v.Export()
		if err != nil {
			values = append(values, v.String())
		} else {
			values = append(values, stringutil.JSONEncodePretty(i))
		}
	}
	s := strings.Join(values, ", ")

	this.output = append(this.output, s)

	return otto.UndefinedValue()
}

func (this *Engine) callRenderChart(call otto.FunctionCall) otto.Value {
	obj := call.Argument(0)
	v, err := obj.Export()
	if err != nil {
		logs.Error(err)
		return otto.UndefinedValue()
	}
	m := maps.NewMap(v)

	options, err := obj.Object().Get("options")
	if err != nil {
		logs.Error(err)
	} else {
		v, err := options.Export()
		if err != nil {
			logs.Error(err)
		} else {
			m["options"] = maps.NewMap(v)
		}
	}

	this.chartOptions = append(this.chartOptions, m)
	return otto.UndefinedValue()
}

func (this *Engine) callSetCache(call otto.FunctionCall) otto.Value {
	key, err := call.Argument(0).ToString()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	key = stringutil.Md5(key)

	value, err := call.Argument(1).Export()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}

	life, err := call.Argument(2).Export()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	lifeSeconds := types.Int64(life)
	engineCache.WriteInterface([]byte(key), value, lifeSeconds)

	return otto.UndefinedValue()
}

// 获取缓存
func (this *Engine) callGetCache(call otto.FunctionCall) otto.Value {
	// 未开启cache
	if !this.cache {
		return otto.UndefinedValue()
	}

	key, err := call.Argument(0).ToString()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	key = stringutil.Md5(key)

	item := engineCache.Read([]byte(key))
	if item == nil {
		return otto.UndefinedValue()
	}
	v, err := this.vm.ToValue(item.ValueInterface)
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	return v
}

// 加载JS库文件
func (this *Engine) loadLib(file string) {
	path := teautils.WebRoot() + Tea.DS + file
	cacheKey := "libfile://" + path
	item := engineCache.Read([]byte(cacheKey))
	var code string
	if item == nil || item.ValueInterface == nil {
		var err error = nil
		code, err = files.NewFile(path).ReadAllString()
		if err != nil {
			logs.Error(err)
			return
		}
		engineCache.WriteInterface([]byte(cacheKey), code, 3600)
	} else {
		code = item.ValueInterface.(string)
	}

	_, err := this.vm.Run(code)
	if err != nil {
		logs.Error(err)
		return
	}
}

// 执行查询
func (this *Engine) callExecuteQuery(call otto.FunctionCall) otto.Value {
	arg, err := call.Argument(0).Export()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	m := maps.NewMap(arg)

	action := m.GetString("action")
	if len(action) == 0 {
		this.throw(errors.New("'action' should not be empty"))
		return otto.UndefinedValue()
	}

	if this.context == nil {
		this.throw(errors.New("'context' should not be nil"))
		return otto.UndefinedValue()
	}

	if this.context.Agent == nil {
		this.throw(errors.New("'context.agent' should not be nil"))
		return otto.UndefinedValue()
	}

	if len(this.context.Agent.Id) == 0 {
		this.throw(errors.New("'context.agent.id' should not be empty"))
		return otto.UndefinedValue()
	}

	query := teadb.NewQuery(teadb.AgentValueDAO().TableName(this.context.Agent.Id))
	if this.context.App != nil {
		query.Attr("appId", this.context.App.Id)
	}
	if this.context.Item != nil {
		query.Attr("itemId", this.context.Item.Id)
	}

	// cond
	cond := m.Get("cond")
	if cond != nil && reflect.TypeOf(cond).Kind() == reflect.Map {
		m, ok := cond.(map[string]interface{})
		if ok {
			for field, ops := range m {
				opsMap, ok := ops.(map[string]interface{})
				if ok {
					for op, v := range opsMap {
						query.Op(field, op, v)
					}
				}
			}
		}
	}

	// avgValues
	aggregationFields := m.GetSlice("aggregationFields")
	aggregationFieldStrings := []string{}
	if aggregationFields != nil {
		for _, f := range aggregationFields {
			aggregationFieldStrings = append(aggregationFieldStrings, types.String(f))
		}
	}

	// offset & size
	timeUnit := ""
	if len(this.context.TimeType) == 0 || this.context.TimeType == "default" { // 默认
		// timePast
		timePast := m.GetMap("timePast")
		if timePast != nil {
			number := timePast.GetInt64("number")
			unit := timePast.GetString("unit")
			switch unit {
			case "SECOND":
				timeUnit = "second"
			case "MINUTE":
				timeUnit = "minute"
			case "HOUR":
				timeUnit = "hour"
			case "DAY":
				timeUnit = "day"
			case "MONTH":
				timeUnit = "month"
			case "YEAR":
				timeUnit = "year"
			}
			query.Gte("createdAt", teaconfigs.TimePastUnixTimeWithUnit(number, unit))
		} else {
			query.Offset(m.GetInt("offset"))
			query.Limit(m.GetInt("size"))
		}
	} else if this.context.TimeType == "past" { // 过去N时间
		timestamp := teaconfigs.TimePastUnixTime(this.context.TimePast)
		query.Gte("createdAt", timestamp)
		timeUnit = strings.ToLower(teaconfigs.TimePastUnit(this.context.TimePast))
	} else if this.context.TimeType == "range" { // 日期范围
		timeUnit = "day"
		if len(this.context.DayFrom) > 0 && dayReg.MatchString(this.context.DayFrom) {
			match := dayReg.FindStringSubmatch(this.context.DayFrom)
			if len(match) == 4 {
				year := types.Int(match[1])
				month := types.Int(match[2])
				day := types.Int(match[3])
				t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
				query.Gte("createdAt", t.Unix())
			}
		} else {
			now := time.Now()
			query.Gte("createdAt", time.Date(now.Year(), now.Month(), now.Day()-14, 0, 0, 0, 0, time.Local).Unix())
		}
		if len(this.context.DayTo) > 0 {
			match := dayReg.FindStringSubmatch(this.context.DayTo)
			if len(match) == 4 {
				year := types.Int(match[1])
				month := types.Int(match[2])
				day := types.Int(match[3])
				t := time.Date(year, time.Month(month), day, 23, 59, 59, 0, time.Local)
				query.Lte("createdAt", t.Unix())
			}
		}
	}

	// sort
	sorts := m.Get("sorts")
	if sorts != nil {
		sortsMap, ok := sorts.([]map[string]interface{})
		if ok {
			for _, m := range sortsMap {
				for k, v := range m {
					if len(k) == 0 {
						k = "createdAt"
					}
					vInt := types.Int(v)
					if vInt < 0 {
						query.Desc(k)
					} else {
						query.Asc(k)
					}
				}
			}
		}
	}

	// 开始执行
	var result interface{} = nil
	if action == "findAll" {
		result, err = teadb.AgentValueDAO().QueryValues(query)
	} else if action == "find" {
		query.Limit(1)
		result, err = teadb.AgentValueDAO().QueryValues(query)
		if len(result.([]*agents.Value)) > 0 {
			result = result.([]*agents.Value)[0]
		}
	} else if action == "avgValues" {
		resultFields := map[string]teadb.Expr{}
		for _, s := range aggregationFieldStrings {
			resultFields[s] = teadb.NewAvgExpr("value." + s)
		}
		result, err = teadb.AgentValueDAO().GroupValuesByTime(query, timeUnit, resultFields)
	}

	if this.isExporting {
		this.result = result
	}

	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}

	jsValue, err := this.toValue(result)
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	return jsValue
}

func (this *Engine) toValue(data interface{}) (v otto.Value, err error) {
	if data == nil {
		return this.vm.ToValue(data)
	}

	// *Value
	if _, ok := data.(*agents.Value); ok {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return this.vm.ToValue(data)
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(jsonData, &m)
		if err != nil {
			logs.Error(err)
			return this.vm.ToValue(data)
		}
		return this.vm.ToValue(m)
	}

	// []*Value
	if _, ok := data.([]*agents.Value); ok {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return this.vm.ToValue(data)
		}
		m := []map[string]interface{}{}
		err = json.Unmarshal(jsonData, &m)
		if err != nil {
			logs.Error(err)
			return this.vm.ToValue(data)
		}
		return this.vm.ToValue(m)
	}

	return this.vm.ToValue(data)
}

func (this *Engine) throw(err error) {
	if err != nil {
		value, _ := this.vm.Call("new Error", nil, err.Error())
		panic(value)
	}
}
