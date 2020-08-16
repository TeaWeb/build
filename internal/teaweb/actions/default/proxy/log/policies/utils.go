package policies

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/maps"
)

func FindAllUsingPolicy(policyId string) []maps.Map {
	// 正在使用此策略的项目
	configItems := []maps.Map{}
	serverList, _ := teaconfigs.SharedServerList()
	if serverList != nil {
		for _, server := range serverList.FindAllServers() {
			if len(server.AccessLog) == 0 {
				continue
			}

			for _, accessLog := range server.AccessLog {
				if !accessLog.On || !accessLog.ContainsStoragePolicy(policyId) {
					continue
				}
				configItems = append(configItems, maps.Map{
					"type":   "server",
					"server": server.Description,
					"link":   "/proxy/detail?serverId=" + server.Id,
				})
			}

			for _, location := range server.Locations {
				if len(location.AccessLog) == 0 {
					continue
				}

				for _, accessLog := range location.AccessLog {
					if !accessLog.On || !accessLog.ContainsStoragePolicy(policyId) {
						continue
					}
					configItems = append(configItems, maps.Map{
						"type":     "location",
						"server":   server.Description,
						"location": location.Pattern,
						"link":     "/proxy/locations/detail?serverId=" + server.Id + "&locationId=" + location.Id,
					})
				}
			}
		}
	}

	return configItems
}
