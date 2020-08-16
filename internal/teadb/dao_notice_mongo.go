package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"golang.org/x/net/context"
	"time"
)

type MongoNoticeDAO struct {
	BaseDAO
}

func (this *MongoNoticeDAO) Init() {
	go this.initIndexes()
}

func (this *MongoNoticeDAO) TableName() string {
	return "notices"
}

func (this *MongoNoticeDAO) InsertOne(notice *notices.Notice) error {
	return NewQuery(this.TableName()).
		InsertOne(notice)
}

func (this *MongoNoticeDAO) NotifyProxyMessage(cond notices.ProxyCond, message string) error {
	notice := notices.NewNotice()
	notice.Message = message
	notice.SetTime(time.Now())
	notice.Proxy = cond
	notice.Hash()
	return NewQuery(this.TableName()).InsertOne(notice)
}

func (this *MongoNoticeDAO) NotifyProxyServerMessage(serverId string, level notices.NoticeLevel, message string) error {
	return this.NotifyProxyMessage(notices.ProxyCond{
		ServerId: serverId,
		Level:    level,
	}, message)
}

func (this *MongoNoticeDAO) CountAllUnreadNotices() (int, error) {
	count, err := NewQuery(this.TableName()).
		Attr("isRead", false).
		Count()
	return int(count), err
}

func (this *MongoNoticeDAO) CountAllReadNotices() (int, error) {
	count, err := NewQuery(this.TableName()).
		Attr("isRead", true).
		Count()
	return int(count), err
}

func (this *MongoNoticeDAO) CountUnreadNoticesForAgent(agentId string) (int, error) {
	count, err := NewQuery(this.TableName()).
		Attr("agent.agentId", agentId).
		Attr("isRead", false).
		Count()
	return int(count), err
}

func (this *MongoNoticeDAO) CountReadNoticesForAgent(agentId string) (int, error) {
	count, err := NewQuery(this.TableName()).
		Attr("agent.agentId", agentId).
		Attr("isRead", true).
		Count()
	return int(count), err
}

func (this *MongoNoticeDAO) CountReceivedNotices(receiverId string, cond map[string]interface{}, minutes int) (int, error) {
	if len(receiverId) == 0 {
		return 0, nil
	}
	if minutes <= 0 {
		return 0, nil
	}
	query := NewQuery(this.TableName()).
		Attr("receivers", receiverId).
		Gte("timestamp", time.Now().Unix()-int64(minutes*60))

	if len(cond) > 0 {
		for k, v := range cond {
			query.Attr(k, v)
		}
	}
	c, err := query.Count()
	return int(c), err
}

func (this *MongoNoticeDAO) ExistNoticesWithHash(hash string, cond map[string]interface{}, duration time.Duration) (bool, error) {
	query := NewQuery(this.TableName())
	query.Attr("messageHash", hash)
	for k, v := range cond {
		query.Attr(k, v)
	}
	query.Gt("timestamp", float64(time.Now().Unix())-duration.Seconds())
	query.Desc("_id")
	one, err := query.FindOne(new(notices.Notice))
	if err != nil {
		return false, err
	}
	if one == nil {
		return false, nil
	}
	notice := one.(*notices.Notice)

	// 中间是否有success级别的
	query2 := NewQuery(this.TableName())
	for k, v := range cond {
		query2.Attr(k, v)
	}
	if len(notice.Proxy.ServerId) > 0 {
		query2.Attr("proxy.level", notices.NoticeLevelSuccess)
		query2.Gt("_id", notice.Id)
	} else if len(notice.Agent.AgentId) > 0 {
		query2.Attr("agent.level", notices.NoticeLevelSuccess)
		query2.Gt("_id", notice.Id)
	}
	result, err := query2.Result("_id").
		FindOne(new(notices.Notice))
	return result == nil, err
}

// 列出消息
func (this *MongoNoticeDAO) ListNotices(isRead bool, offset int, size int) ([]*notices.Notice, error) {
	ones, err := NewQuery(this.TableName()).
		Attr("isRead", isRead).
		Offset(offset).
		Limit(size).
		Desc("_id").
		FindOnes(new(notices.Notice))
	if err != nil {
		return nil, err
	}

	result := []*notices.Notice{}
	for _, one := range ones {
		result = append(result, one.(*notices.Notice))
	}
	return result, err
}

