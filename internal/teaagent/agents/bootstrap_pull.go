package teaagents

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaagent/agentconst"
	"github.com/TeaWeb/build/internal/teaagent/agentutils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"time"
)

// 从主服务器同步数据
func pullEvents() error {
	//logs.Println("pull events ...", connectConfig.Master+"/api/agent/pull")
	master := connectConfig.Master
	if len(master) == 0 {
		return errors.New("'master' should not be empty")
	}
	req, err := http.NewRequest(http.MethodGet, master+"/api/agent/pull", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "TeaWeb Agent")
	req.Header.Set("Tea-Agent-Id", connectConfig.Id)
	req.Header.Set("Tea-Agent-Key", connectConfig.Key)
	req.Header.Set("Tea-Agent-Version", agentconst.AgentVersion)
	req.Header.Set("Tea-Agent-Os", runtime.GOOS)
	req.Header.Set("Tea-Agent-OsName", retrieveOSNameBase64())
	req.Header.Set("Tea-Agent-Arch", runtime.GOARCH)
	req.Header.Set("Tea-Agent-Nano", fmt.Sprintf("%d", time.Now().UnixNano()))
	connectingFailed := false
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// 握手配置
				conn, err := (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 0,
				}).DialContext(ctx, network, addr)
				if err != nil {
					connectingFailed = true
				} else {
					// 恢复连接
					if connectionIsBroken {
						connectionIsBroken = false
						initConnection()
					}
				}
				return conn, err
			},
			IdleConnTimeout: 65 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	defer agentutils.CloseHTTPClient(client)
	resp, err := client.Do(req)
	if err != nil {
		if connectingFailed {
			connectionIsBroken = true
			return err
		}

		// 恢复连接
		if connectionIsBroken {
			connectionIsBroken = false
			initConnection()
		}

		// 如果是超时的则不提示，因为长连接依赖超时设置
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()

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
	events := dataMap.Get("events")
	if events == nil || reflect.TypeOf(events).Kind() != reflect.Slice {
		return nil
	}

	eventsValue := reflect.ValueOf(events)
	count := eventsValue.Len()
	for i := 0; i < count; i++ {
		event := eventsValue.Index(i).Interface()
		if event == nil || reflect.TypeOf(event).Kind() != reflect.Map {
			continue
		}
		eventMap := maps.NewMap(event)
		name := eventMap.GetString("name")
		switch name {
		case "UPDATE_AGENT":
			go downloadConfig()
		case "REMOVE_AGENT":
			os.Exit(0)
		case "ADD_APP":
			go downloadConfig()
		case "UPDATE_APP":
			go downloadConfig()
		case "REMOVE_APP":
			go downloadConfig()
		case "ADD_TASK":
			go downloadConfig()
		case "UPDATE_TASK":
			go downloadConfig()
		case "REMOVE_TASK":
			go downloadConfig()
		case "RUN_TASK":
			eventDataMap := eventMap.GetMap("data")
			if eventDataMap != nil {
				taskId := eventDataMap.GetString("taskId")
				appConfig, taskConfig := runningAgent.FindTask(taskId)
				if taskConfig == nil {
					logs.Println("error:no task with id '" + taskId + " found")
				} else {
					task := NewTask(appConfig.Id, taskConfig)
					go task.RunLog()
				}
			} else {
				logs.Println("invalid event data: should be a map")
			}
		case "ADD_ITEM":
			go downloadConfig()
		case "UPDATE_ITEM":
			go downloadConfig()
		case "DELETE_ITEM":
			go downloadConfig()
		case "RUN_ITEM":
			eventDataMap := eventMap.GetMap("data")
			if eventDataMap != nil {
				itemId := eventDataMap.GetString("itemId")
				found := false
				logs.Println("run item " + itemId)
				for _, item := range runningItems {
					if item.config.Id == itemId {
						found = true
						go func(item *Item) {
							t := time.Now()
							value, err := item.Run()
							costMs := time.Since(t).Seconds() * 1000
							PushEvent(NewItemEvent(runningAgent.Id, item.appId, item.config.Id, value, err, t.Unix(), costMs))
						}(item)
						break
					}
				}
				if !found {
					logs.Println("error:item with id '" + itemId + "' not found")
				}
			}
		}
	}

	return nil
}
