package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	ServerIP   string = "127.0.0.1"
	ServerPort string = "20099"
)

func errorCheck(err error, errStr string) {
	if err != nil {
		fmt.Println(errStr)
		os.Exit(-1)
	}
}

func clientSender(conn net.Conn) {
	fmt.Print("\n\tEnter Your Chat Name : ")
	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadBytes('\n')
		conn.Write(message[0 : len(message)-1])
		fmt.Println("")
	}
}

func clientReceiver(conn net.Conn) {
	msg := make([]byte, 1024)
	for {
		n, err := conn.Read(msg)
		errorCheck(err, "Read error")
		msg := string(msg[:n])
		fmt.Println("$------------------>", msg, "\n")
	}
}

func main() {
	conn, err := net.Dial("tcp", ServerIP+":"+ServerPort)
	errorCheck(err, "tcp connect error")
	defer conn.Close()
	go clientSender(conn)
	go clientReceiver(conn)
	for {
		time.Sleep(1 * 1e9)
	}
}
