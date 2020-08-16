package proxyutils

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teastats"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"strings"
)

// 刷新服务统计
func ReloadServerStats(serverId string) {
	server := teaconfigs.NewServerConfigFromId(serverId)
	if server == nil || !server.On {
		teastats.RestartServerFilters(serverId, nil)
		return
	}

	codes := server.StatItems
	for _, board := range []*teaconfigs.Board{server.RealtimeBoard, server.StatBoard} {
		if board == nil {
			continue
		}
		for _, c := range board.Charts {
			_, chart := c.FindChart()
			if chart == nil || !chart.On {
				continue
			}
			for _, r := range chart.Requirements {
				if lists.ContainsString(codes, r) {
					continue
				}
				codes = append(codes, r)
			}
		}
	}
	teastats.RestartServerFilters(serverId, codes)
}

// 检查图表是否有更新
func CheckChartChanges() bool {
	dir := files.NewFile(Tea.Root + "/web/libs/widgets")
	if !dir.Exists() {
		return false
	}

	for _, file := range dir.List() {
		if !strings.HasPrefix(file.Name(), "widget.") {
			continue
		}
		data, err := file.ReadAll()
		if err != nil {
			logs.Error(err)
			continue
		}

		// 对应配置目录
		configFile := files.NewFile(Tea.ConfigFile("widgets/" + file.Name()))
		if !configFile.Exists() {
			return true
		}

		configData, err := configFile.ReadAll()
		if err != nil {
			logs.Error(err)
			continue
		}

		if bytes.Compare(data, configData) != 0 {
			return true
		}
	}

	return false
}

// 应用图表更新
func ApplyChartChanges() error {
	dir := files.NewFile(Tea.Root + "/web/libs/widgets")
	if !dir.Exists() {
		return nil
	}

	for _, file := range dir.List() {
		if !strings.HasPrefix(file.Name(), "widget.") {
			continue
		}
		data, err := file.ReadAll()
		if err != nil {
			return err
		}

		// 对应配置目录
		configFile := files.NewFile(Tea.ConfigFile("widgets/" + file.Name()))
		if !configFile.Exists() {
			err := configFile.Write(data)
			if err != nil {
				return err
			}
			logs.Println("[proxy]apply the updates for widget '" + configFile.Name() + "'")
			continue
		}

		configData, err := configFile.ReadAll()
		if err != nil {
			return err
		}

		if bytes.Compare(data, configData) != 0 {
			err := configFile.Write(data)
			if err != nil {
				return err
			}
			logs.Println("[proxy]apply the updates for widget '" + configFile.Name() + "'")
			continue
		}
	}

	return nil
}
