package teaconfigs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"github.com/TeaWeb/build/internal/lego/acme"
	"github.com/TeaWeb/build/internal/lego/certcrypto"
	"github.com/TeaWeb/build/internal/lego/certificate"
	"github.com/TeaWeb/build/internal/lego/challenge/dns01"
	"github.com/TeaWeb/build/internal/lego/lego"
	"github.com/TeaWeb/build/internal/lego/log"
	"github.com/TeaWeb/build/internal/lego/registration"
	"github.com/iwind/TeaGo/utils/time"
	"io/ioutil"
	golog "log"
	"strings"
)

// ACME DNS Request
type ACMERequest struct {
	User *ACMELocalUser `yaml:"user" json:"user"`

	Domains    []string         `yaml:"domains" json:"domains"`
	DNSRecords []*ACMEDNSRecord `yaml:"dnsRecords" json:"dnsRecords"`

	CertURL string `yaml:"certURL" json:"certURL"`
	Cert    string `yaml:"cert" json:"cert"`
	Key     string `yaml:"key" json:"key"`
}

// 获取新对象
func NewACMERequest() *ACMERequest {
	return &ACMERequest{}
}

// 获取连接客户端
func (this *ACMERequest) Client() (client *lego.Client, err error) {
	log.Logger.(*golog.Logger).SetOutput(ioutil.Discard)

	if len(this.User.Email) == 0 {
		return nil, errors.New("user email should not be empty")
	}

	if len(this.User.Key) == 0 {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}

		privateKeyData, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return nil, err
		}
		this.User.Key = base64.StdEncoding.EncodeToString(privateKeyData)
	}

	userKeyBytes, err := base64.StdEncoding.DecodeString(this.User.Key)
	if err != nil {
		return nil, err
	}

	userKey, err := x509.ParsePKCS8PrivateKey(userKeyBytes)
	if err != nil {
		return nil, err
	}

	user := ACMEUser{
		Email: this.User.Email,
		Key:   userKey,
	}

	if len(this.User.URI) > 0 {
		user.Registration = &registration.Resource{
			Body: acme.Account{
				Status:               "valid",
				Contact:              []string{"mailto:" + this.User.Email,},
				TermsOfServiceAgreed: true,
			},
			URI: this.User.URI,
		}
	}

	config := lego.NewConfig(&user)
	config.Certificate.KeyType = certcrypto.RSA2048
	client, err = lego.NewClient(config)
	if err != nil {
		return
	}

	err = client.Challenge.SetDNS01Provider(NewACMEDNSProvider("teaweb"), dns01.WrapPreCheck(func(domain, fqdn, value string, check dns01.PreCheckFunc) (b bool, e error) {
		b = true
		return
	}))
	if err != nil {
		return
	}

	if len(this.User.URI) == 0 {
		reg, err := client.Registration.Register(registration.RegisterOptions{
			TermsOfServiceAgreed: true,
		})
		if err != nil {
			return client, err
		}

		user.Registration = reg
		this.User.URI = reg.URI
	}

	return
}

// 获取要设置的DNS记录
func (this *ACMERequest) RetrieveDNSRecords(client *lego.Client) (records []*ACMEDNSRecord, err error) {
	if len(this.Domains) == 0 {
		return nil, errors.New("domains should not be empty")
	}

	recordStrings, err := client.Certificate.GetRecords(certificate.ObtainRequest{
		Domains: this.Domains,
		Bundle:  true,
	})

	if err != nil {
		return nil, err
	}

	for _, recordString := range recordStrings {
		pieces := strings.SplitN(recordString, "|", 2)
		if len(pieces) == 2 {
			records = append(records, &ACMEDNSRecord{
				FQDN:  pieces[0],
				Value: pieces[1],
			})
		}
	}

	return
}

// 获取证书信息
func (this *ACMERequest) Retrieve(client *lego.Client) error {
	if len(this.Domains) == 0 {
		return errors.New("domains should not be empty")
	}

	request := certificate.ObtainRequest{
		Domains: this.Domains,
		Bundle:  true,
	}

	resource, err := client.Certificate.Obtain(request)
	if err != nil {
		return err
	}
	this.Cert = string(resource.Certificate)
	this.Key = string(resource.PrivateKey)
	this.CertURL = resource.CertURL

	return nil
}

// 更新
func (this *ACMERequest) Renew(client *lego.Client) error {
	return this.Retrieve(client)
}

// 获取证书对象
func (this *ACMERequest) CertObject() (*tls.Certificate, error) {
	if len(this.Cert) == 0 || len(this.Key) == 0 {
		return nil, nil
	}

	cert, err := tls.X509KeyPair([]byte(this.Cert), []byte(this.Key))
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// 获取证书时间
func (this *ACMERequest) CertDate() [2]string {
	cert, err := this.CertObject()
	if err != nil {
		return [2]string{}
	}
	for _, data := range cert.Certificate {
		c, err := x509.ParseCertificate(data)
		if err != nil {
			continue
		}
		return [2]string{timeutil.Format("Y-m-d", c.NotBefore), timeutil.Format("Y-m-d", c.NotAfter)}
	}
	return [2]string{}
}

// 写入证书文件
func (this *ACMERequest) WriteCertFile(path string) error {
	return ioutil.WriteFile(path, []byte(this.Cert), 0666)
}

// 写入密钥文件
func (this *ACMERequest) WriteKeyFile(path string) error {
	return ioutil.WriteFile(path, []byte(this.Key), 0666)
}
