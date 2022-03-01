package main

import (
	"bufio"              // to implement a buffered I/O stream
	"fmt"                // to implement formatted I/O like in language C
	"log"                // to simply log
	rpc "net/rpc"        // to access exported methods of objects distributed across a network
	"os"                 // mainly used here for the input buffer, providing interface to OS functionality
	"rpc_assign/commons" // to access the commons folder in our environment namely rpc_assign
)

/*
TODO

1) pooling (the one that is done here)
	- the client will dial the rpc of the coordinating server (done)
	- the client will call the remote procedure on the server to send a message (done)
	- the client can fetch all of the messages history from the server using remote procedure call
	(done via fetching the whole list of messages every time a client sends some message to that specific client)

*/

func main() {

	// connecting to the RPC server located at that provided network
	client, err := rpc.Dial("tcp", "0.0.0.0:7422")

	if err != nil {
		log.Fatal(err) // logging if failure occurs
	}

	in := bufio.NewReader(os.Stdin) // preparing the input buffer

	fmt.Println("======================== Welcome =======================")
	fmt.Println("Enter a nickname to join our chatting pool: ")
	nickName, _, err := in.ReadLine() // getting the nickname from the user

	var N = string(nickName)

	fmt.Println("Welcome " + N + "! You can now type whatever you want, no restrictions in the pool, have fun!")

	for {
		fmt.Println("↓ Your message")
		msg, _, err := in.ReadLine() // getting the message entered by the user

		if err != nil {
			log.Fatal(err)
		}

		// creating an object of type Args, provided by the commons folder
		args := &commons.Args{Message: string(msg), Name: N}

		// creating an object that contains an array of strings type of variable to hold the history
		history := &commons.History{}

		// sending the message along with the client's name and the history containing object
		err = client.Call("Listener.GetLine", args, &history)

		if err != nil {
			log.Fatal(err)
		}

		// printing out the history after modifications, set up by the server
		for _, msg := range history.HistoryOfMessages {
			fmt.Println(msg)
		}
	}

}
