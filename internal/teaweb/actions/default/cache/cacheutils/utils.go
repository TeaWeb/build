package cacheutils

import "github.com/TeaWeb/build/internal/teaconfigs"

// 判断缓存策略是否正在被使用
func IsPolicyUsed(policyFilename string) bool {
	serverList, _ := teaconfigs.SharedServerList()
	if serverList == nil {
		return false
	}

	for _, server := range serverList.FindAllServers() {
		if server.CachePolicy == policyFilename {
			return true
		}

		for _, location := range server.Locations {
			if location.CachePolicy == policyFilename {
				return true
			}
		}
	}
	return false
}
