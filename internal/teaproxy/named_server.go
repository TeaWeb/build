package teaproxy

import "github.com/TeaWeb/build/internal/teaconfigs"

// 域名和服务映射
type NamedServer struct {
	Name   string                   // 匹配后的域名
	Server *teaconfigs.ServerConfig // 匹配后的服务配置
}
