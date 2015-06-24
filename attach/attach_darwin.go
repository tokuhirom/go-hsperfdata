package attach

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// see classes/sun/tools/attach/BsdVirtualMachine.java in openjdk
func force_attach(pid int) error {
	attach_file := filepath.Join(os.TempDir(), fmt.Sprintf(".attach_pid%d", pid))
	f, err := os.Create(attach_file)
	if err != nil {
		return fmt.Errorf("Canot create file:%v:%v", attach_file, err)
	}
	f.Close()

	err = syscall.Kill(pid, syscall.SIGQUIT)
	if err != nil {
		return fmt.Errorf("Canot send sigkill:%v:%v", pid, err)
	}

	sockfile := filepath.Join(os.TempDir(), fmt.Sprintf(".java_pid%d", pid))
	for i := 0; i < 10; i++ {
		if exists(sockfile) {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("Canot attach process:%v", pid)
}

func GetSocketFile(pid int) (string, error) {
	sockfile := filepath.Join(os.TempDir(), fmt.Sprintf(".java_pid%d", pid))
	if !exists(sockfile) {
		err := force_attach(pid)
		if err != nil {
			return "", err
		}
	}
	return sockfile, nil
}

func New(pid int) (*Socket, error) {
	sockfile, err := GetSocketFile(pid)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveUnixAddr("unix", sockfile)
	if err != nil {
		return nil, err
	}

	c, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return nil, err
	}
	return &Socket{c}, nil
}
