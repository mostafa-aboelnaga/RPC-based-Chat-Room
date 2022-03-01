package main

import (
	"fmt"                // to implement formatted I/O like in language C
	"log"                // to simply log
	"net"                // to provude a port-based interface for I/O like TCP
	rpc "net/rpc"        // to access exported methods of objects distributed across a network
	"rpc_assign/commons" // to access the commons folder in our environment namely rpc_assign
)

/*
TODO

the server has either two implementations
1) pooling (the one that is done here)
	- every message sent to the server has to be stored in long list (done)
	- a client may ask for this list or a slice of it to fetch the updates (done)
	(done via fetching the whole list of messages every time a client sends some message to that specific client)

*/

// defining the long list that holds all messages (history of messages)
var messageHistoryEntriesList []string

// defining the listener class (or as a class)
type Listener int

// RPC GetLine method defined for the Listener class
func (l *Listener) GetLine(args *commons.Args, history *commons.History) error {

	// constructing a new entry line in our history of messages list
	// by combining the name of the client and message itself
	// as in the form (name) said: "some message"
	newEntryLine := "(" + args.Name + ")" + " said: " + "\"" + args.Message + "\""

	// appending the new entry line into our list
	messageHistoryEntriesList = append(messageHistoryEntriesList, newEntryLine)

	// updating the history that the user/client sees on his side with the latest
	// entries sent by all connected users in the pool
	(*history).HistoryOfMessages = messageHistoryEntriesList

	// also printing the list ourselves in the server side, to monitor the chatting pool
	for _, entry := range messageHistoryEntriesList {
		fmt.Println(entry)
	}

	return nil
}

func main() {

	// resolving the address given, to an address of TCP end point.
	addrResolved, err := net.ResolveTCPAddr("tcp", "0.0.0.0:7422")

	if err != nil {
		log.Fatal(err) // logging if failure
	}

	// to listen for the TCP
	inbound, err := net.ListenTCP("tcp", addrResolved)

	if err != nil {
		log.Fatal(err)
	}

	serverListener := new(Listener) // allocate memory, creating an instance of our Listener class
	rpc.Register(serverListener)    // publishing the methods of that Listener into the server
	rpc.Accept(inbound)             // accepting connections on the listener, serving requests on the server, etc.

}
