package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
	"time"
)

// 查询对象
type Query struct {
	table        string
	offset       int
	size         int
	operandList  *OperandList
	sortFields   []*SortField
	debug        bool
	timeout      time.Duration
	resultFields []string
	fieldMapping func(field string) string
}

// 构造新查询
func NewQuery(table string) *Query {
	if sharedDBType != "mongo" {
		// 关系型数据库使用下划线分隔
		table = strings.Replace(table, ".", "_", -1)
	}

	query := &Query{
		table: table,
	}
	query.Init()
	return query
}

func (this *Query) Init() *Query {
	this.offset = -1
	this.size = -1
	this.operandList = NewOperandList()
	return this
}

func (this *Query) Table(table string) *Query {
	this.table = table
	return this
}

func (this *Query) Debug() *Query {
	this.debug = true
	return this
}

func (this *Query) Timeout(timeout time.Duration) *Query {
	this.timeout = timeout
	return this
}

func (this *Query) Result(field ...string) *Query {
	this.resultFields = append(this.resultFields, field...)
	return this
}

func (this *Query) Attr(field string, value interface{}) *Query {
	if types.IsSlice(value) {
		return this.Op(field, OperandIn, value)
	} else {
		return this.Op(field, OperandEq, value)
	}
}

func (this *Query) Op(field string, operandCode OperandCode, value interface{}) *Query {
	this.operandList.Add(field, NewOperand(operandCode, value))
	return this
}

func (this *Query) Or(fieldValues []*OperandList) *Query {
	return this.Op("", OperandOr, fieldValues)
}

func (this *Query) Not(field string, value interface{}) *Query {
	if types.IsSlice(value) {
		return this.Op(field, OperandNotIn, value)
	} else {
		return this.Op(field, OperandNeq, value)
	}
}

func (this *Query) Lt(field string, value interface{}) *Query {
	return this.Op(field, OperandLt, value)
}

func (this *Query) Lte(field string, value interface{}) *Query {
	return this.Op(field, OperandLte, value)
}

func (this *Query) Gt(field string, value interface{}) *Query {
	return this.Op(field, OperandGt, value)
}

func (this *Query) Gte(field string, value interface{}) *Query {
	return this.Op(field, OperandGte, value)
}

// SQL查询专用
func (this *Query) sqlCond(expr string, params map[string]interface{}) *Query {
	this.Op("", operandSQLCond, &SQLCond{
		Expr:   expr,
		Params: params,
	})
	return this
}

func (this *Query) Asc(field string) *Query {
	if this.hasSortField(field) {
		this.removeSortField(field)
	}
	this.sortFields = append(this.sortFields, &SortField{
		Name: field,
		Type: SortAsc,
	})
	return this
}

func (this *Query) Desc(field string) *Query {
	if this.hasSortField(field) {
		this.removeSortField(field)
	}
	this.sortFields = append(this.sortFields, &SortField{
		Name: field,
		Type: SortDesc,
	})
	return this
}

func (this *Query) Offset(offset int) *Query {
	this.offset = offset
	return this
}

func (this *Query) Limit(size int) *Query {
	this.size = size
	return this
}

func (this *Query) Node() *Query {
	node := teaconfigs.SharedNodeConfig()
	if node != nil {
		this.Attr("nodeId", node.Id)
	} else {
		this.Attr("nodeId", "")
	}
	return this
}

func (this *Query) FindOne(modelPtr interface{}) (interface{}, error) {
	return sharedDriver.FindOne(this, modelPtr)
}

func (this *Query) FindOnes(modelPtr interface{}) ([]interface{}, error) {
	return sharedDriver.FindOnes(this, modelPtr)
}

func (this *Query) InsertOne(modelPtr interface{}) error {
	return sharedDriver.InsertOne(this.table, modelPtr)
}

func (this *Query) InsertOnes(modelPtrSlice interface{}) error {
	return sharedDriver.InsertOnes(this.table, modelPtrSlice)
}

func (this *Query) Delete() error {
	return sharedDriver.DeleteOnes(this)
}

func (this *Query) Count() (int64, error) {
	return sharedDriver.Count(this)
}

func (this *Query) Sum(field string) (float64, error) {
	return sharedDriver.Sum(this, field)
}

func (this *Query) Min(field string) (float64, error) {
	return sharedDriver.Min(this, field)
}

func (this *Query) Max(field string) (float64, error) {
	return sharedDriver.Max(this, field)
}

func (this *Query) Avg(field string) (float64, error) {
	return sharedDriver.Avg(this, field)
}

func (this *Query) Group(field string, result map[string]Expr) ([]maps.Map, error) {
	return sharedDriver.Group(this, field, result)
}

func (this *Query) FieldMap(mapping func(field string) string) *Query {
	this.fieldMapping = mapping
	return this
}

func (this *Query) hasSortField(field string) bool {
	for _, f := range this.sortFields {
		if f.Name == field {
			return true
		}
	}
	return false
}

func (this *Query) removeSortField(field string) {
	result := []*SortField{}
	for _, f := range this.sortFields {
		if f.Name != field {
			result = append(result, f)
		}
	}
	this.sortFields = result
}
