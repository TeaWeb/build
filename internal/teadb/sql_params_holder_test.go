package teadb

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestSQLParamsHolder_AddParam(t *testing.T) {
	holder := NewSQLParamsHolder("mysql")
	t.Log(holder.Add(123))
	t.Log(holder.Add("456"))
	logs.PrintAsJSON(holder.Params, t)
	t.Log(holder.Parse("SELECT * FROM t WHERE a=:HOLDER0 AND b=:HOLDER1"))
	logs.PrintAsJSON(holder.Args, t)
}
