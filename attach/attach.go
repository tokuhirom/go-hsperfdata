package attach

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

const PROTOCOL_VERSION = "1"
const ATTACH_ERROR_BADVERSION = 101

type Socket struct {
	sock *net.UnixConn
}

func New(pid string) (*Socket, error) {
	sockfile := "/tmp/.java_pid" + pid
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

func (sock *Socket) Read(b []byte) (int, error) {
	return sock.sock.Read(b)
}

func (sock *Socket) ReadString() (string, error) {
	retval := ""
	for {
		buf := make([]byte, 1024)
		read, err := sock.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return retval, err
		}
		retval += string(buf[0 : read-1])
	}
	return retval, nil
}

// see classes/sun/tools/attach/LinuxVirtualMachine.java
// see classes/sun/tools/attach/HotSpotVirtualMachine.java for available commands
func (sock *Socket) Execute(cmd string, args ...string) error {
	err := sock.writeString(PROTOCOL_VERSION)
	if err != nil {
		return err
	}
	err = sock.writeString(cmd)
	if err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		if len(args) > i {
			err = sock.writeString(args[i])
			if err != nil {
				return err
			}
		} else {
			err = sock.writeString("")
			if err != nil {
				return err
			}
		}
	}

	i, err := sock.readInt()
	if i != 0 {
		if i == ATTACH_ERROR_BADVERSION {
			return fmt.Errorf("Protocol mismatch with target VM")
		} else {
			return fmt.Errorf("Command failed in target VM")
		}
	}
	return err
}

func (sock *Socket) readInt() (int, error) {
	b := make([]byte, 1)
	buf := make([]byte, 0)
	for {
		_, err := sock.Read(b)
		if err != nil {
			return 0, err
		}
		if '0' <= b[0] && b[0] <= '9' {
			buf = append(buf, b[0])
			continue
		}

		if len(buf) == 0 {
			return 0, fmt.Errorf("cannot read int")
		} else {
			return strconv.Atoi(string(buf))
		}
	}
}

func (sock *Socket) writeString(s string) error {
	return sock.write([]byte(s))
}

func (sock *Socket) write(bytes []byte) error {
	{
		wrote, err := sock.sock.Write(bytes)
		if err != nil {
			return err
		}
		if wrote != len(bytes) {
			return fmt.Errorf("cannot write")
		}
	}
	{
		wrote, err := sock.sock.Write([]byte("\x00"))
		if err != nil {
			return err
		}
		if wrote != 1 {
			return fmt.Errorf("cannot write null byte")
		}
	}
	return nil
}
