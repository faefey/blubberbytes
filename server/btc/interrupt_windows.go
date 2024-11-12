package btc

import (
	"log"
	"os/exec"

	"golang.org/x/sys/windows"
)

var sysProcAttr = &windows.SysProcAttr{
	CreationFlags: windows.CREATE_NEW_PROCESS_GROUP,
}

// Interrupt a command/process.
func InterruptCmd(cmd *exec.Cmd) {
	err := sendCtrlBreak(cmd.Process.Pid)
	if err != nil {
		log.Println(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(err)
	}
}

// Send Ctrl+Break to a process.
func sendCtrlBreak(pid int) error {
	d, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	defer d.Release()

	p, err := d.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return err
	}

	r, _, err := p.Call(windows.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		return err
	}

	return nil
}
