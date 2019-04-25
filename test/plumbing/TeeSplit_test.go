//
//  TeeSplit_test.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import (
	"github.com/puremvc/puremvc-go-util-pipes/src/interfaces"
	"github.com/puremvc/puremvc-go-util-pipes/src/messages"
	"github.com/puremvc/puremvc-go-util-pipes/src/plumbing"
	"testing"
)

/**
Test the TeeSplit class.
*/

/**
  Test connecting and disconnecting I/O Pipes.

  Connect an input and several output pipes to a splitting tee.
  Then disconnect all outputs in LIFO order by calling disconnect
  repeatedly.
*/
func TestConnectingAndDisconnectingIOPipes(t *testing.T) {
	// create input pipe
	input1 := plumbing.Pipe{}

	// create output pipes 1, 2, 3 and 4
	pipe1 := plumbing.Pipe{}
	pipe2 := plumbing.Pipe{}
	pipe3 := plumbing.Pipe{}
	pipe4 := plumbing.Pipe{}

	// create splitting tee (args are first two output fittings of tee)
	teeSplit := plumbing.TeeSplit{}

	// connect 2 extra outputs for a total of 4
	connected1 := teeSplit.Connect(&pipe1)
	connected2 := teeSplit.Connect(&pipe2)
	connected3 := teeSplit.Connect(&pipe3)
	connected4 := teeSplit.Connect(&pipe4)

	// connect the single input
	inputConnected := input1.Connect(&teeSplit)

	// test assertions
	if connected1 != true {
		t.Error("Expecting connected pipe 1")
	}
	if connected2 != true {
		t.Error("Expecting connected pipe 2")
	}
	if connected3 != true {
		t.Error("Expecting connected pipe 3")
	}
	if connected4 != true {
		t.Error("Expecting connected pipe 4")
	}
	if inputConnected != true {
		t.Error("Expecting connected input1 to teeSplit")
	}

	//test LIFO order of output disconnection
	if teeSplit.Disconnect() != &pipe4 {
		t.Error("Expecting disconnected pipe 4")
	}
	if teeSplit.Disconnect() != &pipe3 {
		t.Error("Expecting disconnected pipe 3")
	}
	if teeSplit.Disconnect() != &pipe2 {
		t.Error("Expecting disconnected pipe 2")
	}
	if teeSplit.Disconnect() != &pipe1 {
		t.Error("Expecting disconnected pipe 1")
	}
}

/**
  Test disconnectFitting method.

  Connect several output pipes to a splitting tee.
  Then disconnect specific outputs, making sure that once
  a fitting is disconnected using disconnectFitting, that
  it isn't returned when disconnectFitting is called again.
  Finally, make sure that the when a message is sent to
  the tee that the correct number of output messages is
  written.
*/
func TestDisconnectFitting(t *testing.T) {
	callback := Callback{}

	// create output pipes 1, 2, 3 and 4
	pipe1 := plumbing.Pipe{}
	pipe2 := plumbing.Pipe{}
	pipe3 := plumbing.Pipe{}
	pipe4 := plumbing.Pipe{}

	// setup pipelisteners
	pipe1.Connect(&plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod})
	pipe2.Connect(&plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod})
	pipe3.Connect(&plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod})
	pipe4.Connect(&plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod})

	// create splitting tee
	teeSplit := plumbing.TeeSplit{}

	// add outputs
	if !teeSplit.Connect(&pipe1) {
		t.Error("Expecting pipe1 connection")
	}
	if !teeSplit.Connect(&pipe2) {
		t.Error("Expecting pipe2 connection")
	}
	if !teeSplit.Connect(&pipe3) {
		t.Error("Expecting pipe3 connection")
	}
	if !teeSplit.Connect(&pipe4) {
		t.Error("Expecting pipe4 connection")
	}

	// test assertions
	if teeSplit.DisconnectFitting(&pipe4) != &pipe4 {
		t.Error("Expecting teeSplit.disconnectFitting(&pipe4) == &pipe4")
	}

	// Write a message to the tee
	teeSplit.Write(messages.NewMessage(messages.NORMAL, nil, nil, messages.PRIORITY_MED))
	if len(callback.messagesReceived) != 3 {
		t.Error("Expecting messagesReceived.count == 3")
	}

	if teeSplit.DisconnectFitting(&pipe3) != &pipe3 {
		t.Error("Expecting teeSplit.disconnectFitting(&pipe3) == &pipe3")
	}
	if teeSplit.DisconnectFitting(&pipe2) != &pipe2 {
		t.Error("Expecting teeSplit.disconnectFitting(&pipe2) == &pipe2")
	}
	if teeSplit.DisconnectFitting(&pipe1) != &pipe1 {
		t.Error("Expecting teeSplit.disconnectFitting(&pipe1) == &pipe1")
	}
	if teeSplit.DisconnectFitting(&pipe4) != nil {
		t.Error("Expecting teeSplit.disconnectFitting(&pipe4) == nil")
	}
}

/**
  Test receiving messages from two pipes using a TeeMerge.
*/
func TestReceiveMessagesFromTwoTeeSplitOutputs(t *testing.T) {
	callback := Callback{}

	// create a message to send on pipe 1
	message := messages.NewMessage(messages.NORMAL, Test{testVal: 1}, "", messages.PRIORITY_MED)

	// create output pipes 1 and 2
	pipe1 := plumbing.Pipe{}
	pipe2 := plumbing.Pipe{}

	// create and connect anonymous listeners
	connected1 := pipe1.Connect(&plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod})
	connected2 := pipe2.Connect(&plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod})

	// create splitting tee (args are first two output fittings of tee)
	teeSplit := plumbing.TeeSplit{}
	teeSplit.Connect(&pipe1)
	teeSplit.Connect(&pipe2)

	// write messages to their respective pipes
	written := teeSplit.Write(message)

	// test assertions
	if message == nil {
		t.Error("Expecting message not nil")
	}
	if &pipe1 == nil {
		t.Error("Expecting pipe1 not nil")
	}
	if &pipe2 == nil {
		t.Error("Expecting pipe2 not nil")
	}
	if &teeSplit == nil {
		t.Error("Expecting teeSplit is not nil")
	}
	if connected1 != true {
		t.Error("Expecting connected anonymous listener to pipe 1")
	}
	if connected2 != true {
		t.Error("Expecting connected anonymous listener to pipe 2")
	}
	if written != true {
		t.Error("Expecting wrote single message to tee")
	}

	// test that both messages were received, then test
	// FIFO order by inspecting the messages themselves
	if len(callback.messagesReceived) != 2 {
		t.Error("Expecting received 2 messages")
	}

	// test message 1 assertions
	var message1 interfaces.IPipeMessage
	message1, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:]
	if message1 == nil {
		t.Error("Expecting message1 not nil")
	}
	if message1 != message {
		t.Error("Expecting message1 == message")
	}

	// test message 2 assertions
	var message2 interfaces.IPipeMessage
	message2, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:]
	if message2 == nil {
		t.Error("Expecting message2 not nil")
	}
	if message2 != message {
		t.Error("Expecting message2 == message")
	}
}
