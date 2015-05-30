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

type Message struct {
	Type int
	Name string
	Data string
}

const (
	CONTROL_MSG = 1
	BROAD_MSG   = 2
	CHAT_MSG    = 3
)

var (
	ServerIP   string = "127.0.0.1"
	ServerPort string = "20099"
	debug      bool   = true
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

func clientHandler(conn net.Conn, ch_msg chan Message, l *list.List) {
	var stop bool = false
	msg := Message{}

	/*Read name for the connection from client*/
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	/* Convert byte array to string */
	name := string(buf[:n])
	/* check same name already available in list */
	avail := checkAvail(l, name)
	if avail == true {
		Log("Client Name already available in the list")
		conn.Write([]byte(strconv.Itoa(1)))
		return
	}
	/*Send success control message to client*/
	conn.Write([]byte(strconv.Itoa(0)))

	/* Add the client property to property list for future use */
	newclient := &ClientProp{name, conn}
	l.PushBack(*newclient)

	/* Send client add message to all connected clients */
	msg.Type = BROAD_MSG
	msg.Name = name
	msg.Data = "joined the chat"
	ch_msg <- msg
	Log("Client " + name + " Connected to the Sever")

	/* Receive messages continuously untill connection is active */
	for !stop {
		n2, err := conn.Read(buf)

		if err != nil {
			stop = true
			continue
		}
		msg.Type = CHAT_MSG
		msg.Name = name
		msg.Data = string(buf[:n2])
		//Log(name + " sending->" + msg)
		ch_msg <- msg
	}

	/* Remove the client and send status messge */
	removeClient(l, name)

	msg.Type = BROAD_MSG
	msg.Data = "left the chat"
	ch_msg <- msg

	fmt.Println("Client", name, "Closed the Connection")
	conn.Close()
}

func allClientSend(ch_msg chan Message, l *list.List) {
	var write_msg string

	for {
		/* receive the data from channel */
		msg := <-ch_msg

		/* Format the message based on the type to sent */
		if msg.Type == BROAD_MSG {
			write_msg = msg.Name + " " + msg.Data
		} else if msg.Type == CHAT_MSG {
			write_msg = "[ " + msg.Name + "-> " + msg.Data + " ]"
		}

		Log("send-> " + write_msg)

		/* Send the message all connected client except sender */
		for val := l.Front(); val != nil; val = val.Next() {
			client := val.Value.(ClientProp)
			if msg.Name != client.Name {
				client.Conn.Write([]byte(write_msg))
			}
		}
	}
}

func checkAvail(l *list.List, name string) bool {
	for val := l.Front(); val != nil; val = val.Next() {
		client := val.Value.(ClientProp)
		if client.Name == name {
			return true
		}
	}
	return false
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

func main() {
	client_list := list.New()
	ch_msg := make(chan Message)

	netlisten, err := net.Listen("tcp", ServerIP+":"+ServerPort)
	errorCheck(err, "Failed to listen.")
	defer netlisten.Close()

	go allClientSend(ch_msg, client_list)

	fmt.Println("Server Wait for the client to connect.")
	for {
		conn, err := netlisten.Accept()
		errorCheck(err, "Accept Failed")

		go clientHandler(conn, ch_msg, client_list)
	}

}
