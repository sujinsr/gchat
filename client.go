package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	ServerIP   string = "127.0.0.1"
	ServerPort string = "20099"
	stop       bool   = false
	debug      bool   = true
)

func clientSender(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	status := make([]byte, 8)
	it := 0

	fmt.Print("\n\tEnter Your Chat Name : ")

	for !stop {
		message, _ := reader.ReadBytes('\n')
		conn.Write(message[0 : len(message)-1])
		if it == 0 {

			n, _ := conn.Read(status)
			ret, _ := strconv.Atoi(string(status[:n]))
			if ret == 0 {
				go clientReceiver(conn)
			} else {
				fmt.Println("Same user alread available")
				stop = true
			}
			it++
		}
		fmt.Println("")
	}
	fmt.Println("Connection disconnected")
}

func clientReceiver(conn net.Conn) {
	msg := make([]byte, 1024)
	for !stop {
		n, err := conn.Read(msg)
		if err != nil {
			stop = true
		}
		msg := string(msg[:n])
		fmt.Println("------------------$", msg, "\n")
	}
	fmt.Println("Server disconnected")
}

func main() {
	conn, err := net.Dial("tcp", ServerIP+":"+ServerPort)
	if err != nil {
		fmt.Println(err, "tcp connect error")
		os.Exit(-1)
	}
	defer conn.Close()

	go clientSender(conn)

	/* hang on main untill go routine finish */
	for !stop {
		time.Sleep(1 * 1e9)
	}

	if debug == true {
		panic("Display stack trace")
	}
}
