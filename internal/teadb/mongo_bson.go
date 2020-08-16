package teadb

import (
	"encoding/json"
	"github.com/iwind/TeaGo/maps"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BSONArrayBytes(data []byte) (interface{}, error) {
	m := []interface{}{}
	err := json.Unmarshal(data, &m)
	return m, err
}

func BSONDecode(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case primitive.D:
		m := map[string]interface{}{}
		for k, v1 := range v.Map() {
			v2, err := BSONDecode(v1)
			if err != nil {
				return nil, err
			}
			m[k] = v2
		}
		return m, nil
	case primitive.A:
		arr := []interface{}{}
		for _, v1 := range v {
			v2, err := BSONDecode(v1)
			if err != nil {
				return nil, err
			}
			arr = append(arr, v2)
		}
		return arr, nil
	case primitive.ObjectID:
		return v.Hex(), nil
	case primitive.Null:
		return nil, nil
	case primitive.Binary:
		return v.Data, nil
	case primitive.Timestamp:
		return v.T, nil
	case map[string]interface{}:
		for itemKey, itemValue := range v {
			r, err := BSONDecode(itemValue)
			if err != nil {
				return nil, err
			}
			v[itemKey] = r
		}
	case maps.Map:
		for itemKey, itemValue := range v {
			r, err := BSONDecode(itemValue)
			if err != nil {
				return nil, err
			}
			v[itemKey] = r
		}
	}

	return value, nil
}
