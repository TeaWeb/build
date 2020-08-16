package wafutils

import "github.com/TeaWeb/build/internal/teaconfigs"

// 判断缓存策略是否正在被使用
func IsPolicyUsed(wafId string) bool {
	serverList, _ := teaconfigs.SharedServerList()
	if serverList == nil {
		return false
	}

	for _, server := range serverList.FindAllServers() {
		if server.WafId == wafId {
			return true
		}

		for _, location := range server.Locations {
			if location.WafId == wafId {
				return true
			}
		}
	}
	return false
}
