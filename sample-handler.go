// Package main
//
// This package implements the client side of a sample client/server request/response
// process. The client requests are created in the main() function.
// The Serve() function is responsible for spawning goroutines to handle the
// client requests -- thereby emulating a server-side process. The main
// purpose of this package is to illustrate the use of channels as a message-passing
// mechanism, and the use of goroutines to do concurrent procesing.
//
// TODO: separate out the server portion as a separate Package
// TODO: implement command-line arguments in place of the hard-coded MaxOutstanding and MaxClients constants
package main

import (
	"fmt"
	"time"
	"github.com/infiniteprairie/channels/server"
)


// The constant, MaxOutstanding, determines the request group size for an individual
// goroutine to handle. Set MaxOutstanding to 1 if you want each thread to handle a
// single request only
const MaxOutstanding int = 3
const MaxClients int = 10


// main sets up the request channel (and the "quit channel"), then calls the Serve()
// function, which ultimately serves the client requests. Next, the (simplistic) client
// Requests are constructed add added to the request channel.
//
// Note: currently, the "quit channel" is not used. It comes from an idea in the
// Go tutorial that, in practice, would allow us to interrupt or shut down the
// request flow. It merits a further investigation in a true remote client /
// remote server model.
func main() {
	clientRequests := make(chan *server.Request)
	quitChannel := make(chan bool)
	server.Serve(MaxOutstanding, clientRequests, quitChannel)

	// fire off 10 dummy Request requests
	for i := 0; i < MaxClients; i++ {
		//fmt.Println("Creating client request threads...")
		req := new(server.Request)
		req.RequestID = i
		req.Header = fmt.Sprintf("Request #%v", i)
		req.Content = fmt.Sprintf("Content: this is the content for Request #%v)", i)
		clientRequests <- req
		// impose a delay, but only every 3rd request cycle
		if i%3 == 0 {
			fmt.Println("Client req ", i, "created ... Sleeping for 2 seconds before the next batch...")
			time.Sleep(2000 * time.Millisecond)
		}
		// serverResponse
	}
	fmt.Println("Exited the request loop...")
	fmt.Println("Sleeping for 5 seconds...")
	time.Sleep(5000 * time.Millisecond)
	//quitChannel <- true

}
