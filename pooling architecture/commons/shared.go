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

// didnt need the get server address function, hard coded this part with an address
// that has a port number that consists of both my fav numbers, 74 and 22 😎
