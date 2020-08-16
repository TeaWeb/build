package utils

// 菜单项
type MenuItem struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	SubName    string `json:"subName"` // 副标题
	SupName    string `json:"supName"` // 头标
	URL        string `json:"url"`
	IsActive   bool   `json:"isActive"`
	Icon       string `json:"icon"`
	IsSortable bool   `json:"isSortable"`
	SubColor   string `json:"subColor"`
}
