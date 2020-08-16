package proxyutils

import "testing"

func TestCheckChartChanges(t *testing.T) {
	t.Log(CheckChartChanges())
}

func TestApplyChartChanges(t *testing.T) {
	t.Log(ApplyChartChanges())
}