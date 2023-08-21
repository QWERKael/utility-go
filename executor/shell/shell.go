package shell

import (
	"bytes"
	"errors"
	"os/exec"
)

type CommandStatus uint8

const (
	Null       CommandStatus = iota // 未设置命令
	Prepare                         // 已设置命令，未执行
	Done                            // 执行完成，不允许再次执行
	Repeatable                      // 可重复执行的命令
	Running                         // 执行中
)

type ShCommand struct {
	Cmd        *exec.Cmd
	Status     CommandStatus
	IsCombined bool
}

func NewShCommand(isCombined bool) *ShCommand {
	return &ShCommand{
		Status:     Null,
		IsCombined: isCombined,
	}
}

func (sc *ShCommand) SetCommand(name string, args ...string) {
	sc.Cmd = exec.Command(name, args...)
	sc.Status = Prepare
}

func (sc *ShCommand) Execute() (bytes.Buffer, bytes.Buffer, error) {
	if sc.Status != Prepare && sc.Status != Repeatable {
		return bytes.Buffer{}, bytes.Buffer{}, errors.New("ShCommand的状态不正确")
	}
	var stdout, stderr bytes.Buffer
	if sc.IsCombined {
		sc.Cmd.Stdout = &stdout
		sc.Cmd.Stderr = &stdout
	} else {
		sc.Cmd.Stdout = &stdout
		sc.Cmd.Stderr = &stderr
	}
	err := sc.Cmd.Run()
	sc.Status = Done
	return stdout, stderr, err
}

func (sc *ShCommand) AsyncExecuteWithoutOutput() error {
	if sc.Status != Prepare && sc.Status != Repeatable {
		return errors.New("ShCommand的状态不正确")
	}
	err := sc.Cmd.Start()
	sc.Status = Running
	return err
}
