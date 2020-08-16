package certs

import (
	"archive/zip"
	"bytes"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"path/filepath"
)

type DownloadAction actions.Action

// 下载
func (this *DownloadAction) RunGet(params struct {
	CertId string
	Type   string
}) {
	cert := teaconfigs.SharedSSLCertList().FindCert(params.CertId)
	if cert == nil {
		this.WriteString("找不到要查看的证书")
		return
	}

	// 下载zip
	if params.Type == "zip" {
		certData, err := cert.ReadCert()
		if err != nil {
			this.Fail("读取证书失败：" + err.Error())
		}

		keyData := []byte{}
		if len(cert.KeyFile) > 0 {
			keyData, err = cert.ReadKey()
			if err != nil {
				this.Fail("读取私钥失败：" + err.Error())
			}
		}

		buffer := bytes.NewBuffer([]byte{})
		z := zip.NewWriter(buffer)

		{
			w, err := z.Create(filepath.Base(cert.FullCertPath()))
			if err != nil {
				this.Fail("创建ZIP失败：" + err.Error())
			}
			_, err = w.Write(certData)
			if err != nil {
				this.Fail("写入证书失败：" + err.Error())
			}
		}

		if len(keyData) > 0 {
			w, err := z.Create(filepath.Base(cert.FullKeyPath()))
			if err != nil {
				this.Fail("创建ZIP失败：" + err.Error())
			}
			_, err = w.Write(keyData)
			if err != nil {
				this.Fail("写入私钥失败：" + err.Error())
			}
		}

		err = z.Flush()
		if err != nil {
			this.Fail("创建ZIP失败：" + err.Error())
		}

		err = z.Close()
		if err != nil {
			logs.Error(err)
		}

		this.AddHeader("Content-Disposition", "attachment; filename=\"certificate.zip\";")
		this.Write(buffer.Bytes())
		return
	}

	// 下载证书
	if params.Type == "cert" {
		data, err := cert.ReadCert()
		if err != nil {
			this.WriteString(err.Error())
			return
		}

		this.AddHeader("Content-Disposition", "attachment; filename=\""+filepath.Base(cert.CertFile)+"\";")
		this.Write(data)
		return
	}

	// 下载私钥
	if params.Type == "key" {
		data, err := cert.ReadKey()
		if err != nil {
			this.WriteString(err.Error())
			return
		}

		this.AddHeader("Content-Disposition", "attachment; filename=\""+filepath.Base(cert.KeyFile)+"\";")
		this.Write(data)
		return
	}

	// 查看证书
	if params.Type == "viewCert" {
		data, err := cert.ReadCert()
		if err != nil {
			this.WriteString(err.Error())
			return
		}

		this.Write(data)
		return
	}

	// 查看私钥
	if params.Type == "viewKey" {
		data, err := cert.ReadKey()
		if err != nil {
			this.WriteString(err.Error())
			return
		}

		this.Write(data)
		return
	}
}
