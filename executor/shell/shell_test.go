package shell

import (
	"log"
	"testing"
)

func TestExecute(t *testing.T) {
	sc := NewShCommand(false)
	sc.SetCommand("ls", "-lhz")
	stdout, stderr, err := sc.Execute()
	log.Printf("\nstdout: %s\nstderr: %s", string(stdout.Bytes()), string(stderr.Bytes()))
	if err != nil {
		log.Fatalf("执行失败！")
	}
}
