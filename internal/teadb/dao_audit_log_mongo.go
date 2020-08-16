package teadb

import "github.com/TeaWeb/build/internal/teaconfigs/audits"

type MongoAuditLogDAO struct {
	BaseDAO
}

func (this *MongoAuditLogDAO) Init() {

}

func (this *MongoAuditLogDAO) CountAllAuditLogs() (int64, error) {
	return NewQuery(this.collName()).Count()
}

func (this *MongoAuditLogDAO) ListAuditLogs(offset int, size int) ([]*audits.Log, error) {
	ones, err := NewQuery(this.collName()).
		Offset(offset).
		Limit(size).
		Desc("_id").
		FindOnes(new(audits.Log))
	if err != nil {
		return nil, err
	}

	result := []*audits.Log{}
	for _, one := range ones {
		result = append(result, one.(*audits.Log))
	}

	return result, nil
}

func (this *MongoAuditLogDAO) InsertOne(auditLog *audits.Log) error {
	return NewQuery(this.collName()).InsertOne(auditLog)
}

func (this *MongoAuditLogDAO) collName() string {
	return "logs.audit"
}
