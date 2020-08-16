package teaconfigs

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/types"
	"io/ioutil"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// TLS Version
type TLSVersion = string

// Cipher Suites
type TLSCipherSuite = string

// SSL配置
type SSLConfig struct {
	On bool `yaml:"on" json:"on"` // 是否开启

	Certificate    string `yaml:"certificate" json:"certificate"`       // 证书文件, deprecated in v0.1.4
	CertificateKey string `yaml:"certificateKey" json:"certificateKey"` // 密钥, deprecated in v0.1.4

	Certs           []*SSLCertConfig  `yaml:"certs" json:"certs"`
	CertTasks       []*SSLCertTask    `yaml:"certTasks" json:"certTasks"`
	ClientAuthType  SSLClientAuthType `yaml:"clientAuthType" json:"clientAuthType"`   // 客户端认证类型
	ClientCACertIds []string          `yaml:"clientCACertIds" json:"clientCACertIds"` // 客户端认证CA

	Listen       []string         `yaml:"listen" json:"listen"`             // 网络地址
	MinVersion   TLSVersion       `yaml:"minVersion" json:"minVersion"`     // 支持的最小版本
	CipherSuites []TLSCipherSuite `yaml:"cipherSuites" json:"cipherSuites"` // 加密算法套件

	HSTS          *HSTSConfig `yaml:"hsts2" json:"hsts"`                  // HSTS配置，yaml之所以使用hsts2，是因为要和以前的版本分开
	HTTP2Disabled bool        `yaml:"http2Disabled" json:"http2Disabled"` // 是否禁用HTTP2

	nameMapping map[string]*tls.Certificate // dnsName => cert

	minVersion   uint16
	cipherSuites []uint16

	clientCAPool *x509.CertPool
}

// 获取新对象
func NewSSLConfig() *SSLConfig {
	return &SSLConfig{}
}

// 校验配置
func (this *SSLConfig) Validate() error {
	if !this.On {
		return nil
	}

	if len(this.Certs) == 0 {
		if len(this.Certificate) > 0 && len(this.CertificateKey) > 0 {
			newCert := NewSSLCertConfig(this.Certificate, this.CertificateKey)
			newCert.Id = "old_version_cert"
			this.Certs = append(this.Certs, newCert)
		}
	}

	if len(this.Certs) == 0 {
		return errors.New("no certificates in https config")
	}

	for _, cert := range this.Certs {
		err := cert.Validate()
		if err != nil {
			return err
		}
	}

	if this.Listen == nil {
		this.Listen = []string{}
	} else {
		for index, addr := range this.Listen {
			_, _, err := net.SplitHostPort(addr)
			if err != nil {
				this.Listen[index] = strings.TrimSuffix(addr, ":") + ":443"
			}
		}
	}

	// min version
	this.convertMinVersion()

	// cipher suite categories
	this.initCipherSuites()

	// hsts
	if this.HSTS != nil {
		err := this.HSTS.Validate()
		if err != nil {
			return err
		}
	}

	// CA证书
	if len(this.ClientCACertIds) > 0 && this.ClientAuthType != SSLClientAuthTypeNoClientCert {
		this.clientCAPool = x509.NewCertPool()
		list := SharedSSLCertList()
		for _, certId := range this.ClientCACertIds {
			cert := list.FindCert(certId)
			if cert == nil {
				continue
			}
			if !cert.On {
				continue
			}
			data, err := ioutil.ReadFile(cert.FullCertPath())
			if err != nil {
				return err
			}
			this.clientCAPool.AppendCertsFromPEM(data)
		}
	}

	return nil
}

// 取得最小版本
func (this *SSLConfig) TLSMinVersion() uint16 {
	return this.minVersion
}

// 套件
func (this *SSLConfig) TLSCipherSuites() []uint16 {
	return this.cipherSuites
}

// 校验是否匹配某个域名
func (this *SSLConfig) MatchDomain(domain string) (cert *tls.Certificate, ok bool) {
	for _, cert := range this.Certs {
		if cert.MatchDomain(domain) {
			return cert.CertObject(), true
		}
	}
	return nil, false
}

// 取得第一个证书
func (this *SSLConfig) FirstCert() *tls.Certificate {
	for _, cert := range this.Certs {
		return cert.CertObject()
	}
	return nil
}

// 是否包含某个证书或密钥路径
func (this *SSLConfig) ContainsFile(file string) bool {
	for _, cert := range this.Certs {
		if cert.CertFile == file || cert.KeyFile == file {
			return true
		}
	}
	return false
}

// 删除证书文件
func (this *SSLConfig) DeleteFiles() error {
	var resultErr error = nil

	if len(this.Certificate) > 0 {
		err := files.NewFile(Tea.ConfigFile(this.Certificate)).Delete()
		if err != nil {
			resultErr = err
		}
	}

	if len(this.CertificateKey) > 0 {
		err := files.NewFile(Tea.ConfigFile(this.CertificateKey)).Delete()
		if err != nil {
			resultErr = err
		}
	}

	for _, cert := range this.Certs {
		err := cert.DeleteFiles()
		if err != nil {
			resultErr = err
		}
	}

	return resultErr
}

// 查找单个证书配置
func (this *SSLConfig) FindCert(certId string) *SSLCertConfig {
	for _, cert := range this.Certs {
		if cert.Id == certId {
			return cert
		}
	}
	return nil
}

// 添加证书
func (this *SSLConfig) AddCert(cert *SSLCertConfig) {
	this.Certs = append(this.Certs, cert)
}

// 添加证书任务
func (this *SSLConfig) AddCertTask(certTask *SSLCertTask) {
	this.CertTasks = append(this.CertTasks, certTask)
}

// 删除证书任务
func (this *SSLConfig) RemoveCertTask(certTaskId string) {
	result := []*SSLCertTask{}
	for _, task := range this.CertTasks {
		if task.Id == certTaskId {
			continue
		}
		result = append(result, task)
	}
	this.CertTasks = result
}

// 查找证书任务
func (this *SSLConfig) FindCertTask(certTaskId string) *SSLCertTask {
	for _, task := range this.CertTasks {
		if task.Id == certTaskId {
			return task
		}
	}
	return nil
}

// CA证书Pool，用于TLS对客户端进行认证
func (this *SSLConfig) CAPool() *x509.CertPool {
	return this.clientCAPool
}

// 分解所有监听地址
func (this *SSLConfig) ParseListenAddresses() []string {
	result := []string{}
	var reg = regexp.MustCompile(`\[\s*(\d+)\s*[,:-]\s*(\d+)\s*]$`)
	for _, addr := range this.Listen {
		match := reg.FindStringSubmatch(addr)
		if len(match) == 0 {
			result = append(result, addr)
		} else {
			min := types.Int(match[1])
			max := types.Int(match[2])
			if min > max {
				min, max = max, min
			}
			for i := min; i <= max; i++ {
				newAddr := reg.ReplaceAllString(addr, ":"+strconv.Itoa(i))
				result = append(result, newAddr)
			}
		}
	}
	return result
}
