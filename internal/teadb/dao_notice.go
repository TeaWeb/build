package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"time"
)

// 通知DAO
type NoticeDAOInterface interface {
	// 设置驱动
	SetDriver(driver DriverInterface)

	// 初始化
	Init()

	// 写入一个通知
	InsertOne(notice *notices.Notice) error

	// 发送一个代理的通知（形式1）
	NotifyProxyMessage(cond notices.ProxyCond, message string) error

	// 发送一个代理的通知（形式2）
	NotifyProxyServerMessage(serverId string, level notices.NoticeLevel, message string) error

	// 获取所有未读通知数
	CountAllUnreadNotices() (int, error)

	// 获取所有已读通知数
	CountAllReadNotices() (int, error)

	// 获取某个Agent的未读通知数
	CountUnreadNoticesForAgent(agentId string) (int, error)

	// 获取某个Agent已读通知数
	CountReadNoticesForAgent(agentId string) (int, error)

	// 获取某个接收人在某个时间段内接收的通知数
	CountReceivedNotices(receiverId string, cond map[string]interface{}, minutes int) (int, error)

	// 通过Hash判断是否存在相同的消息
	ExistNoticesWithHash(hash string, cond map[string]interface{}, duration time.Duration) (bool, error)

	// 列出消息
	ListNotices(isRead bool, offset int, size int) ([]*notices.Notice, error)

	// 列出某个Agent相关的消息
	ListAgentNotices(agentId string, isRead bool, offset int, size int) ([]*notices.Notice, error)

	// 删除Agent相关通知
	DeleteNoticesForAgent(agentId string) error

	// 更改某个通知的接收人
	UpdateNoticeReceivers(noticeId string, receiverIds []string) error

	// 设置全部已读
	UpdateAllNoticesRead() error

	// 设置一组通知已读
	UpdateNoticesRead(noticeIds []string) error

	// 设置Agent的一组通知已读
	UpdateAgentNoticesRead(agentId string, noticeIds []string) error

	// 设置Agent所有通知已读
	UpdateAllAgentNoticesRead(agentId string) error
}
