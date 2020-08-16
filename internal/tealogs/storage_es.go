package tealogs

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/logs"
	"github.com/pquerna/ffjson/ffjson"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ElasticSearch存储策略
type ESStorage struct {
	Storage `yaml:", inline"`

	Endpoint    string `yaml:"endpoint" json:"endpoint"`
	Index       string `yaml:"index" json:"index"`
	MappingType string `yaml:"mappingType" json:"mappingType"`
	Username    string `yaml:"username" json:"username"`
	Password    string `yaml:"password" json:"password"`
}

// 开启
func (this *ESStorage) Start() error {
	if len(this.Endpoint) == 0 {
		return errors.New("'endpoint' should not be nil")
	}
	if !regexp.MustCompile(`(?i)^(http|https)://`).MatchString(this.Endpoint) {
		this.Endpoint = "http://" + this.Endpoint
	}
	if len(this.Index) == 0 {
		return errors.New("'index' should not be nil")
	}
	if len(this.MappingType) == 0 {
		return errors.New("'mappingType' should not be nil")
	}
	return nil
}

// 写入日志
func (this *ESStorage) Write(accessLogs []*accesslogs.AccessLog) error {
	if len(accessLogs) == 0 {
		return nil
	}

	bulk := &strings.Builder{}
	id := time.Now().UnixNano()
	indexName := this.FormatVariables(this.Index)
	typeName := this.FormatVariables(this.MappingType)
	for _, accessLog := range accessLogs {
		id++
		opData, err := ffjson.Marshal(map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
				"_type":  typeName,
				"_id":    fmt.Sprintf("%d", id),
			},
		})
		if err != nil {
			logs.Error(err)
			continue
		}

		data, err := this.FormatAccessLogBytes(accessLog)
		if err != nil {
			logs.Error(err)
			continue
		}

		if this.Format != StorageFormatJSON {
			m := map[string]interface{}{
				"log": teautils.UnsafeBytesToString(data),
			}
			mData, err := ffjson.Marshal(m)
			if err != nil {
				logs.Error(err)
				continue
			}

			bulk.Write(opData)
			bulk.WriteString("\n")
			bulk.Write(mData)
			bulk.WriteString("\n")
		} else {
			bulk.Write(opData)
			bulk.WriteString("\n")
			bulk.Write(data)
			bulk.WriteString("\n")
		}
	}

	if bulk.Len() == 0 {
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, this.Endpoint+"/_bulk", strings.NewReader(bulk.String()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", teaconst.TeaProductName+"/"+teaconst.TeaVersion)
	if len(this.Username) > 0 || len(this.Password) > 0 {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(this.Username+":"+this.Password)))
	}
	client := teautils.SharedHttpClient(10 * time.Second)
	defer func() {
		_ = req.Body.Close()
	}()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		bodyData, _ := ioutil.ReadAll(resp.Body)
		return errors.New("ElasticSearch response status code: " + fmt.Sprintf("%d", resp.StatusCode) + " content: " + string(bodyData))
	}

	return nil
}

// 关闭
func (this *ESStorage) Close() error {
	return nil
}
