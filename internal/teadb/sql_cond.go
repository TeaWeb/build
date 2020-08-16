package teadb

// SQL条件
type SQLCond struct {
	Expr   string
	Params map[string]interface{}
}
