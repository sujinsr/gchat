package main

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

type ClientProp struct {
	Name string
	Conn net.Conn
}

var (
	debug bool = true
)

func errorCheck(err error, errStr string) {
	if err != nil {
		fmt.Println(err, errStr)
		os.Exit(-1)
	}
}

func Log(v ...interface{}) {
	if debug == true {
		log.Print(v)
	}
}

func clientHandler(conn net.Conn, ch_msg chan string, l *list.List) {
	var stop bool = false

	/*Read name for the connection */
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	/* Convert byte array to string */
	name := string(buf[:n])

	/* Add to client list*/
	newclient := &ClientProp{name, conn}
	l.PushBack(*newclient)

	ch_msg <- name + " joined to chat"

	/* Receive messages continuously untill connection is active */
	for !stop {
		n2, err := conn.Read(buf)

		if err != nil {
			stop = true
			continue
		}
		msg := name + ":" + string(buf[:n2])
		//Log(name + " sending->" + msg)
		ch_msg <- msg
	}
	removeClient(l, name)
	fmt.Println("Closing the Client Connection")
	conn.Close()
}

func removeClient(l *list.List, name string) {
	for val := l.Front(); val != nil; val = val.Next() {
		client := val.Value.(ClientProp)
		if client.Name == name {
			l.Remove(val)
			Log("Client " + name + "removed from Client List")

		}
	}
}

func allClientSend(ch_msg chan string, l *list.List) {
	for {
		msg := <-ch_msg

		for val := l.Front(); val != nil; val = val.Next() {
			client := val.Value.(ClientProp)
			Log("send-> " + msg + " L" + strconv.Itoa(len(msg)))
			client.Conn.Write([]byte(msg))
		}
	}
}

func main() {
	client_list := list.New()
	ch_msg := make(chan string)

	netlisten, err := net.Listen("tcp", "127.0.0.1:20099")
	errorCheck(err, "Failed to listen.")
	defer netlisten.Close()

	go allClientSend(ch_msg, client_list)

	for {
		fmt.Println("Server Wait for the client to connect.")
		conn, err := netlisten.Accept()
		errorCheck(err, "Accept Failed")

		go clientHandler(conn, ch_msg, client_list)
	}

}

/*func clientReceiver(conn net.Conn, ch_msg chan string, l *list.List) {
	var stop bool = false
	message := make([]byte, 1024)

	for !stop {
		_, err := conn.Read(message)
		if err != nil {
			stop = true
			continue
		}
		ch_msg <- string(message)
	}
	fmt.Println("Closing the Client Connection")
	conn.Close()
}*/
