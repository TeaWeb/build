package teaconfigs

import "github.com/iwind/TeaGo/lists"

// 角色定义
type NodeRole = string

const (
	NodeRoleMaster = "MASTER"
	NodeRoleSlave  = "SLAVE"
)

// 所有角色
func AllNodeRoles() []string {
	return []string{NodeRoleMaster, NodeRoleSlave}
}

// 判断某个角色是否存在
func ExistNodeRole(role NodeRole) bool {
	return lists.ContainsString(AllNodeRoles(), role)
}
