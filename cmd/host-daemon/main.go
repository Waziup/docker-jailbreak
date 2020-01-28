package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

const maxStrSize = 5120

var socketFile = "/var/run/host.sock"

func main() {

	if len(os.Args) >= 2 {
		socketFile = os.Args[1]
	}

	listener, err := net.Listen("unix", socketFile)
	if err != nil {
		log.Fatal(err)
	}

	notifyDaemon()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go serve(conn)
	}
}

type byteReader struct {
	io.Reader
}

func (r *byteReader) ReadByte() (byte, error) {
	var c [1]byte
	_, err := io.ReadFull(r, c[:])
	return c[0], err
}

var errMaxStrSize = errors.New("maximum string size exceeded")
var errMaxArgs = errors.New("maximum number of arguments exceeded")

func readString(r io.Reader) (string, error) {
	var br = byteReader{r}
	s, err := binary.ReadUvarint(&br)
	if err != nil {
		return "", err
	}
	if s == 0 {
		return "", nil
	}
	if s > maxStrSize {
		return "", errMaxStrSize
	}
	buf := make([]byte, s)
	_, err = io.ReadFull(r, buf)
	return string(buf), err
}

func serve(conn net.Conn) {
	addr := conn.RemoteAddr()
	defer conn.Close()

	name, err := readString(conn)
	if err != nil {
		log.Printf("[%s] Err %v", addr, err)
		return
	}
	var args []string
	for i := 0; ; i++ {
		if i == 64 {
			log.Printf("[%s] Err %v", addr, errMaxArgs)
		}
		arg, err := readString(conn)
		if err != nil {
			log.Printf("[%s] Err %v", addr, err)
			return
		}
		if arg == "" {
			break
		}
		args = append(args, arg)
	}

	log.Printf("[%s] $ %s %v", addr, name, args)
	cmd := exec.Command(name, args...)
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	err = cmd.Run()
	if err != nil {
		log.Printf("[%s] Err %v", addr, err)
	} else {
		log.Printf("[%s] Completed successfully.", addr)
	}
}
