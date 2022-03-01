package main

import (
	"bufio"              // to implement a buffered I/O stream
	"fmt"                // to implement formatted I/O like in language C
	"log"                // to simply log
	"net"                // to provude a port-based interface for I/O like TCP
	rpc "net/rpc"        // to access exported methods of objects distributed across a network
	"os"                 // mainly used here for the input buffer, providing interface to OS functionality
	"rpc_assign/commons" // to access the commons folder in our environment namely rpc_assign
	"strconv"            // mainly to convert int to strings properly
)

/*
TODO

2) event-driven (this is the one that is done here)
	- a client starts by looking for a port to establish it's server on (like giving my phone number to my friends to call me)
		this is done via a function namely GetFreePort, defined at the start of this file.
		background:
			after searching for a bit, I found that when we set the port number to 0
			as in networkAdress:0, we end up getting a free port AUTOMATICALLY
			so I decided (with the help of some GitHub repos) to resolve an automatically
			port generated address into a TCP end point one, then return that resolved address's port
			guaranteeing that no other apps could use it.
	- a client can send a message through an infinite loop waiting for input text, this message will be broadcasted to other clients through an rpc call on the server
		this is done by maintainng a list of all connected ports represeting all connected clients
		looping in each of them, to send the new message, more like broadcasting each new message.
	- a client can also receive messages simultaneously using the GO keyword
		this is done by setting the clientListen function, we declared at the end, with a GO keyword
		in order to achieve such required concurrency.
	- so a client here is a server + a client at the same time
		Roger that.

*/

// defining the port variable, that holds the client's auto-selected free port
var clientPort int

// defining the free port selection function
func GetFreePort() (int, error) {

	// resolving the address given, which is an address that ends with a free-port
	// to an address of TCP end point.
	newAddress, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")

	if err != nil {
		return 0, err // returning 0 as the port, along with the error, if unforunately a failure occurs
	}

	// to listen for the TCP network
	inbound, err := net.ListenTCP("tcp", newAddress)

	if err != nil {
		return 0, err
	}

	defer inbound.Close() // if failure

	// returning the guaranteed free port along with nil value for the error parameter.
	return inbound.Addr().(*net.TCPAddr).Port, nil

}

func main() {

	// connecting to the RPC server located and defined as the return value of
	// the Get_server_address function that returns our main server's address
	client, err := rpc.Dial("tcp", commons.Get_server_address())

	if err != nil {
		log.Fatal(err) // logging if failure occurs
	}

	in := bufio.NewReader(os.Stdin) // preparing the input buffer

	fmt.Println("======================== Welcome =======================")
	fmt.Println("Enter a nickname to join our public chatting room: ")
	nickName, _, err := in.ReadLine() // getting the nickname from the user

	var N = string(nickName)

	// getting the free port ready, assigning it to our global variable clientPort
	clientPort, err = GetFreePort()

	// now starting the listening functionality of the client
	// with the help of GO keyword, we can either run this function on the same OS thread
	// or it can run on a different OS thread but either way concurrently.
	go clientListen()

	fmt.Println("Welcome " + N + `! You can now type whatever you want, but
be aware that everyone is reading whatever you send immediately, have fun!`)

	// defining a variable to acknowledge the registeration of the new client
	var newClientAcknowledged bool

	// sending the selected port and the ack variable (to be updated)
	err = client.Call("Listener.RegisterClient", clientPort, &newClientAcknowledged)

	if err != nil {
		log.Fatal(err)
	}

	// to point out that the user has to type something next (if he wants)
	fmt.Println("↓ Your message")

	for {
		msg, _, err := in.ReadLine() // getting the message entered by the user

		if err != nil {
			log.Fatal(err)
		}

		// creating an object of type Args, provided by the commons folder
		args := &commons.Args{Message: string(msg), Name: N}

		// defining a variable to acknowledge that the message has been sent
		var messageSentAcknowledgement bool

		// sending the message along with the client's name and the ack variable (to be updated)
		err = client.Call("Listener.SendMessage", args, &messageSentAcknowledgement)

		if err != nil {
			log.Fatal(err)
		}
	}

}

//////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////// here lies the listening part of the client /////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////

// defining the client listener class (or as a class)
type ClientListener int

// RPC ReceiveMessage method defined for the ClientListener class
func (l *ClientListener) ReceiveMessage(args *commons.Args, ack *bool) error {

	// constructing a new message that is sent by a client, broadcasted by our server
	// by combining the name of the sending client and his message
	// as in the form (name) said: "some message"
	newMessage := "(" + args.Name + ")" + " said: " + "\"" + args.Message + "\""

	// now printing it in the client's side (supposedly all connected clients)
	fmt.Println(newMessage)

	// to update the cursor's place of each client
	fmt.Println("↓ Your message")

	// acknowledging receiving the message on the client
	*ack = true

	return nil
}

func clientListen() {

	// resolving the address given to the client, to an address of TCP end point.
	addrResolved, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+strconv.Itoa(clientPort))

	if err != nil {
		log.Fatal(err)
	}

	// to listen for the TCP
	inbound, err := net.ListenTCP("tcp", addrResolved)

	if err != nil {
		log.Fatal(err)
	}

	clientListener := new(ClientListener) // allocate memory, creating an instance of our ClientListener class
	rpc.Register(clientListener)          // publishing the methods of that ClientListener into the server
	rpc.Accept(inbound)                   // accepting connections on the client listener, serving requests done probably by only the server.

}
