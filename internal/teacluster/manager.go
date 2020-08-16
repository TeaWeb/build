package teacluster

import (
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teacluster/configs"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/vmihailenco/msgpack"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

var SharedManager = NewManager()

// cluster communication manager
type Manager struct {
	RestartChan chan bool

	conn    net.Conn
	encoder *msgpack.Encoder
	decoder *msgpack.Decoder

	isStarting  bool
	startLocker sync.Mutex

	prevAction  ActionInterface
	queueLocker sync.Mutex

	error string

	isActive  bool
	isChanged bool
}

func NewManager() *Manager {
	return &Manager{
		RestartChan: make(chan bool, 10),
	}
}

// start manager
func (this *Manager) Start() error {
	this.startLocker.Lock()
	teaconst.AgentEnabled = true

	if this.isStarting {
		this.startLocker.Unlock()
		return nil
	}
	this.isStarting = true
	defer func() {
		this.isStarting = false
		this.startLocker.Unlock()
	}()

	node := teaconfigs.SharedNodeConfig()
	if node == nil {
		return nil
	}

	teaconst.AgentEnabled = !node.On
	if !node.On {
		return nil
	}

	if len(node.ClusterAddr) == 0 {
		return errors.New("'clusterAddr' should not be empty")
	}

	conn, err := net.DialTimeout("tcp", node.ClusterAddr, 10*time.Second)
	if err != nil {
		this.error = err.Error()
		return err
	}
	defer func() {
		// close connection
		_ = conn.Close()
	}()

	this.isActive = true

	this.conn = conn
	this.encoder = msgpack.NewEncoder(this.conn)
	this.decoder = msgpack.NewDecoder(this.conn)

	// register
	err = this.Write(&RegisterAction{
		ClusterId:     node.ClusterId,
		ClusterSecret: node.ClusterSecret,
		NodeId:        node.Id,
		NodeName:      node.Name,
		NodeRole:      node.Role,
	})
	if err != nil {
		logs.Error(errors.New("fail to register node"))
	}

	this.Read(func(action ActionInterface) {
		logs.Println("[cluster]receive action", "'"+action.Name()+"'")
		if action.Name() == "success" || action.Name() == "fail" {
			if this.prevAction != nil && action.BaseAction().RequestId == this.prevAction.BaseAction().Id {
				switch action.Name() {
				case "success":
					this.error = ""
					err = this.prevAction.OnSuccess(action.(*SuccessAction))
					if err != nil {
						logs.Error(err)
					}
				case "fail":
					this.error = action.(*FailAction).Message
					err = this.prevAction.OnFail(action.(*FailAction))
					if err != nil {
						logs.Error(err)
					}
				}
			}
		}
		err = action.Execute()
		if err != nil {
			logs.Error(err)
		}
	})

	return nil
}

// read action from cluster
func (this *Manager) Read(f func(action ActionInterface)) {
	for {
		typeId, _, err := this.decoder.DecodeExtHeader()
		if err != nil {
			if err == io.EOF {
				break
			}
			this.error = err.Error()
			break
		}
		instance := FindActionInstance(typeId)
		if instance == nil {
			logs.Error(errors.New("can not find action type '" + fmt.Sprintf("%d", typeId) + "'"))
			continue
		}
		err = this.decoder.Decode(instance)
		if err != nil {
			if err == io.EOF {
				break
			}
			this.error = err.Error()
			break
		}
		f(instance)
	}

	this.isActive = false
}

// write action message to cluster manager
func (this *Manager) Write(action ActionInterface) error {
	if this.conn == nil {
		return errors.New("no connection to cluster")
	}
	if action.Name() != "ping" {
		logs.Println("[cluster]send action '" + action.Name() + "'")
	}
	this.queueLocker.Lock()
	this.prevAction = action
	action.BaseAction().Id = GenerateActionId()
	err := this.encoder.Encode(action)
	this.queueLocker.Unlock()
	return err
}

// stop manager
func (this *Manager) Stop() error {
	conn := this.conn
	if conn != nil {
		err := conn.Close()
		this.conn = nil
		return err
	}
	return nil
}

// is active
func (this *Manager) IsActive() bool {
	return this.isActive
}

func (this *Manager) Error() string {
	return this.error
}

func (this *Manager) Restart() {
	err := this.Stop()
	if err != nil {
		logs.Error(err)
	}
	this.RestartChan <- true
}

func (this *Manager) PushItems() {
	action := &PushAction{}

	// cluster sum
	clusterSumMap := this.clusterSumMap()

	// proxy
	RangeFiles(func(file *files.File, relativePath string) {
		item := configs.NewItem()
		item.Id = relativePath
		sum, err := file.Md5()
		if err != nil {
			logs.Error(err)
			return
		}
		item.Sum = sum

		clusterSum, ok := clusterSumMap[item.Id]
		if !ok || clusterSum != sum {
			data, err := file.ReadAll()
			if err != nil {
				logs.Error(err)
				return
			}
			item.Data = data
			logs.Println("[cluster]push '"+relativePath+"'", len(data), "bytes")
		}

		stat, err := file.Stat()
		if err == nil {
			item.Flags = []int{int(stat.Mode.Perm())}
		}
		action.AddItem(item)
	})

	if len(action.Items) == 0 {
		return
	}

	err := SharedManager.Write(action)
	if err != nil {
		logs.Error(err)
	} else {
		this.isChanged = false
	}
}

func (this *Manager) PullItems() {
	action := &PullAction{}

	for id, sum := range this.nodeSumMap() {
		action.LocalItems = append(action.LocalItems, &configs.Item{
			Id:  id,
			Sum: sum,
		})
	}
	err := SharedManager.Write(action)
	if err != nil {
		logs.Error(err)
	} else {
		this.isChanged = false
	}
}

func (this *Manager) BuildSum() []byte {
	sumList := []string{}
	RangeFiles(func(file *files.File, relativePath string) {
		sum, err := file.Md5()
		if err != nil {
			logs.Error(err)
			return
		}
		sumList = append(sumList, relativePath+"|"+sum)
	})

	sumData := []byte(strings.Join(sumList, "\n"))
	sumFile := files.NewFile(Tea.ConfigFile("node.sum"))
	err := sumFile.Write(sumData)
	if err != nil {
		logs.Error(err)
	}
	return sumData
}

// determine cluster data changes
func (this *Manager) IsChanged() bool {
	if this.isChanged {
		return true
	}

	node := teaconfigs.SharedNodeConfig()
	if node == nil || !node.On {
		return false
	}

	map1 := this.clusterSumMap()
	map2 := this.nodeSumMap()

	for k1, v1 := range map1 {
		v2, ok := map2[k1]
		if !ok {
			return true
		}
		if v1 != v2 {
			return true
		}
	}

	for k2, v2 := range map2 {
		v1, ok := map1[k2]
		if !ok {
			return true
		}
		if v1 != v2 {
			return true
		}
	}
	return false
}

func (this *Manager) SetIsChanged(isChanged bool) {
	this.isChanged = isChanged
}

func (this *Manager) clusterSumMap() map[string]string {
	return this.sumMap(Tea.ConfigFile("cluster.sum"))
}

func (this *Manager) nodeSumMap() map[string]string {
	return this.sumMap(Tea.ConfigFile("node.sum"))
}

func (this *Manager) sumMap(path string) map[string]string {
	file := files.NewFile(path)
	if !file.Exists() {
		return nil
	}
	s, err := file.ReadAllString()
	if err != nil {
		return nil
	}

	if len(s) == 0 {
		return nil
	}

	result := map[string]string{}
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		pieces := strings.SplitN(line, "|", 2)
		if len(pieces) != 2 {
			continue
		}
		result[pieces[0]] = pieces[1]
	}
	return result
}
