package teaproxy

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/TeaWeb/build/internal/teaproxy/mitm"
	"github.com/iwind/TeaGo/Tea"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestRequest_Forward_CA(t *testing.T) {
	ca, privateKey, err := mitm.NewAuthority("teaweb.teaos.cn", "TeaWeb Authority", 10*365*24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", ca.PublicKey)

	{
		writer, err := os.OpenFile("teaweb.proxy.pem", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			t.Fatal(err)
		}

		err = pem.Encode(writer, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: ca.Raw,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		writer, err := os.OpenFile("teaweb.proxy.key", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			t.Fatal(err)
		}

		err = pem.Encode(writer, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("OK")
}

func TestRequest_Forward_Decode(t *testing.T) {
	{
		data, err := ioutil.ReadFile(Tea.Root + "/web/certs/teaweb.proxy.pem")
		if err != nil {
			t.Fatal(err)
		}

		for {
			block, rest := pem.Decode(data)
			if err != nil {
				t.Fatal(err)
			}

			t.Log(x509.ParseCertificate(block.Bytes))

			if len(rest) == 0 {
				break
			}
			data = rest
		}
	}

	{
		data, err := ioutil.ReadFile(Tea.Root + "/web/certs/teaweb.proxy.key")
		if err != nil {
			t.Fatal(err)
		}

		block, _ := pem.Decode(data)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(x509.ParsePKCS1PrivateKey(block.Bytes))
	}
}
