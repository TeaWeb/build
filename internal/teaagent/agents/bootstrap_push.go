package teaagents

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/TeaWeb/build/internal/teaagent/agentconst"
	"github.com/TeaWeb/build/internal/teautils/logbuffer"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

// 向Master同步事件
var db *logbuffer.Buffer

func pushEvents() {
	db = logbuffer.NewBuffer(Tea.Root + "/logs/agent.event")

	// 读取本地数据库日志并发送到Master
	go func() {
		for {
			data, err := db.Read()
			if err != nil {
				logs.Println("[push]" + err.Error())
				time.Sleep(1 * time.Second)
				continue
			}

			if len(data) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			// Push到Master服务器
			req, err := http.NewRequest(http.MethodPost, connectConfig.Master+"/api/agent/push", bytes.NewReader(data))
			if err != nil {
				logs.Println("[push]" + err.Error())
				continue
			}

			err = func() error {
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("User-Agent", "TeaWeb Agent")
				req.Header.Set("Tea-Agent-Id", connectConfig.Id)
				req.Header.Set("Tea-Agent-Key", connectConfig.Key)
				req.Header.Set("Tea-Agent-Version", agentconst.AgentVersion)
				req.Header.Set("Tea-Agent-Os", runtime.GOOS)
				req.Header.Set("Tea-Agent-Arch", runtime.GOARCH)
				resp, err := HTTPClient.Do(req)

				if err != nil {
					return err
				}
				defer func() {
					_ = resp.Body.Close()
				}()
				if resp.StatusCode != 200 {
					return errors.New("response code '" + strconv.Itoa(resp.StatusCode) + "'")
				}

				respBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				respJSON := maps.Map{}
				err = json.Unmarshal(respBody, &respJSON)
				if err != nil {
					return err
				}

				if respJSON.GetInt("code") != 200 {
					return errors.New("error response from master:" + string(respBody))
				}

				return nil
			}()
			if err != nil {
				logs.Println("[push]", err.Error())
				time.Sleep(60 * time.Second)
			}
		}
	}()

	// 读取日志并写入到本地数据库
	logId := time.Now().UnixNano()
	for {
		event := <-eventQueue

		if runningAgent.Id != "local" {
			// 进程事件
			if event, found := event.(*ProcessEvent); found {
				if event.EventType == ProcessEventStdout || event.EventType == ProcessEventStderr {
					logs.Println("[" + findTaskName(event.TaskId) + "]" + event.Data)
				} else if event.EventType == ProcessEventStart {
					logs.Println("[" + findTaskName(event.TaskId) + "]start")
				} else if event.EventType == ProcessEventStop {
					logs.Println("[" + findTaskName(event.TaskId) + "]stop")
				}
			}
		}

		jsonData, err := event.AsJSON()
		if err != nil {
			logs.Println("error:", err.Error())
			continue
		}

		if db != nil {
			logId++
			_, err = db.Write(jsonData)
			if err != nil {
				logs.Println("[ERROR]" + err.Error())
			}
		}
	}
}
