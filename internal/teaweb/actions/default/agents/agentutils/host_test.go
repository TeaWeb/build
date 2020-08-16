package agentutils

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
	"time"
)

func TestParseHostRules(t *testing.T) {
	result := ParseHostRules("a\nb", 1)
	t.Log(len(result))
	t.Log(result)
}

func TestParseHostRules_Range(t *testing.T) {
	result := ParseHostRules("192.[01-10].1.[2-20]", -1)
	t.Log(len(result))
	logs.PrintAsJSON(result, t)
}

func TestParseHostRules_Range2(t *testing.T) {
	result := ParseHostRules("web[001-022]", 10)
	t.Log(len(result))
	logs.PrintAsJSON(result, t)
}

func TestCheckHostConnectivity(t *testing.T) {
	t.Log(CheckHostConnectivity("127.0.0.1", 28, 3*time.Second))
	t.Log(CheckHostConnectivity("192.168.2.33", 22, 3*time.Second))
}

func TestCheckHostConnectivity_Timeout(t *testing.T) {
	t.Log(CheckHostConnectivity("192.168.2.1", 22, 3*time.Second))
}
