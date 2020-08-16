package notices

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"testing"
)

func TestNoticeScriptMedia_Send(t *testing.T) {
	script := `#!/usr/bin/env bash

echo  "subject:${NoticeSubject}"
echo "body:${NoticeBody}"
`

	tmp := files.NewFile(Tea.Root + "/web/tmp/media_test.sh")
	err := tmp.WriteString(script)
	if err != nil {
		t.Fatal(err)
	}
	_ = tmp.Chmod(0777)
	defer func() {
		_ = tmp.Delete()
	}()

	media := NewNoticeScriptMedia()
	media.Path = tmp.Path()
	_, err = media.Send("zhangsan", "this is subject", "this is body")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoticeScriptMedia_Send2(t *testing.T) {
	media := NewNoticeScriptMedia()
	media.ScriptType = "code"
	media.Script = `#!/usr/bin/env bash

echo  "subject:${NoticeSubject}"
echo "body:${NoticeBody}"
`
	_, err := media.Send("zhangsan", "this is subject", "this is body")
	if err != nil {
		t.Fatal(err)
	}
}
