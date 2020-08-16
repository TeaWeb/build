package teautils

import (
	"github.com/iwind/TeaGo/maps"
	json "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
)

// 通过YAML把map转换为object
func MapToObjectYAML(fromMap map[string]interface{}, toPtr interface{}) error {
	data, err := yaml.Marshal(fromMap)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, toPtr)
	return err
}

// 通过JSON把map转换为object
func MapToObjectJSON(fromMap map[string]interface{}, toPtr interface{}) error {
	data, err := json.Marshal(fromMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, toPtr)
	return err
}

// 通过JSON把object转换为map
func ObjectToMapJSON(fromPtr interface{}, toMap *map[string]interface{}) error {
	data, err := json.Marshal(fromPtr)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, toMap)
	return err
}

// 获取所有的键值
func MapKeys(s maps.Map) (keys []string) {
	for k := range s {
		keys = append(keys, k)
	}
	return
}
