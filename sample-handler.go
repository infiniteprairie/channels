// Package main - the only package in this project so far.
//
// This package implements a sample client/server request/response
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
)

type Request struct {
	RequestID int
	Header    string
	Content   string
}

type Response struct {
	Header  string
	Content string
}

type ServerReqResp struct {
	SReq  *Request
	SResp *Response
}

// The process function merely emulates what the server-side processing of an
// individual request might look like
func process(gid int, payload *ServerReqResp) {
	req := payload.SReq
	fmt.Println("**processing on threadID: ", gid, " **")
	fmt.Println("requestID: ", req.RequestID)
	fmt.Println("request header/content: ", req.Header, "/ ", req.Content)
	resp := new(Response)
	resp.Content = req.Content
	payload.SResp = resp
}

// The handle function (spawned as a goroutine) takes the available requests from the channel and calls a
// process function to perform whatever work we want done on the individual requests.
func handle(gid int, queue chan *Request) {
	fmt.Println("threadID: ", gid)
	for r := range queue {
		payload := new(ServerReqResp)
		payload.SReq = r
		process(gid, payload)
		fmt.Println("request processed. Response=", payload.SResp.Content)
	}
}

// The constant, MaxOutstanding, determines the request group size for an individual
// goroutine to handle. Set MaxOutstanding to 1 if you want each thread to handle a
// single request only
const MaxOutstanding int = 3
const MaxClients int = 10

// Serve takes the request channel as input (also, the "quit channel" -- currently unused)
// and dispatches groups of requests to a handler function (handle). The constant,
// MaxOutstanding, determines the request group size
func Serve(clientRequests chan *Request, quit chan bool) {
	time.Sleep(1000 * time.Millisecond) // artificial delay
	// Start handler
	for i := 0; i < MaxOutstanding; i++ {
		go handle(i, clientRequests)
	}
	// <-quit  // Wait to be told to exit.
}

// main sets up the request channel (and the "quit channel"), then calls the Serve()
// function, which ultimately serves the client requests. Next, the (simplistic) client
// Requests are constructed add added to the request channel.
//
// Note: currently, the "quit channel" is not used. It comes from an idea in the
// Go tutorial that, in practice, would allow us to interrupt or shut down the
// request flow. It merits a further investigation in a true remote client /
// remote server model.
func main() {
	clientRequests := make(chan *Request)
	quitChannel := make(chan bool)
	Serve(clientRequests, quitChannel)

	// fire off 10 dummy Request requests
	for i := 0; i < MaxClients; i++ {
		//fmt.Println("Creating client request threads...")
		req := new(Request)
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
