//go:build windows

package sampledconn

import (
	"errors"
	"golang.org/x/sys/windows"
	"syscall"
)

func OSPeekConn(conn syscall.Conn) (PeekedBytes, error) {
	s := PeekedBytes{}

	rawConn, err := conn.SyscallConn()
	if err != nil {
		return s, err
	}

	readBytes := 0
	var readErr error
	err = rawConn.Read(func(fd uintptr) bool {
		for readBytes < peekSize {
			var n uint32
			flags := uint32(windows.MSG_PEEK)
			wsabuf := windows.WSABuf{
				Len: uint32(len(s) - readBytes),
				Buf: &s[readBytes],
			}

			readErr = windows.WSARecv(windows.Handle(fd), &wsabuf, 1, &n, &flags, nil, nil)
			if errors.Is(readErr, windows.WSAEWOULDBLOCK) {
				return false
			}
			if readErr != nil {
				return true
			}
			readBytes += int(n)
		}
		return true
	})
	if readErr != nil {
		return s, readErr
	}
	if err != nil {
		return s, err
	}

	return s, nil
}
