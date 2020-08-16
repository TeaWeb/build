package teawaf

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"github.com/dchest/siphash"
	"github.com/iwind/TeaGo/assert"
	"net/http"
	"regexp"
	"runtime"
	"testing"
)

func TestRuleSet_MatchRequest(t *testing.T) {
	set := NewRuleSet()
	set.Connector = RuleConnectorAnd

	set.Rules = []*Rule{
		{
			Param:    "${arg.name}",
			Operator: RuleOperatorEqString,
			Value:    "lu",
		},
		{
			Param:    "${arg.age}",
			Operator: RuleOperatorEq,
			Value:    "20",
		},
	}

	err := set.Init()
	if err != nil {
		t.Fatal(err)
	}

	rawReq, err := http.NewRequest(http.MethodGet, "http://teaos.cn/hello?name=lu&age=20", nil)
	if err != nil {
		t.Fatal(err)
	}
	req := requests.NewRequest(rawReq)
	t.Log(set.MatchRequest(req))
}

func TestRuleSet_MatchRequest2(t *testing.T) {
	a := assert.NewAssertion(t)

	set := NewRuleSet()
	set.Connector = RuleConnectorOr

	set.Rules = []*Rule{
		{
			Param:    "${arg.name}",
			Operator: RuleOperatorEqString,
			Value:    "lu",
		},
		{
			Param:    "${arg.age}",
			Operator: RuleOperatorEq,
			Value:    "21",
		},
	}

	err := set.Init()
	if err != nil {
		t.Fatal(err)
	}

	rawReq, err := http.NewRequest(http.MethodGet, "http://teaos.cn/hello?name=lu&age=20", nil)
	if err != nil {
		t.Fatal(err)
	}
	req := requests.NewRequest(rawReq)
	a.IsTrue(set.MatchRequest(req))
}

func BenchmarkRuleSet_MatchRequest(b *testing.B) {
	runtime.GOMAXPROCS(1)

	set := NewRuleSet()
	set.Connector = RuleConnectorOr

	set.Rules = []*Rule{
		{
			Param:    "${requestAll}",
			Operator: RuleOperatorMatch,
			Value:    `(onmouseover|onmousemove|onmousedown|onmouseup|onerror|onload|onclick|ondblclick|onkeydown|onkeyup|onkeypress)\s*=`,
		},
		{
			Param:    "${requestAll}",
			Operator: RuleOperatorMatch,
			Value:    `\b(eval|system|exec|execute|passthru|shell_exec|phpinfo)\s*\(`,
		},
		{
			Param:    "${arg.name}",
			Operator: RuleOperatorEqString,
			Value:    "lu",
		},
		{
			Param:    "${arg.age}",
			Operator: RuleOperatorEq,
			Value:    "21",
		},
	}

	err := set.Init()
	if err != nil {
		b.Fatal(err)
	}

	rawReq, err := http.NewRequest(http.MethodPost, "http://teaos.cn/hello?name=lu&age=20", bytes.NewBuffer(bytes.Repeat([]byte("HELLO"), 1024)))
	if err != nil {
		b.Fatal(err)
	}
	req := requests.NewRequest(rawReq)
	for i := 0; i < b.N; i++ {
		_, _ = set.MatchRequest(req)
	}
}

func BenchmarkRuleSet_MatchRequest_Regexp(b *testing.B) {
	runtime.GOMAXPROCS(1)

	set := NewRuleSet()
	set.Connector = RuleConnectorOr

	set.Rules = []*Rule{
		{
			Param:             "${requestBody}",
			Operator:          RuleOperatorMatch,
			Value:             `\b(eval|system|exec|execute|passthru|shell_exec|phpinfo)\s*\(`,
			IsCaseInsensitive: false,
		},
	}

	err := set.Init()
	if err != nil {
		b.Fatal(err)
	}

	rawReq, err := http.NewRequest(http.MethodPost, "http://teaos.cn/hello?name=lu&age=20", bytes.NewBuffer(bytes.Repeat([]byte("HELLO"), 2048)))
	if err != nil {
		b.Fatal(err)
	}
	req := requests.NewRequest(rawReq)
	for i := 0; i < b.N; i++ {
		_, _ = set.MatchRequest(req)
	}
}

func BenchmarkRuleSet_MatchRequest_Regexp2(b *testing.B) {
	reg, err := regexp.Compile(`(?iU)\b(eval|system|exec|execute|passthru|shell_exec|phpinfo)\b`)
	if err != nil {
		b.Fatal(err)
	}

	buf := bytes.Repeat([]byte(" HELLO "), 10240)

	for i := 0; i < b.N; i++ {
		_ = reg.Match(buf)
	}
}

func BenchmarkRuleSet_MatchRequest_Regexp3(b *testing.B) {
	reg, err := regexp.Compile(`(?iU)^(eval|system|exec|execute|passthru|shell_exec|phpinfo)`)
	if err != nil {
		b.Fatal(err)
	}

	buf := bytes.Repeat([]byte(" HELLO "), 1024)

	for i := 0; i < b.N; i++ {
		_ = reg.Match(buf)
	}
}

func BenchmarkHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = siphash.Hash(0, 0, bytes.Repeat([]byte("HELLO"), 10240))
	}
}
