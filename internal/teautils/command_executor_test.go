package teautils

import "testing"

func TestExec(t *testing.T) {
	exec := NewCommandExecutor()
	exec.Add("ps", "ax")
	exec.Add("grep", "mysql")
	exec.Add("wc", "-l")
	output, err := exec.Run()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)
}

func TestExec2(t *testing.T) {
	exec := NewCommandExecutor()
	exec.Add("ps", "ax")
	exec.Add("grep", "mysql2")
	exec.Add("wc", "-l")
	output, err := exec.Run()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)
}
