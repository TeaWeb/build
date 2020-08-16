package agents

import (
	"database/sql"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/maps"
	"strconv"
)

// MySQL SQL
type MySQLSource struct {
	Source `yaml:",inline"`

	Addr           string `yaml:"addr" json:"addr"`
	Username       string `yaml:"username" json:"username"`
	Password       string `yaml:"password" json:"password"`
	DatabaseName   string `yaml:"databaseName" json:"databaseName"`
	TimeoutSeconds int    `yaml:"timeoutSeconds" json:"timeoutSeconds"`
	SQL            string `yaml:"sql" json:"sql"`

	db *sql.DB
}

// 获取新对象
func NewMySQLSource() *MySQLSource {
	return &MySQLSource{}
}

// 名称
func (this *MySQLSource) Name() string {
	return "MySQL SQL"
}

// 代号
func (this *MySQLSource) Code() string {
	return "mysql.sql"
}

// 描述
func (this *MySQLSource) Description() string {
	return "通过SQL语句从MySQL中获取信息，可以使用数据格式单行或多行来控制返回数据的行数"
}

// 执行
func (this *MySQLSource) Execute(params map[string]string) (value interface{}, err error) {
	if this.TimeoutSeconds <= 0 {
		this.TimeoutSeconds = 5
	}

	var db *sql.DB
	if this.db != nil {
		db = this.db
	} else {
		// 超时时间使用ms，并除以3，因为会自动尝试3次连接
		db, err = sql.Open("mysql", this.Username+":"+this.Password+"@tcp("+this.Addr+")/"+this.DatabaseName+"?timeout="+strconv.Itoa(this.TimeoutSeconds*1000/3)+"ms")
		if err != nil {
			return nil, err
		}
		db.SetMaxIdleConns(1)
		this.db = db
	}

	rows, err := db.Query(this.SQL)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := []map[string]interface{}{}
	for rows.Next() {
		values := []interface{}{}
		for range cols {
			var ptr interface{} = nil
			values = append(values, &ptr)
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		result := map[string]interface{}{}
		for index, v := range values {
			value := *(v.(*interface{}))
			if b, ok := value.([]byte); ok {
				result[cols[index]] = string(b)
			} else {
				result[cols[index]] = value
			}
		}
		results = append(results, result)
	}

	if this.DataFormat == SourceDataFormatSingeLine {
		if len(results) > 0 {
			return results[0], nil
		}
		return maps.Map{}, nil
	}

	return results, nil
}

// 表单信息
func (this *MySQLSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	{
		group := form.NewGroup()

		{
			field := forms.NewTextField("地址", "Host")
			field.Comment = "带端口的地址，比如 127.0.0.1:3306"
			field.IsRequired = true
			field.Code = "addr"
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入MySQL数据库地址");
}`
			group.Add(field)
		}

		{
			field := forms.NewTextField("用户名", "Username")
			field.IsRequired = true
			field.Code = "username"
			field.Value = "root"
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入用户名");
}`
			group.Add(field)
		}

		{
			field := forms.NewTextField("密码", "Password")
			field.Code = "password"
			group.Add(field)
		}

		{
			field := forms.NewTextField("数据库名称", "Database")
			field.Code = "databaseName"
			field.IsRequired = true
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入数据库名称");
}
`
			group.Add(field)
		}

		{
			field := forms.NewTextField("连接超时时间", "Timeout")
			field.MaxLength = 4
			field.Attr("style", "width:5em")
			field.RightLabel = "秒"
			field.Value = 10
			field.ValidateCode = `
var intValue = parseInt(value);
if (isNaN(intValue)) {
	throw new Error("超时时间需要是一个整数");
}

return intValue;
`
			field.Code = "timeoutSeconds"

			group.Add(field)
		}

		{
			field := forms.NewTextBox("SQL语句", "SQL")
			field.Code = "sql"
			field.Rows = 3
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入SQL语句");
}`
			field.Value = "SELECT 1"
			group.Add(field)
		}
	}

	return form
}

// 显示信息
func (this *MySQLSource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
<tr>
	<td>地址<em>（Host）</em></td>
	<td>{{source.addr}}</td>
</tr>
<tr>
	<td>用户名<em>（Username）</em></td>
	<td>{{source.username}}</td>
</tr>
<tr>
	<td>密码<em>（Password）</em></td>
	<td><span v-if="source.password.length == 0">没有设置</span><span v-if="source.password.length > 0">{{source.password}}</span></td>
</tr>
<tr>
	<td>数据库名称<em>（Database）</em></td>
	<td>{{source.databaseName}}</td>
</tr>
<tr>
	<td>连接超时时间<em>（Timeout）</em></td>
	<td>{{source.timeoutSeconds}}s</td>
</tr>
<tr>
	<td>SQL<em>（Host）</em></td>
	<td>{{source.sql}}</td>
</tr>
`
	return p
}

// 变量
func (this *MySQLSource) Variables() []*SourceVariable {
	return []*SourceVariable{

	}
}

// 阈值
func (this *MySQLSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	return result
}

// 图表
func (this *MySQLSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}
	return charts
}

// 停止
func (this *MySQLSource) Stop() error {
	if this.db != nil {
		return this.db.Close()
	}
	return nil
}
