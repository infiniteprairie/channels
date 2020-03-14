// Package server - the package that implements the (sample) server.
//
// This package implements a sample server-side request/response
// process.
//
// The Serve() function is responsible for spawning goroutines to handle the
// client requests -- thereby emulating a server-side process. The main
// purpose of this package is to illustrate the use of channels as a message-passing
// mechanism, and the use of goroutines to do concurrent procesing.
//
// TODO: implement command-line arguments in place of the hard-coded MaxOutstanding and MaxClients constants
// 		Note: now that we've split up the client (main) and server (this package),
//		we will need to change the way in which the server uses MaxOutstanding
package server

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


// Serve takes the request channel as input (also, the "quit channel" -- currently unused)
// and dispatches groups of requests to a handler function (handle). The parameter
// maxOutstanding sets the request group size
func Serve(maxOutstanding int,clientRequests chan *Request, quit chan bool) {
	time.Sleep(1000 * time.Millisecond) // artificial delay
	// Start handler
	for i := 0; i < maxOutstanding; i++ {
		go handle(i, clientRequests)
	}
	// <-quit  // Wait to be told to exit.
}
