package main

import (
	"fmt"
	"github.com/dotcloud/docker/pkg/system"
	"github.com/docker/libcontainer/namespaces"
	"os/exec"
	"syscall"
	"os"
)

func main() {
	term := namespaces.NewTerminal(os.Stdin, os.Stdout, os.Stderr, true)

	master, slave, err := system.CreateMasterAndConsole()
	if err != nil {
		fmt.Println("system.CreateMasterAndConsole failed: ", err)
		return
	}

	term.SetMaster(master)

	err = term.Attach(nil)
	if err != nil {
		fmt.Println("term.Attach failed: ", err)
		return
	}
	defer term.Close()

	fmt.Println("os.Args = ", os.Args)

	command := exec.Command(os.Args[1], slave)
	system.SetCloneFlags(command, uintptr(syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNS |
				syscall.CLONE_NEWPID | syscall.CLONE_NEWNET))

	if output, err := command.Output(); err == nil {
		fmt.Println("command.Output: ", string(output[:]))
	} else {
		fmt.Println("command.Output error = ", err)
		return
	}

	fmt.Println("testpts succeeded")
}
