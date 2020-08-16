package cmd

import "testing"

func TestWebShell_CheckPid(t *testing.T) {
	shell := &WebShell{}
	t.Log(shell.checkPid())
}
