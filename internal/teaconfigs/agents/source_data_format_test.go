package agents

import "testing"

func TestDecodeSource(t *testing.T) {
	{
		data := []byte("123")
		t.Log(DecodeSource(data, SourceDataFormatSingeLine))
	}

	{
		data := []byte(" 123 \n")
		t.Log(DecodeSource(data, SourceDataFormatSingeLine))
	}

	{
		data := []byte("123\n 456 \n  789")
		t.Log(DecodeSource(data, SourceDataFormatMultipleLine))
	}

	{
		data := []byte(`{
	"name": "lu",
	"age": 22
}`)
		t.Log(DecodeSource(data, SourceDataFormatJSON))
	}

	{
		data := []byte(`
name: lu
age: 22
`)
		t.Log(DecodeSource(data, SourceDataFormatYAML))
	}
}

func TestFindSourceDataFormat(t *testing.T) {
	t.Log(FindSourceDataFormat(SourceDataFormatJSON))
}
