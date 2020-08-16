package tealogs

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teautils/logbuffer"
	"github.com/iwind/TeaGo/logs"
	"github.com/mailru/easyjson"
	"time"
)

// 导入数据库的访问日志的步长
const (
	DBBatchLogSize = 256
)

// 访问日志队列
type AccessLogQueue struct {
	index  int
	buffer *logbuffer.Buffer
}

// 创建队列对象
func NewAccessLogQueue(buffer *logbuffer.Buffer, index int) *AccessLogQueue {
	return &AccessLogQueue{
		buffer: buffer,
		index:  index,
	}
}

// 从队列中接收日志
func (this *AccessLogQueue) Receive(ch chan *accesslogs.AccessLog) {
	for log := range ch {
		if log == nil {
			continue
		}
		if log.ShouldStat() || log.ShouldWrite() {
			log.Parse()

			// 统计
			if log.ShouldStat() {
				CallAccessLogHooks(log)
			}

			// 保存到文件
			if log.ShouldWrite() {
				log.CleanFields()
				data, err := easyjson.Marshal(log)
				if err != nil {
					logs.Error(err)
					continue
				}
				_, err = this.buffer.Write(data)
				if err != nil {
					logs.Error(err)
				}
			}
		}
	}
}

// 导出日志到别的媒介
func (this *AccessLogQueue) Dump() {
	ticker := teautils.NewTicker(1 * time.Second)
	for range ticker.C {
		this.dumpInterval()
	}
}

// 导出日志定时内容
func (this *AccessLogQueue) dumpInterval() {
	accessLogsList := []interface{}{}

	storageLogs := map[string][]*accesslogs.AccessLog{} // policyId => accessLogs
	for i := 0; i < 4096; i++ {
		data, err := this.buffer.Read()
		if err != nil {
			logs.Error(err)
			break
		}
		if len(data) == 0 {
			break
		}
		accessLog := new(accesslogs.AccessLog)
		err = easyjson.Unmarshal(data, accessLog)
		if err != nil {
			logs.Error(err)
			continue
		}

		// 如果非storageOnly则可以存储到数据库中
		if !accessLog.StorageOnly {
			accessLog.Id = shared.NewObjectId()
			accessLogsList = append(accessLogsList, accessLog)
		}

		// 日志存储策略
		if len(accessLog.StoragePolicyIds) > 0 {
			for _, policyId := range accessLog.StoragePolicyIds {
				_, ok := storageLogs[policyId]
				if !ok {
					storageLogs[policyId] = []*accesslogs.AccessLog{}
				}
				storageLogs[policyId] = append(storageLogs[policyId], accessLog)
			}
		}
	}

	if len(storageLogs) > 0 {
		for policyId, storageAccessLogs := range storageLogs {
			storage := FindPolicyStorage(policyId)
			if storage == nil {
				continue
			}
			err := storage.Write(storageAccessLogs)
			if err != nil {
				logs.Println("access log storage policy '"+policyId+"/"+FindPolicyName(policyId)+"'", err.Error())
			}
		}
	}

	// 导入数据库
	if len(accessLogsList) > 0 {
		count := len(accessLogsList)
		offset := 0
		to := offset + DBBatchLogSize
		for {
			if to > count {
				to = count
			}
			err := teadb.AccessLogDAO().InsertAccessLogs(accessLogsList[offset:to])
			if err != nil {
				logs.Println("[logger]insert access logs:", err.Error())
			}

			offset += DBBatchLogSize
			if offset >= count {
				break
			}
			to = offset + DBBatchLogSize
		}
	}
}
