package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func errorCheck(err error, errStr string) {
	if err != nil {
		fmt.Println(errStr)
		os.Exit(-1)
	}
}

func clientSender(conn net.Conn) {
	fmt.Print("Enter Your Name : ")
	/*reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadBytes('\n')
	conn.Write(name[:len(name)-1])*/
	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadBytes('\n')

		//fmt.Println("I am going to send->", message)
		conn.Write(message[0 : len(message)-1])
	}
}

func clientReceiver(conn net.Conn) {
	msg := make([]byte, 1024)
	for {
		n, err := conn.Read(msg)
		errorCheck(err, "Read error")
		fmt.Println("[$", string(msg[:n]))
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:20099")
	errorCheck(err, "tcp connect error")
	defer conn.Close()
	go clientSender(conn)
	go clientReceiver(conn)
	for {
		time.Sleep(1 * 1e9)
	}
}
