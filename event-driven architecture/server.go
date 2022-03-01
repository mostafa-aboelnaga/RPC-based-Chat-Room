package main

import (
	"fmt"                // to implement formatted I/O like in language C
	"log"                // to simply log
	"net"                // to provude a port-based interface for I/O like TCP
	rpc "net/rpc"        // to access exported methods of objects distributed across a network
	"rpc_assign/commons" // to access the commons folder in our environment namely rpc_assign
	"strconv"            // mainly to convert int to strings properly
)

/*

TODO

2) event-driven (this is the one that is done here)
	- a server is more like a coordinator
		Yup, it does what a normal server is expected to do, listens.
	- the server waits for clients wanting to register themselves as listeners
		this is done by the RegisterClient method, defined in the listener of our server.
	- a client sends a message by calling an rpc responsible for broadcasting, the client calls a function to loop on all registered clients and send his own message to each of them separately
		this is done by the SendMessage method, that loops on all registered ports, and simply
		broadcasts the message sent by the client to all of them.
	- the client on the other side has a server listening for messages being pushed
		Yup.

*/

// defining the list that holds port numbers of each registered client
var registeredPorts []int

// defining the listener class (or as a class)
type Listener int

// RPC SendMessage method defined for the Listener class
func (l *Listener) SendMessage(args *commons.Args, ack *bool) error {

	// here, we loop on the registered ports, forwarding the message to each client that is registered
	// in our server, using the client's own receive message method provided by the ClientListener
	// named ReceiveMessage

	for _, port := range registeredPorts {

		// connecting to the address constructed by the network address and each registered client's port
		registeredClient, err := rpc.Dial("tcp", "0.0.0.0:"+strconv.Itoa(port))

		if err != nil {
			log.Fatal(err) // logging if failure occurs
		}

		// defining a variable to acknowledge that the message has been received by this registered client
		var messageReceivedAcknowledgement bool

		// sending the arguments of the served message (to be broadcasted) and the ack variable (to be updated)
		err = registeredClient.Call("ClientListener.ReceiveMessage", args, &messageReceivedAcknowledgement)

		if err != nil {
			log.Fatal(err)
		}
	}

	// also printing the message ourselves in the server side, to monitor the chatting room
	// by constructing the messaage in the same manner
	// as in the form (name) said: "some message"
	newMessage := "(" + args.Name + ")" + " said: " + "\"" + args.Message + "\""
	fmt.Println(newMessage)

	// acknowledging receiving the message on the client
	*ack = true

	return nil
}

// RPC RegisterClient method defined for the Listener class
func (l *Listener) RegisterClient(clientPortNumber int, ack *bool) error {

	// we append the selected port number to our list of registered ports (to be registered)
	registeredPorts = append(registeredPorts, clientPortNumber)

	// send out a notification in the server itself that a new client has been registered
	// with its specific port number, mainly for monitoring purposes
	fmt.Println("A new client with port " + strconv.Itoa(clientPortNumber) + " has been registered")

	// acknowledging the registeration of the client
	*ack = true

	return nil
}

func main() {

	// resolving the address of our RPC server located and defined as the return value of
	// the Get_server_address function that returns our main server's address
	// to an address of TCP end point.
	addrResolved, err := net.ResolveTCPAddr("tcp", commons.Get_server_address())

	if err != nil {
		log.Fatal(err)
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
