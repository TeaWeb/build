package teadb

import (
	"reflect"
	"regexp"
	"strconv"
)

var holderReg = regexp.MustCompile(`:\w+`)

// SQL参数holder
type SQLParamsHolder struct {
	Params map[string]interface{} // holder => value
	Args   []interface{}
	index  int
	driver string
}

// 获取新对象
func NewSQLParamsHolder(driver string) *SQLParamsHolder {
	return &SQLParamsHolder{
		Params: map[string]interface{}{},
		index:  0,
		driver: driver,
	}
}

// 添加参数值
func (this *SQLParamsHolder) Add(value interface{}) (holder string) {
	holder = ":HOLDER" + strconv.Itoa(this.index)
	this.Params[holder] = value
	this.index++
	return
}

// 添加自定义参数值
func (this *SQLParamsHolder) AddHolder(holder string, value interface{}) {
	this.Params[":"+holder] = value
}

// 参加slice参数值
func (this *SQLParamsHolder) AddSlice(s interface{}) (holder string) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Slice {
		return this.Add(s)
	}

	l := v.Len()
	holder = "("
	for i := 0; i < l; i++ {
		h := ":HOLDER" + strconv.Itoa(this.index)
		if i == 0 {
			holder += h
		} else {
			holder += ", " + h
		}
		this.Params[h] = v.Index(i).Interface()
		this.index++
	}

	holder += ")"

	return
}

// 分析数据
func (this *SQLParamsHolder) Parse(sqlString string) string {
	if len(this.Params) == 0 {
		return sqlString
	}
	index := 0
	sqlString = holderReg.ReplaceAllStringFunc(sqlString, func(s string) string {
		v, _ := this.Params[s]
		this.Args = append(this.Args, v)

		switch this.driver {
		case "mysql":
			return "?"
		case "postgres":
			index++
			return "$" + strconv.Itoa(index)
		}
		return "?"
	})
	return sqlString
}
