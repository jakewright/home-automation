package ola

import (
	"fmt"
	"os/exec"
	"testing"
)

func Test_test(t *testing.T) {
	cmd := exec.Command("ls", "-l", "-a")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	fmt.Printf("Output: %s\n", string(out))

	t.Fail()
}