// 列出某个Agent相关的消息
func (this *MongoNoticeDAO) ListAgentNotices(agentId string, isRead bool, offset int, size int) ([]*notices.Notice, error) {
	ones, err := NewQuery(this.TableName()).
		Attr("agent.agentId", agentId).
		Attr("isRead", isRead).
		Offset(offset).
		Limit(size).
		Desc("_id").
		FindOnes(new(notices.Notice))
	if err != nil {
		return nil, err
	}

	result := []*notices.Notice{}
	for _, one := range ones {
		result = append(result, one.(*notices.Notice))
	}
	return result, err
}

func (this *MongoNoticeDAO) DeleteNoticesForAgent(agentId string) error {
	return NewQuery(this.TableName()).
		Attr("agent.agentId", agentId).
		Delete()
}

func (this *MongoNoticeDAO) UpdateNoticeReceivers(noticeId string, receiverIds []string) error {
	idObject, err := shared.ObjectIdFromHex(noticeId)
	if err != nil {
		return err
	}

	coll, err := this.driver.(*MongoDriver).SelectColl(this.TableName())
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = coll.UpdateMany(ctx, map[string]interface{}{
		"_id": idObject,
	}, maps.Map{
		"$set": maps.Map{
			"isNotified": true,
			"receivers":  receiverIds,
		},
	})
	return err
}

func (this *MongoNoticeDAO) UpdateAllNoticesRead() error {
	coll, err := this.driver.(*MongoDriver).SelectColl(this.TableName())
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = coll.UpdateMany(ctx, map[string]interface{}{}, maps.Map{
		"$set": maps.Map{
			"isRead": true,
		},
	})
	return err
}

// 设置一组通知已读
func (this *MongoNoticeDAO) UpdateNoticesRead(noticeIds []string) error {
	if len(noticeIds) == 0 {
		return nil
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	noticeIdObjects := []shared.ObjectId{}
	for _, noticeId := range noticeIds {
		idObject, err := shared.ObjectIdFromHex(noticeId)
		if err != nil {
			return err
		}
		noticeIdObjects = append(noticeIdObjects, idObject)
	}
	coll, err := this.driver.(*MongoDriver).SelectColl(this.TableName())
	if err != nil {
		return err
	}
	_, err = coll.UpdateMany(ctx, map[string]interface{}{
		"_id": map[string]interface{}{
			"$in": noticeIdObjects,
		},
	}, maps.Map{
		"$set": maps.Map{
			"isRead": true,
		},
	})
	return err
}

// 设置Agent的一组通知已读
func (this *MongoNoticeDAO) UpdateAgentNoticesRead(agentId string, noticeIds []string) error {
	if len(noticeIds) == 0 {
		return nil
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	noticeIdObjects := []shared.ObjectId{}
	for _, noticeId := range noticeIds {
		idObject, err := shared.ObjectIdFromHex(noticeId)
		if err != nil {
			return err
		}
		noticeIdObjects = append(noticeIdObjects, idObject)
	}
	coll, err := this.driver.(*MongoDriver).SelectColl(this.TableName())
	if err != nil {
		return err
	}
	_, err = coll.UpdateMany(ctx, map[string]interface{}{
		"agent.agentId": agentId,
		"_id": map[string]interface{}{
			"$in": noticeIdObjects,
		},
	}, maps.Map{
		"$set": maps.Map{
			"isRead": true,
		},
	})
	return err
}

// 设置Agent所有通知已读
func (this *MongoNoticeDAO) UpdateAllAgentNoticesRead(agentId string) error {
	coll, err := this.driver.(*MongoDriver).SelectColl(this.TableName())
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = coll.UpdateMany(ctx, map[string]interface{}{
		"agent.agentId": agentId,
	}, maps.Map{
		"$set": maps.Map{
			"isRead": true,
		},
	})
	return err
}

// 初始化索引
func (this *MongoNoticeDAO) initIndexes() {
	if isInitializedTable(this.TableName()) {
		return
	}
	coll, err := this.driver.(*MongoDriver).SelectColl(this.TableName())
	if err != nil {
		logs.Error(err)
		return
	}
	_ = coll.CreateIndex(shared.NewIndexField("proxy.serverId", true))
	_ = coll.CreateIndex(shared.NewIndexField("agent.agentId", true))
	_ = coll.CreateIndex(
		shared.NewIndexField("agent.agentId", true),
		shared.NewIndexField("agent.appId", true),
		shared.NewIndexField("agent.itemId", true),
	)
}
