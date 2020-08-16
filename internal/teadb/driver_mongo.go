package teadb

import (
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/processes"
	"github.com/iwind/TeaGo/types"
	stringutil "github.com/iwind/TeaGo/utils/string"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"github.com/shirou/gopsutil/mem"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"golang.org/x/net/context"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

type MongoDriver struct {
	BaseDriver

	sharedClient       *mongo.Client
	sharedClientLocker sync.Mutex
	dbName             string

	collMap    map[string]*MongoCollection
	collLocker sync.Mutex
}

func (this *MongoDriver) Init() {
	this.collMap = map[string]*MongoCollection{}

	agentValueDAO = new(MongoAgentValueDAO)
	agentValueDAO.SetDriver(this)
	agentValueDAO.Init()

	agentLogDAO = new(MongoAgentLogDAO)
	agentLogDAO.SetDriver(this)
	agentLogDAO.Init()

	serverValueDAO = new(MongoServerValueDAO)
	serverValueDAO.SetDriver(this)
	serverValueDAO.Init()

	auditLogDAO = new(MongoAuditLogDAO)
	auditLogDAO.SetDriver(this)
	auditLogDAO.Init()

	accessLogDAO = new(MongoAccessLogDAO)
	accessLogDAO.SetDriver(this)
	accessLogDAO.Init()

	noticeDAO = new(MongoNoticeDAO)
	noticeDAO.SetDriver(this)
	noticeDAO.Init()

	this.initDB()
}

func (this *MongoDriver) FindOne(query *Query, modelPtr interface{}) (interface{}, error) {
	if !this.IsAvailable() {
		return nil, ErrorDBUnavailable
	}

	if len(query.table) == 0 {
		return nil, errors.New("'table' should not be empty")
	}

	currentDB := this.DB()
	if currentDB == nil {
		return nil, errors.New("can not select db")
	}

	opt := options.Find()
	if query.offset > -1 {
		opt.SetSkip(int64(query.offset))
	}
	opt.SetLimit(1)

	if len(query.resultFields) > 0 {
		projection := map[string]interface{}{}
		for _, field := range query.resultFields {
			projection[field] = 1
		}
		opt.SetProjection(projection)
	}

	if len(query.sortFields) > 0 {
		s := map[string]int{}
		for _, f := range query.sortFields {
			if f.IsAsc() {
				s[f.Name] = 1
			} else {
				s[f.Name] = -1
			}
		}
		opt.SetSort(s)
	}

	filter, err := this.buildFilter(query)
	if err != nil {
		return nil, err
	}
	if query.debug {
		logs.Println("===filter===")
		logs.PrintAsJSON(filter)
	}

	cursor, err := currentDB.Collection(query.table).Find(this.timeoutContext(5*time.Second), filter, opt)
	if err != nil {
		if this.isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	defer func(cursor *mongo.Cursor) {
		_ = cursor.Close(context.Background())
	}(cursor)

	if !cursor.Next(context.Background()) {
		return nil, nil
	}

	err = cursor.Decode(modelPtr)
	if err != nil {
		if this.isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return modelPtr, nil
}

func (this *MongoDriver) FindOnes(query *Query, modelPtr interface{}) ([]interface{}, error) {
	if !this.IsAvailable() {
		return nil, ErrorDBUnavailable
	}

	if len(query.table) == 0 {
		return nil, errors.New("'table' should not be empty")
	}

	currentDB := this.DB()
	if currentDB == nil {
		return nil, errors.New("can not select db")
	}

	// 查询选项
	opt := options.Find()
	if query.offset > -1 {
		opt.SetSkip(int64(query.offset))
	}
	if query.size > -1 {
		opt.SetLimit(int64(query.size))
	}

	if len(query.resultFields) > 0 {
		projection := map[string]interface{}{}
		for _, field := range query.resultFields {
			projection[field] = 1
		}
		opt.SetProjection(projection)
	}

	if len(query.sortFields) > 0 {
		s := map[string]int{}
		for _, f := range query.sortFields {
			if f.IsAsc() {
				s[f.Name] = 1
			} else {
				s[f.Name] = -1
			}
		}
		opt.SetSort(s)
	}

	filter, err := this.buildFilter(query)
	if err != nil {
		return nil, err
	}
	if query.debug {
		logs.Println("===filter===")
		logs.PrintAsJSON(filter)
	}

	cursor, err := currentDB.Collection(query.table).Find(this.timeoutContext(5*time.Second), filter, opt)
	if err != nil {
		if this.isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	defer func(cursor *mongo.Cursor) {
		_ = cursor.Close(context.Background())
	}(cursor)

	modelType := reflect.TypeOf(modelPtr).Elem()
	result := []interface{}{}
	for cursor.Next(context.Background()) {
		m := reflect.New(modelType).Interface()
		err = cursor.Decode(m)
		if err != nil {
			if this.isNotFoundError(err) {
				continue
			}
			return nil, err
		}

		result = append(result, m)
	}
	return result, nil
}

func (this *MongoDriver) DeleteOnes(query *Query) error {
	if !this.IsAvailable() {
		return ErrorDBUnavailable
	}

	if len(query.table) == 0 {
		return errors.New("'table' should not be empty")
	}

	filter, err := this.buildFilter(query)
	if err != nil {
		return err
	}

	_, err = this.DB().Collection(query.table).DeleteMany(this.timeoutContext(5*time.Second), filter)
	return err
}

func (this *MongoDriver) InsertOne(table string, modelPtr interface{}) error {
	if !this.IsAvailable() {
		return ErrorDBUnavailable
	}
	if len(table) == 0 {
		return errors.New("'table' should not be empty")
	}
	if modelPtr == nil {
		return errors.New("insertOne: modelPtr should not be nil")
	}
	_, err := this.DB().Collection(table).InsertOne(this.timeoutContext(5*time.Second), modelPtr)
	return err
}

func (this *MongoDriver) InsertOnes(table string, modelPtrSlice interface{}) error {
	if !this.IsAvailable() {
		return ErrorDBUnavailable
	}
	if len(table) == 0 {
		return errors.New("'table' should not be empty")
	}
	if modelPtrSlice == nil {
		return nil
	}

	t := reflect.ValueOf(modelPtrSlice)
	if t.IsNil() {
		return nil
	}
	if t.Kind() != reflect.Slice {
		return errors.New("insertOnes: only slice is accepted")
	}

	s := []interface{}{}
	l := t.Len()
	for i := 0; i < l; i++ {
		s = append(s, t.Index(i).Interface())
	}

	_, err := this.DB().Collection(table).InsertMany(this.timeoutContext(10*time.Second), s)
	return err
}

func (this *MongoDriver) Count(query *Query) (int64, error) {
	if !this.IsAvailable() {
		return 0, ErrorDBUnavailable
	}
	if len(query.table) == 0 {
		return 0, errors.New("'table' should not be empty")
	}

	currentDB := this.DB()
	if currentDB == nil {
		return 0, errors.New("can not select db")
	}

	// 查询选项
	opts := options.Count()
	if query.offset > -1 {
		opts.SetSkip(int64(query.offset))
	}
	if query.size > -1 {
		opts.SetLimit(int64(query.size))
	}

	filter, err := this.buildFilter(query)
	if err != nil {
		return 0, err
	}

	return this.DB().Collection(query.table).CountDocuments(this.timeoutContext(10*time.Second), filter, opts)
}

func (this *MongoDriver) Sum(query *Query, field string) (float64, error) {
	return this.aggregate("sum", query, field)
}

func (this *MongoDriver) Avg(query *Query, field string) (float64, error) {
	return this.aggregate("avg", query, field)
}

func (this *MongoDriver) Min(query *Query, field string) (float64, error) {
	return this.aggregate("min", query, field)
}

func (this *MongoDriver) Max(query *Query, field string) (float64, error) {
	return this.aggregate("max", query, field)
}

func (this *MongoDriver) Group(query *Query, field string, result map[string]Expr) ([]maps.Map, error) {
	if !this.IsAvailable() {
		return nil, ErrorDBUnavailable
	}

	group := map[string]interface{}{
		"_id": "$" + field,
	}

	for name, expr := range result {
		// 处理点符号
		name = strings.Replace(name, ".", "@", -1)

		switch e := expr.(type) {
		case *SumExpr:
			group[name] = map[string]interface{}{
				"$sum": this.convertArrayElement(e.Field),
			}
		case *AvgExpr:
			group[name] = map[string]interface{}{
				"$avg": this.convertArrayElement(e.Field),
			}
		case *MaxExpr:
			group[name] = map[string]interface{}{
				"$max": this.convertArrayElement(e.Field),
			}
		case *MinExpr:
			group[name] = map[string]interface{}{
				"$min": this.convertArrayElement(e.Field),
			}
		case string:
			group[name] = map[string]interface{}{
				"$first": this.convertArrayElement(e),
			}
		case maps.Map:
			group[name] = e
		}
	}

	sorts := map[string]interface{}{}
	if len(query.sortFields) > 0 {
		for _, sortField := range query.sortFields {
			if sortField.IsAsc() {
				sorts[sortField.Name] = 1
			} else {
				sorts[sortField.Name] = -1
			}
		}
	}

	filter, err := this.buildFilter(query)
	if err != nil {
		return nil, err
	}
	pipelines := []interface{}{
		map[string]interface{}{
			"$match": filter,
		},
		map[string]interface{}{
			"$limit": 100000, // 限制进入下一个pipeline的记录数量，以避免查询超时
		},
		map[string]interface{}{
			"$group": group,
		},
	}
	if len(sorts) > 0 {
		pipelines = append(pipelines, map[string]interface{}{
			"$sort": sorts,
		})
	}

	if query.debug {
		logs.Println("===pipelines===")
		logs.PrintAsJSON(pipelines)
	}

	cursor, err := this.DB().Collection(query.table).Aggregate(this.timeoutContext(30*time.Second), pipelines)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = cursor.Close(context.Background())
	}()

	ones := []maps.Map{}

	for cursor.Next(context.Background()) {
		m := maps.Map{}
		err = cursor.Decode(&m)
		if err != nil {
			return nil, err
		}

		// 处理@符号（从上面的点符号转换过来）
		for k, v := range m {
			if strings.Contains(k, "@") {
				this.setMapValue(m, strings.Split(k, "@"), v)
				delete(m, k)
			}
		}

		ones = append(ones, m)
	}

	return ones, nil
}

// 测试数据库连接
func (this *MongoDriver) Test() error {
	client, err := this.connect()
	if err != nil {
		return err
	}

	// 尝试查询
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err = client.Database(this.dbName).
		Collection("logs").
		Find(ctx, map[string]interface{}{}, options.Find().SetLimit(1))

	// 重置客户端
	if err != nil {
		this.sharedClientLocker.Lock()
		this.sharedClient = nil
		this.sharedClientLocker.Unlock()
	}

	return err
}

// 删除表
func (this *MongoDriver) DropTable(table string) error {
	coll, err := this.SelectColl(table)
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return coll.Drop(ctx)
}

// 测试URI
func (this *MongoDriver) TestConfig(config *db.MongoConfig) (message string, ok bool) {
	opts := options.Client().ApplyURI(config.ComposeURI())
	if config.AuthEnabled {
		opts.SetAuth(options.Credential{
			Username:                config.Username,
			Password:                config.Password,
			AuthMechanism:           config.AuthMechanism,
			AuthMechanismProperties: config.AuthMechanismPropertiesMap(),
			AuthSource:              config.DBName,
		})
	}
	client, err := mongo.NewClient(opts)
	if err != nil {
		message = "尝试分析配置错误：" + err.Error()
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		message = "尝试连接数据库失败：" + err.Error()
		return
	}

	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	coll := client.Database(config.DBName).Collection("test")
	_, err = coll.InsertOne(context.Background(), map[string]interface{}{
		"a": 1,
	})
	if err != nil {
		message = "尝试写入数据失败：" + err.Error()
		return
	}

	err = coll.Drop(context.Background())
	if err != nil {
		message = "尝试删除集合失败：" + err.Error()
		return
	}

	ok = true
	return
}

// 关闭服务
func (this *MongoDriver) Shutdown() error {
	this.collLocker.Lock()
	this.collMap = map[string]*MongoCollection{}
	this.collLocker.Unlock()

	this.sharedClientLocker.Lock()
	if this.sharedClient != nil {
		oldClient := this.sharedClient
		go func() {
			if oldClient != nil {
				_ = oldClient.Disconnect(context.Background())
			}
		}()
		this.sharedClient = nil
	}
	this.sharedClientLocker.Unlock()

	return nil
}

// 选择数据集合
func (this *MongoDriver) SelectColl(name string) (*MongoCollection, error) {
	this.collLocker.Lock()
	defer this.collLocker.Unlock()

	coll, ok := this.collMap[name]
	if ok {
		return coll, nil
	}

	currentDB := this.DB()
	if currentDB != nil {
		coll = &MongoCollection{
			currentDB.Collection(name),
		}
		this.collMap[name] = coll
		return coll, nil
	}
	return nil, errors.New("can not select collection '" + name + "'")
}

// 列出所有表
func (this *MongoDriver) ListTables() ([]string, error) {
	if !this.isAvailable {
		return nil, ErrorDBUnavailable
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := this.DB().ListCollections(ctx, maps.Map{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(context.Background())
	}()

	names := []string{}
	for cursor.Next(context.Background()) {
		m := maps.Map{}
		err := cursor.Decode(&m)
		if err != nil {
			return nil, err
		}
		name := m.GetString("name")
		if len(name) == 0 {
			continue
		}
		names = append(names, name)
	}

	sort.Strings(names)

	return names, nil
}

// 统计数据表
func (this *MongoDriver) StatTables(tables []string) (map[string]*TableStat, error) {
	statMap := map[string]*TableStat{}

	currentDB := this.DB()
	if currentDB == nil {
		return statMap, errors.New("database is not available now")
	}

	for _, collName := range tables {
		if len(collName) == 0 {
			continue
		}

		result := currentDB.RunCommand(context.Background(), bsonx.Doc{{"collStats", bsonx.String(collName)}, {"verbose", bsonx.Boolean(false)}})
		if result.Err() != nil {
			return statMap, result.Err()
		}

		m1 := maps.Map{}
		err := result.Decode(&m1)
		if err != nil {
			if result.Err() != nil {
				return statMap, result.Err()
			}
		}
		if m1.GetInt("ok") != 1 {
			continue
		}
		size := float64(m1.GetInt("size"))
		formattedSize := ""
		if size < 1024 {
			formattedSize = fmt.Sprintf("%.2fB", size)
		} else if size < 1024*1024 {
			formattedSize = fmt.Sprintf("%.2fKB", size/1024)
		} else if size < 1024*1024*1024 {
			formattedSize = fmt.Sprintf("%.2fMB", size/1024/1024)
		} else {
			formattedSize = fmt.Sprintf("%.2fGB", size/1024/1024/1024)
		}
		statMap[collName] = &TableStat{
			Count:         m1.GetInt64("count"),
			Size:          m1.GetInt64("size"),
			FormattedSize: formattedSize,
		}
	}
	return statMap, nil
}

// 选取数据库
func (this *MongoDriver) DB() *mongo.Database {
	client, _ := this.connect()
	if client == nil {
		return nil
	}
	return client.Database(this.dbName)
}

// 获取共享的Client
func (this *MongoDriver) connect() (*mongo.Client, error) {
	if this.sharedClient != nil {
		return this.sharedClient, nil
	}

	this.sharedClientLocker.Lock()
	defer this.sharedClientLocker.Unlock()

	if this.sharedClient != nil {
		return this.sharedClient, nil
	}

	config, err := db.LoadMongoConfig()
	if err != nil {
		return nil, err
	}

	if len(config.DBName) == 0 {
		config.DBName = "teaweb"
	}
	this.dbName = config.DBName

	opts := options.Client().ApplyURI(config.URI)
	if config.PoolSize > 0 {
		opts.SetMaxPoolSize(uint64(config.PoolSize))
	} else {
		opts.SetMaxPoolSize(32)
	}
	if config.Timeout > 0 {
		opts.SetConnectTimeout(time.Duration(5) * time.Second)
	} else {
		opts.SetConnectTimeout(5 * time.Second)
	}
	sharedConfig, err := db.LoadMongoConfig()
	if err != nil {
		return nil, err
	}

	if sharedConfig != nil && sharedConfig.AuthEnabled && len(sharedConfig.AuthMechanism) > 0 {
		opts.SetAuth(options.Credential{
			Username:                sharedConfig.Username,
			Password:                sharedConfig.Password,
			AuthMechanism:           sharedConfig.AuthMechanism,
			AuthMechanismProperties: sharedConfig.AuthMechanismPropertiesMap(),
			AuthSource:              config.DBName,
		})
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	this.sharedClient = client

	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil

}

func (this *MongoDriver) buildFilter(query *Query) (filter map[string]interface{}, err error) {
	if query.operandList.Len() > 0 {
		return this.buildOperandMap(query.operandList)
	}
	return map[string]interface{}{}, nil
}

func (this *MongoDriver) buildOperandMap(operandList *OperandList) (filter map[string]interface{}, err error) {
	filter = map[string]interface{}{}
	operandList.Range(func(field string, operands []*Operand) {
		fieldQuery := map[string]interface{}{}
		for _, op := range operands {
			switch op.Code {
			case OperandEq:
				fieldQuery["$eq"] = op.Value
			case OperandLt:
				fieldQuery["$lt"] = op.Value
			case OperandLte:
				fieldQuery["$lte"] = op.Value
			case OperandGt:
				fieldQuery["$gt"] = op.Value
			case OperandGte:
				fieldQuery["$gte"] = op.Value
			case OperandIn:
				fieldQuery["$in"] = op.Value
			case OperandNotIn:
				fieldQuery["$nin"] = op.Value
			case OperandNeq:
				fieldQuery["$ne"] = op.Value
			case OperandOr:
				if op.Value != nil {
					operandLists, ok := op.Value.([]*OperandList)
					if ok {
						result := []map[string]interface{}{}
						for _, operandList := range operandLists {
							f, err := this.buildOperandMap(operandList)
							if err != nil {
								return
							}
							result = append(result, f)
						}
						filter["$or"] = result
					} else {
						err = errors.New("or: should be a valid []OperandMap")
						return
					}
				} else {
					err = errors.New("or: should be a valid []OperandMap")
					return
				}
			}
		}
		if len(fieldQuery) > 0 {
			filter[field] = fieldQuery
		}
	})

	return
}

func (this *MongoDriver) isNotFoundError(err error) bool {
	return err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments
}

func (this *MongoDriver) timeoutContext(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}

func (this *MongoDriver) aggregate(funcName string, query *Query, field string) (float64, error) {
	if !this.IsAvailable() {
		return 0, ErrorDBUnavailable
	}

	filter, err := this.buildFilter(query)
	if err != nil {
		return 0, err
	}

	pipelines, err := BSONArrayBytes([]byte(`[
	{
		"$match": ` + stringutil.JSONEncode(filter) + `
	},
	{
		"$group": {
			"_id": null,
			"result": { "$` + funcName + `": ` + stringutil.JSONEncode("$"+field) + `}
		}
	}
]`))
	if err != nil {
		return 0, err
	}

	cursor, err := this.DB().Collection(query.table).Aggregate(this.timeoutContext(30*time.Second), pipelines)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = cursor.Close(context.Background())
	}()

	m := maps.Map{}
	if !cursor.Next(context.Background()) {
		return 0, nil
	}
	err = cursor.Decode(&m)
	if err != nil {
		return 0, err
	}

	return m.GetFloat64("result"), nil
}

func (this *MongoDriver) setMapValue(m maps.Map, keys []string, v interface{}) {
	l := len(keys)
	if l == 0 {
		return
	}
	if l == 1 {
		m[keys[0]] = v
		return
	}
	subM, ok := m[keys[0]]
	if ok {
		subV, ok := subM.(maps.Map)
		if ok {
			this.setMapValue(subV, keys[1:], v)
		} else {
			m[keys[0]] = maps.Map{}
			this.setMapValue(m[keys[0]].(maps.Map), keys[1:], v)
		}
	} else {
		m[keys[0]] = maps.Map{}
		this.setMapValue(m[keys[0]].(maps.Map), keys[1:], v)
	}
}

// 初始化MongoDB
func (this *MongoDriver) initDB() {
	go func() {
		this.startInstalledMongo()
		this.cleanAccessLogs()
	}()
}

// 启动本机安装的Mongo
func (this *MongoDriver) startInstalledMongo() {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		return
	}

	config, err := db.LoadMongoConfig()
	if err != nil {
		logs.Println("load mongodb configure failed: " + err.Error())
		return
	}
	if strings.HasSuffix(config.Addr, "127.0.0.1") || strings.HasSuffix(config.Addr, "localhost") {
		return
	}

	err = this.Test()

	if err != nil {
		mongodbDir := Tea.Root + "/mongodb"

		// 是否已安装
		if !files.NewFile(mongodbDir + "/bin/mongod").Exists() {
			return
		}

		// 启动
		args := []string{"--dbpath=" + mongodbDir + "/data", "--fork", "--logpath=" + mongodbDir + "/data/fork.log"}

		// 控制内存不能超过1G
		stat, err := mem.VirtualMemory()
		if err == nil && stat.Total > 0 {
			count := stat.Total / 1024 / 1024 / 1024
			if count >= 6 {
				args = append(args, "--wiredTigerCacheSizeGB=2")
			} else if count >= 3 {
				args = append(args, "--wiredTigerCacheSizeGB=1")
			}
		}

		p := processes.NewProcess(mongodbDir+"/bin/mongod", args...)
		p.SetPwd(mongodbDir)

		logs.Println("start mongo:", mongodbDir+"/bin/mongod", strings.Join(args, " "))

		err = p.StartBackground()
		if err != nil {
			logs.Println("[mongo]start error: " + err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

// 清理访问日志任务
func (this *MongoDriver) cleanAccessLogs() {
	reg := regexp.MustCompile("^logs\\.\\d{8}$")

	teautils.Every(1*time.Minute, func(ticker *teautils.Ticker) {
		config, _ := db.LoadMongoConfig()
		if config == nil {
			return
		}

		if config.AccessLog == nil {
			return
		}

		now := time.Now()
		if config.AccessLog.CleanHour != now.Hour() ||
			now.Minute() != 0 ||
			config.AccessLog.KeepDays < 1 {
			return
		}

		compareDay := "logs." + timeutil.Format("Ymd", time.Now().Add(-time.Duration(config.AccessLog.KeepDays*24)*time.Hour))
		logs.Println("[mongo]clean access logs before '" + compareDay + "'")

		currentDB := this.DB()
		if currentDB == nil {
			return
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		cursor, err := currentDB.ListCollections(ctx, maps.Map{})
		if err != nil {
			logs.Error(err)
			return
		}

		defer func() {
			err = cursor.Close(context.Background())
			if err != nil {
				logs.Error(err)
			}
		}()

		for cursor.Next(context.Background()) {
			m := maps.Map{}
			err := cursor.Decode(&m)
			if err != nil {
				logs.Error(err)
				return
			}
			name := m.GetString("name")
			if len(name) == 0 {
				continue
			}
			if !reg.MatchString(name) {
				continue
			}

			if name < compareDay {
				logs.Println("[mongo]clean collection '" + name + "'")
				ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
				err := currentDB.Collection(name).Drop(ctx)
				if err != nil {
					logs.Error(err)
				}
			}
		}
	})
}

// 转换字段中的数字，方便聚合使用
func (this *MongoDriver) convertArrayElement(field string) (result interface{}) {
	if len(field) == 0 {
		return ""
	}
	pieces := strings.Split(field, ".")
	for index, piece := range pieces {
		if teautils.RegexpDigitNumber.MatchString(piece) {
			return maps.Map{
				"$convert": maps.Map{
					"input": maps.Map{
						"$arrayElemAt": []interface{}{
							"$" + strings.Join(pieces[:index], "."),
							types.Int(piece),
						},
					},
					"to": "double",
				},
			}
		}
	}
	return "$" + field
}
