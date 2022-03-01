package commons

/*

TODO
define any structs here to be used by the rpc
- defined a history struct to contain an array of strings to be able to point at it easily
- define an Args struct to hold the arguments needed in each message sent

*/

type History struct {
	HistoryOfMessages []string
}

type Args struct {
	Message string
	Name    string
}

// need to have the server address fixed between clients and the coordinating server
// used to define the common server address, used by all connecting client
// although this, might still be hard coded, I left it as is.

func Get_server_address() string {
	return "0.0.0.0:7422"
}
