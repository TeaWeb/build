package agents

import (
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"reflect"
	"sync"
)

// 获取所有的数据源信息
var allDataSources = []maps.Map{}
var allDataSourcesLocker = sync.Mutex{}

func AllDataSources() []maps.Map {
	return allDataSources
}

// 注册内置的数据源信息
func RegisterAllDataSources() {
	// basic
	RegisterDataSource(NewScriptSource(), SourceCategoryBasic)
	RegisterDataSource(NewWebHookSource(), SourceCategoryBasic)
	RegisterDataSource(NewFileSource(), SourceCategoryBasic)

	// common
	RegisterDataSource(NewCPUSource(), SourceCategoryCommon)
	RegisterDataSource(NewMemorySource(), SourceCategoryCommon)
	RegisterDataSource(NewLoadSource(), SourceCategoryCommon)
	RegisterDataSource(NewNetworkSource(), SourceCategoryCommon)
	RegisterDataSource(NewDiskSource(), SourceCategoryCommon)
	RegisterDataSource(NewIOStatSource(), SourceCategoryCommon)
	RegisterDataSource(NewConnectionsSource(), SourceCategoryCommon)
	RegisterDataSource(NewProcessesSource(), SourceCategoryCommon)
	RegisterDataSource(NewDateSource(), SourceCategoryCommon)
	RegisterDataSource(NewURLConnectivitySource(), SourceCategoryCommon)
	RegisterDataSource(NewSocketConnectivitySource(), SourceCategoryCommon)
	RegisterDataSource(NewDNSSource(), SourceCategoryCommon)
	RegisterDataSource(NewPingSource(), SourceCategoryCommon)
	RegisterDataSource(NewAppProcessesSource(), SourceCategoryCommon)
	RegisterDataSource(NewFileChangeSource(), SourceCategoryCommon)
	RegisterDataSource(NewMySQLSource(), SourceCategoryCommon)
	RegisterDataSource(NewPostgreSQLSource(), SourceCategoryCommon)
	RegisterDataSource(NewDockerSource(), SourceCategoryCommon)
	RegisterDataSource(NewTeaWebSource(), SourceCategoryCommon)
	RegisterDataSource(NewNginxStatusSource(), SourceCategoryCommon)

	// plugin
}

// 单个数据源信息
func RegisterDataSource(dataSource SourceInterface, category SourceCategory) {
	allDataSourcesLocker.Lock()
	defer allDataSourcesLocker.Unlock()

	m := maps.Map{
		"name":        dataSource.Name(),
		"code":        dataSource.Code(),
		"description": dataSource.Description(),
		"type":        reflect.TypeOf(dataSource).Elem(),
		"instance":    dataSource,
		"category":    category,
	}
	allDataSources = append(allDataSources, m)
}

// 查找单个数据源信息
func FindDataSource(code string) maps.Map {
	for _, summary := range AllDataSources() {
		if summary["code"] == code {
			return summary
		}
	}
	return nil
}

// 查找单个数据源实例
func FindDataSourceInstance(code string, options map[string]interface{}) SourceInterface {
	for _, summary := range AllDataSources() {
		if summary["code"] == code {
			instance := reflect.New(summary["type"].(reflect.Type)).Interface().(SourceInterface)
			if options != nil {
				err := teautils.MapToObjectJSON(options, instance)
				if err != nil {
					logs.Error(err)
				}
			}
			return instance
		}
	}
	return nil
}

// 将Source转换为Map
func ConvertSourceToMap(source SourceInterface) map[string]interface{} {
	m := map[string]interface{}{}
	err := teautils.ObjectToMapJSON(source, &m)
	if err != nil {
		logs.Error(err)
	}
	return m
}
