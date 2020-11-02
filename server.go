package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	for {
		conn, _ := ln.Accept()
		conns <- conn
	}
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	reader := bufio.NewReader(client)
	for {
		msg, _ := reader.ReadString('\n')
		fmt.Println(msg)
		var newMsg Message
		newMsg.sender = clientid
		newMsg.message = msg
		msgs <- newMsg
	}
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()
	ln, _ := net.Listen("tcp", *portPtr)

	//TODO Create a Listener for TCP connections on the port given above.

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	i := 0
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			go func(conn net.Conn, clients map[int]net.Conn, msgs chan Message) {
				fmt.Println("adding")
				clients[i] = conn
				go handleClient(conn, i, msgs)
				fmt.Println("map", clients)
				fmt.Println("added")
				i++
			}(conn, clients, msgs)
			// - assign a client ID
			// - add the client to the clients channel
			// - start to asynchronously handle messages from this client
		case msg := <-msgs:
			go func(clients map[int]net.Conn, msg Message) {
				fmt.Println("sending...")
				for k := range clients {
					fmt.Println("sending to", k)
					if k != msg.sender {
						fmt.Println(msg.message)
						fmt.Fprintln(clients[k], msg.message)
					}
				}
				fmt.Println("sent")
			}(clients, msg)
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
		}
	}
}
