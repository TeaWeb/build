package checkpoints

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"testing"
)

func TestRequestUploadCheckpoint_RequestValue(t *testing.T) {
	body := bytes.NewBuffer([]byte{})

	writer := multipart.NewWriter(body)

	{
		part, err := writer.CreateFormField("name")
		if err == nil {
			part.Write([]byte("lu"))
		}
	}

	{
		part, err := writer.CreateFormField("age")
		if err == nil {
			part.Write([]byte("20"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile", "hello.txt")
		if err == nil {
			part.Write([]byte("Hello, World!"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile2", "hello.PHP")
		if err == nil {
			part.Write([]byte("Hello, World, PHP!"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile3", "hello.asp")
		if err == nil {
			part.Write([]byte("Hello, World, ASP Pages!"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile4", "hello.asp")
		if err == nil {
			part.Write([]byte("Hello, World, ASP Pages!"))
		}
	}

	writer.Close()

	rawReq, err := http.NewRequest(http.MethodPost, "http://teaos.cn/", body)
	if err != nil {
		t.Fatal()
	}

	req := requests.NewRequest(rawReq)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	checkpoint := new(RequestUploadCheckpoint)
	t.Log(checkpoint.RequestValue(req, "field", nil))
	t.Log(checkpoint.RequestValue(req, "minSize", nil))
	t.Log(checkpoint.RequestValue(req, "maxSize", nil))
	t.Log(checkpoint.RequestValue(req, "name", nil))
	t.Log(checkpoint.RequestValue(req, "ext", nil))

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
