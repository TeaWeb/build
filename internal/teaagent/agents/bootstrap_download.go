package teaagents

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaagent/agentconst"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime"
)

// 下载配置
func downloadConfig() error {
	// 本地
	if connectConfig.Id == "local" {
		loadLocalConfig()

		return nil
	}

	// 远程的
	master := connectConfig.Master
	if len(master) == 0 {
		return errors.New("'master' should not be empty")
	}
	req, err := http.NewRequest(http.MethodGet, master+"/api/agent", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "TeaWeb Agent")
	req.Header.Set("Tea-Agent-Id", connectConfig.Id)
	req.Header.Set("Tea-Agent-Key", connectConfig.Key)
	req.Header.Set("Tea-Agent-Version", agentconst.AgentVersion)
	req.Header.Set("Tea-Agent-Os", runtime.GOOS)
	req.Header.Set("Tea-Agent-Arch", runtime.GOARCH)
	req.Header.Set("Tea-Agent-Group-Key", connectConfig.GroupKey)

	hostName, _ := os.Hostname()
	req.Header.Set("Tea-Agent-Hostname", base64.StdEncoding.EncodeToString([]byte(hostName)))

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid status response from master '" + fmt.Sprintf("%d", resp.StatusCode) + "'")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respMap := maps.Map{}
	err = json.Unmarshal(data, &respMap)
	if err != nil {
		return err
	}

	if respMap == nil {
		return errors.New("response data should not be nil")
	}

	if respMap.GetInt("code") != 200 {
		return errors.New("invalid response from master:" + string(data))
	}

	jsonData := respMap.Get("data")
	if jsonData == nil || reflect.TypeOf(jsonData).Kind() != reflect.Map {
		return errors.New("response json data should be a map")
	}

	dataMap := maps.NewMap(jsonData)
	config := dataMap.GetString("config")

	agent := &agents.AgentConfig{}
	err = yaml.Unmarshal([]byte(config), agent)
	if err != nil {
		return err
	}

	if len(agent.Id) == 0 {
		return errors.New("invalid agent id")
	}

	err = agent.Validate()
	if err != nil {
		return err
	}

	// 保存
	agentsDir := files.NewFile(Tea.ConfigFile("agents/"))
	if !agentsDir.IsDir() {
		err = agentsDir.Mkdir()
		if err != nil {
			return err
		}
	}
	agentFile := files.NewFile(Tea.ConfigFile("agents/agent." + agent.Id + ".conf"))
	err = agentFile.WriteString(config)
	if err != nil {
		return err
	}

	runningAgent = agent

	// 保存ID和Key信息
	if connectConfig.Id != agent.Id || connectConfig.Key != agent.Key {
		connectConfig.Id = agent.Id
		connectConfig.Key = agent.Key
		err = connectConfig.Save()
		if err != nil {
			return errors.New("save config error: " + err.Error())
		}
	}

	if !isBooting {
		// 定时任务
		scheduleTasks()

		// 监控项数据
		scheduleItems()
	}

	return nil
}
