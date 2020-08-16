package utils

// 子菜单定义
type Menu struct {
	Id               string      `json:"id"`
	Name             string      `json:"name"`
	Items            []*MenuItem `json:"items"`
	IsActive         bool        `json:"isActive"`
	AlwaysActive     bool        `json:"alwaysActive"`
	Index            int         `json:"index"`
	CountNormalItems int         `json:"countNormalItems"`
}

// 获取新对象
func NewMenu() *Menu {
	return &Menu{
		Items: []*MenuItem{},
	}
}

// 添加菜单项
func (this *Menu) Add(name string, subName string, url string, isActive bool) *MenuItem {
	item := &MenuItem{
		Name:     name,
		SubName:  subName,
		URL:      url,
		IsActive: isActive,
	}
	this.CountNormalItems ++
	this.Items = append(this.Items, item)

	if isActive {
		this.IsActive = true
	}

	return item
}

// 添加特殊菜单项，不计数
func (this *Menu) AddSpecial(name string, subName string, url string, isActive bool) *MenuItem {
	item := this.Add(name, subName, url, isActive)
	this.CountNormalItems --
	return item
}
