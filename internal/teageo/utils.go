package teageo

import (
	"github.com/TeaWeb/build/internal/teamemory"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/logs"
	"github.com/oschwald/geoip2-golang"
	"net"
	"strings"
)

var ip2cityGrid = teamemory.NewGrid(32, teamemory.NewLimitCountOpt(100_0000))

// 加载Geo-City数据库
func SetupDB() {
	db, err := geoip2.Open(teautils.WebRoot() + "/resources/GeoLite2-City/GeoLite2-City.mmdb")
	if err != nil {
		logs.Error(err)
		return
	}
	DB = db
}

// 获取City数据
func IP2City(ip string, cache bool) (city *geoip2.City, err error) {
	// skip IPv6
	if strings.Contains(ip, ":") {
		return
	}

	ipObj := net.ParseIP(ip)
	if ipObj == nil {
		return
	}

	cacheKey := []byte(ip)
	if cache {
		item := ip2cityGrid.Read(cacheKey)
		if item != nil {
			if item.ValueInterface == nil {
				return nil, nil
			}
			return item.ValueInterface.(*geoip2.City), nil
		}
	}

	// 参考：https://dev.maxmind.com/geoip/geoip2/geolite2/
	city, err = DB.City(ipObj)
	if err != nil {
		return nil, err
	}

	if cache {
		ip2cityGrid.WriteInterface(cacheKey, city, 3600)
	}

	return city, nil
}
