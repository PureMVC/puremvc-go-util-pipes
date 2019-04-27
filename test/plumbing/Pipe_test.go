//
//  Pipe_test.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import (
	"github.com/puremvc/puremvc-go-util-pipes/src/interfaces"
	"github.com/puremvc/puremvc-go-util-pipes/src/plumbing"
	"testing"
)

/*
Test the Pipe class.
*/

/*
 * Test the constructor.
 */
func TestConstructor(t *testing.T) {
	var pipe interfaces.IPipeFitting = &plumbing.Pipe{}

	// test assertions
	if pipe == nil {
		t.Error("Expecting pipe not nil")
	}
}

/*
  Test connecting and disconnecting two pipes.
*/
func TestConnectingAndDisconnectingTwoPipes(t *testing.T) {
	// create two pipes
	var pipe1 interfaces.IPipeFitting = &plumbing.Pipe{}
	var pipe2 interfaces.IPipeFitting = &plumbing.Pipe{}

	// connect them
	success := pipe1.Connect(pipe2)

	// test assertions
	if pipe1 == nil {
		t.Error("Expecting pipe1 not nil")
	}
	if pipe2 == nil {
		t.Error("Expecting pipe2 not nil")
	}
	if success != true {
		t.Error("Expecting connected pipe1 to pipe2")
	}

	// disconnect pipe 2 from pipe 1
	disconnectedPipe := pipe1.Disconnect()
	if disconnectedPipe != pipe2 {
		t.Error("Expecing disconnected pipe2 from pipe1")
	}
}

/*
  Test attempting to connect a pipe to a pipe with an output already connected.
*/
func TestConnectingToAConnectedPipe(t *testing.T) {
	// create three pipes
	var pipe1 interfaces.IPipeFitting = &plumbing.Pipe{}
	var pipe2 interfaces.IPipeFitting = &plumbing.Pipe{}
	var pipe3 interfaces.IPipeFitting = &plumbing.Pipe{}

	// connect them
	success := pipe1.Connect(pipe2)

	// test assertions
	if success != true {
		t.Error("Expecing connected pipe1 to pipe2")
	}
	if pipe1.Connect(pipe3) != false {
		t.Error("Expecting can't connect pipe3 to pipe1")
	}
}
