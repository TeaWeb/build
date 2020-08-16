package teautils

import (
	"reflect"
)

// 拷贝同类型struct指针对象中的字段
func CopyStructObject(destPtr, sourcePtr interface{}) {
	value := reflect.ValueOf(destPtr)
	value2 := reflect.ValueOf(sourcePtr)

	countFields := value2.Elem().NumField()
	for i := 0; i < countFields; i++ {
		v := value2.Elem().Field(i)
		if !v.IsValid() || !v.CanSet() {
			continue
		}
		value.Elem().Field(i).Set(v)
	}
}
