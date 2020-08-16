package mongo

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/nets"
	"github.com/iwind/TeaGo/processes"
	"github.com/shirou/gopsutil/mem"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var downloader = nets.NewDownloader()
var isInstalling = false
var installStatus = ""   // 安装状态
var installPercent = 0   // 安装百分比
var installStartAt int64 // 安装开始时间

type InstallAction actions.Action

func (this *InstallAction) Run(params struct{}) {
	this.Data["isInstalling"] = isInstalling

	err := teadb.SharedDB().Test()
	if err == nil {
		this.Data["isConnected"] = true
	} else {
		this.Data["isConnected"] = false
	}

	// 是否已安装
	mongodbDir := Tea.Root + "/mongodb"
	if files.NewFile(mongodbDir + "/bin/mongod").Exists() {
		this.Data["isInstalled"] = true
		installStatus = "start"
		installPercent = 10
	} else {
		this.Data["isInstalled"] = false
	}

	this.Show()
}

func (this *InstallAction) RunPost(params struct{}) {
	if isInstalling {
		return
	}

	defer func() {
		isInstalling = false
	}()

	isInstalling = true

	// 新建目录
	mongodbDir := Tea.Root + "/mongodb"
	mongodbDirFile := files.NewFile(mongodbDir)
	if mongodbDirFile.Exists() {
		if !mongodbDirFile.IsDir() {
			this.Fail("./mongodb应该是一个目录")
		}

		// 是否已安装
		if files.NewFile(mongodbDir + "/bin/mongod").Exists() {
			installStatus = "install"
			installPercent = 100

			this.start(mongodbDir)

			this.Success()
		}

		dataFile := files.NewFile(mongodbDir + "/data")
		if !dataFile.Exists() {
			err := dataFile.Mkdir()
			if err != nil {
				this.Fail("创建mongodb/data目录失败：", err.Error())
			}
		}
	} else {
		err := mongodbDirFile.Mkdir()
		if err != nil {
			this.Fail("创建mongodb/目录失败：", err.Error())
		}

		err = files.NewFile(mongodbDir + "/data").Mkdir()
		if err != nil {
			this.Fail("创建mongodb/data目录失败：", err.Error())
		}
	}

	resp, err := http.Get("http://dl.teaos.cn/mongodb.json")
	if err != nil {
		this.Fail("发生错误：", err.Error())
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		this.Fail("服务器状态码返回错误：", fmt.Sprintf("%d", resp.StatusCode))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		this.Fail("发生错误：", err.Error())
	}

	m := maps.Map{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		this.Fail("发生错误：", err.Error())
	}

	fileName := m.GetString(runtime.GOOS + "-" + runtime.GOARCH)
	if len(fileName) == 0 {
		this.Fail("目前不支持：", runtime.GOOS+"-"+runtime.GOARCH)
	}

	target := mongodbDir + "/" + fileName

	// 是否已下载完毕
	if files.NewFile(target).Exists() {
		installStatus = "download"
		installPercent = 100

		this.install(mongodbDir, target)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	url := "http://dl.teaos.cn/" + fileName

	downloader.Add(url, "", target+".tmp")
	downloader.OnProgress(func(item *nets.DownloaderItem) {
		installStatus = "download"
		installPercent = int(math.Ceil(float64(item.Progress()) * 100))

		// logs.Println("progress:", installPercent)
	})
	var isDownloaded = false
	downloader.OnCompleteFn(func(item *nets.DownloaderItem) {
		if item.Success() {
			installStatus = "download"
			installPercent = 100
			isDownloaded = true
		}

		wg.Add(-1)
	})
	installStartAt = time.Now().Unix()
	downloader.Start()
	wg.Wait()
	tmpFile := files.NewFile(target + ".tmp")
	if !isDownloaded || !tmpFile.Exists() {
		this.Fail("下载失败，请刷新当前页面后重试")
	}
	err = os.Rename(target+".tmp", target)
	if err != nil {
		this.Fail("移动文件失败：", err.Error())
	}

	logs.Println("download ok")
	this.install(mongodbDir, target)

	this.Success()
}

func (this *InstallAction) install(mongodbDir string, target string) {
	installStatus = "install"
	installPercent = 10

	// 是否已安装
	if files.NewFile(mongodbDir + "/bin/mongod").Exists() {
		this.start(mongodbDir)
		return
	}

	// unzip
	logs.Println("unzip target:", target)

	fp, err := os.OpenFile(target, os.O_RDONLY, 0444)
	if err != nil {
		this.Fail("压缩包读取失败：", err.Error())
	}
	defer func() {
		_ = fp.Close()
	}()

	gzReader, err := gzip.NewReader(fp)
	if err != nil {
		this.Fail("压缩包读取失败：", err.Error())
	}
	defer func() {
		_ = gzReader.Close()
	}()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			break
		}

		installPercent += 5
		if installPercent > 100 {
			installPercent = 100
		}
		info := header.FileInfo()
		filePath := mongodbDir + string(os.PathSeparator) + header.Name[strings.Index(header.Name, string(os.PathSeparator))+1:]

		if files.NewFile(filePath).Exists() {
			continue
		}

		logs.Println("install", header.Name)

		if info.IsDir() {
			err := os.MkdirAll(filePath, info.Mode())
			if err != nil {
				this.Fail("安装失败：", err.Error())
			}
		} else {
			parent := filepath.Dir(filePath)
			if !files.NewFile(parent).Exists() {
				err := os.MkdirAll(parent, 0751)
				if err != nil {
					this.Fail("安装失败：", err.Error())
				}
			}

			fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				this.Fail("安装失败：", err.Error())
			}
			_, _ = io.Copy(fp, tarReader)
			_ = fp.Close()
		}
	}

	this.start(mongodbDir)
}

func (this *InstallAction) start(mongodbDir string) {
	installStatus = "start"
	installPercent = 10

	// 检查 data 目录
	dataFile := files.NewFile(mongodbDir + "/data")
	if !dataFile.Exists() {
		err := dataFile.Mkdir()
		if err != nil {
			this.Fail("创建mongodb/data目录失败：", err.Error())
		}
	}

	// 启动
	args := []string{"--dbpath=" + mongodbDir + "/data", "--fork", "--logpath=" + mongodbDir + "/data/fork.log"}

	// 控制内存不能超过1G
	stat, err := mem.VirtualMemory()
	if err == nil && stat.Total > 0 {
		size := stat.Total / 1024 / 1024 / 1024
		if size >= 6 {
			args = append(args, "--wiredTigerCacheSizeGB=2")
		} else if size > 3 {
			args = append(args, "--wiredTigerCacheSizeGB=1")
		}
	}

	p := processes.NewProcess(mongodbDir+"/bin/mongod", args...)
	p.SetPwd(mongodbDir)

	logs.Println("start mongo:", mongodbDir+"/bin/mongod", strings.Join(args, " "))

	err = p.StartBackground()
	if err != nil {
		this.Fail("试图启动失败：", err.Error())
	}

	_ = p.Wait()

	this.check()
}

func (this *InstallAction) check() {
	// 保存数据库设置
	config := db.SharedDBConfig()
	config.Type = db.DBTypeMongo
	config.IsInitialized = true
	err := config.Save()
	if err != nil {
		logs.Error(err)
	}
	teadb.ChangeDB()

	// 检查状态
	installStatus = "check"
	installPercent = 10

	err = teadb.SharedDB().Test()
	if err != nil {
		this.Fail("仍然无法连接到MongoDb，请尝试手动安装")
	}

	installPercent = 100

	this.Success()
}
