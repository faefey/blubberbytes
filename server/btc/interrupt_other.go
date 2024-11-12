//go:build !windows
// +build !windows

package btc

import (
	"log"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

var sysProcAttr = &unix.SysProcAttr{
	Setpgid: true,
}

// Interrupt a command/process.
func InterruptCmd(cmd *exec.Cmd) {
	err := cmd.Process.Signal(os.Interrupt)
	if err != nil {
		log.Println(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(err)
	}
}
