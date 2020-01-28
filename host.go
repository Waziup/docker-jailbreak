// Package host provides a function to run commands on the host machine.
// It communicates with a daemon service that is running on the host machine.
// This can be used to run commands from inside a docker container.
//
//   cmd, err := host.Exec("date")
//   if err != nil {
//     log.Fatal(err)
//   }
//   defer cmd.Close()
//   io.Copy(os.Stdout, cmd)
//
package host

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

// SocketFile is the unix socket name that the host serves on.
var SocketFile = "/var/run/host.sock"

func writeString(w io.Writer, str string) error {
	var buf [8]byte
	data := []byte(str)
	n := binary.PutUvarint(buf[:], uint64(len(data)))
	_, err := w.Write(buf[0:n])
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

var errEmptyName = errors.New("can not send empty name")
var errEmptyArg = errors.New("can not send empty arg")

// Exec runs a new command on the host system.
func Exec(name string, args ...string) (io.ReadWriteCloser, error) {
	conn, err := net.Dial("unix", SocketFile)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errEmptyName
	}
	if err = writeString(conn, name); err != nil {
		return nil, err
	}
	for _, arg := range args {
		if arg == "" {
			return nil, errEmptyArg
		}
		if err = writeString(conn, arg); err != nil {
			return nil, err
		}
	}
	writeString(conn, "")
	return conn, nil
}
