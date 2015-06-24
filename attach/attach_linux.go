package attach

import (
	"net"
	"os"
	"path/filepath"
)

func New(pid string) (*Socket, error) {
	sockfile := filepath.Join(os.TempDir(), ".java_pid"+pid)
	addr, err := net.ResolveUnixAddr("unix", sockfile)
	if err != nil {
		return nil, err
	}

	// TODO: implement force attach feature.
	// see classes/sun/tools/attach/LinuxVirtualMachine.java in openjdk

	c, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return nil, err
	}
	return &Socket{c}, nil
}
