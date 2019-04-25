//
//  TeeMerge_test.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import (
	"encoding/xml"
	"github.com/puremvc/puremvc-go-util-pipes/src/interfaces"
	"github.com/puremvc/puremvc-go-util-pipes/src/messages"
	"github.com/puremvc/puremvc-go-util-pipes/src/plumbing"
	"testing"
)

/**
Test the TeeMerge class.
*/

/**
  Test connecting an output and several input pipes to a merging tee.
*/
func TestConnectingIOPipes(t *testing.T) {
	// create input pipe
	output1 := &plumbing.Pipe{}

	// create input pipes 1, 2, 3 and 4
	pipe1 := &plumbing.Pipe{}
	pipe2 := &plumbing.Pipe{}
	pipe3 := &plumbing.Pipe{}
	pipe4 := &plumbing.Pipe{}

	// create splitting tee (args are first two input fittings of tee)
	teeMerge := &plumbing.TeeMerge{}

	// connect 2 extra inputs for a total of 4
	connectedExtra1 := teeMerge.ConnectInput(pipe1)
	connectedExtra2 := teeMerge.ConnectInput(pipe2)
	connectedExtra3 := teeMerge.ConnectInput(pipe3)
	connectedExtra4 := teeMerge.ConnectInput(pipe4)

	// connect the single output
	connected := output1.Connect(teeMerge)

	// test assertions
	if pipe1 == nil {
		t.Error("Expecting pipe1 is not nil")
	}
	if pipe2 == nil {
		t.Error("Expecting pipe2 is not nil")
	}
	if pipe3 == nil {
		t.Error("Expecting pipe3 is not nil")
	}
	if pipe4 == nil {
		t.Error("Expecting pipe4 is not nil")
	}
	if teeMerge == nil {
		t.Error("Expecting teeMerge is TeeMerge")
	}
	if connected == false {
		t.Error("Expecting connected")
	}
	if connectedExtra1 != true {
		t.Error("Expecting connected extra input 1")
	}
	if connectedExtra2 != true {
		t.Error("Expecting connected extra input 2")
	}
	if connectedExtra3 != true {
		t.Error("Expecting connected extra input 3")
	}
	if connectedExtra4 != true {
		t.Error("Expecting connected extra input 4")
	}
}

/**
  Test receiving messages from two pipes using a TeeMerge.
*/
func TestReceiveMessagesFromTwoPipesViaTeeMerge(t *testing.T) {
	// create a message to send on pipe 1
	pipe1Message := messages.NewMessage(messages.NORMAL, Test{testVal: 1}, []byte(`<testMessage testAtt='Pipe1Message'/>`), messages.PRIORITY_LOW)

	// create a message to send on pipe 2
	pipe2Message := messages.NewMessage(messages.NORMAL, Test{testVal: 2}, []byte(`<testMessage testAtt='Pipe2Message'/>`), messages.PRIORITY_HIGH)

	// create pipes 1 and 2
	pipe1 := &plumbing.Pipe{}
	pipe2 := &plumbing.Pipe{}

	// create merging tee (args are first two input fittings of tee)
	teeMerge := &plumbing.TeeMerge{}
	teeMerge.ConnectInput(pipe1)
	teeMerge.ConnectInput(pipe2)

	// create listener
	callback := &Callback{}
	listener := &plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod}

	// connect the listener to the tee and write the messages
	connected := teeMerge.Connect(listener)

	// write messages to their respective pipes
	pipe1Written := pipe1.Write(pipe1Message)
	pipe2Written := pipe2.Write(pipe2Message)

	// test assertions
	if pipe1Message == nil {
		t.Error("Expecting pipe1Message is not nil")
	}
	if pipe2Message == nil {
		t.Error("Expecting pipe2Message is not nil")
	}
	if pipe1 == nil {
		t.Error("Expecting pipe1 not nil")
	}
	if pipe2 == nil {
		t.Error("Expecting pipe2 not nil")
	}
	if teeMerge == nil {
		t.Error("Expecting teeMerge not nil")
	}
	if listener == nil {
		t.Error("Expecting listener not nil")
	}
	if connected != true {
		t.Error("Expecting listener connected to teeMerge")
	}
	if pipe1Written != true {
		t.Error("Expecting wrote message to pipe 1")
	}
	if pipe2Written != true {
		t.Error("Expecting wrote message to pipe 2")
	}

	// test that both messages were received, then test
	// FIFO order by inspecting the messages themselves
	if len(callback.messagesReceived) != 2 {
		t.Error("Expecting received 2 messages")
	}

	// test message 1 assertions
	var message1 interfaces.IPipeMessage
	message1, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:] // shift
	testMessage1 := TestMessage{}
	xml.Unmarshal(message1.Body().([]byte), &testMessage1)

	if message1 == nil {
		t.Error("Expecting message1 not nil")
	}
	if message1 != pipe1Message {
		t.Error("Expecting message1 == pipe1Message")
	}
	if message1.Type() != messages.NORMAL {
		t.Error("Expecting message1.Type() == messages.NORMAL")
	}
	if message1.Header().(Test).testVal != 1 {
		t.Error("Expecting message1.Header().(Test).testVal != '1'")
	}
	if testMessage1.TestAtt != "Pipe1Message" {
		t.Error("Expecting testMessage1.TestAtt != 'Pipe1Message'")
	}
	if message1.Priority() != messages.PRIORITY_LOW {
		t.Error("Expecting message1.Priority() == messages.PRIORITY_LOW")
	}

	// test message 2 assertions
	var message2 interfaces.IPipeMessage
	message2, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:]
	testMessage2 := TestMessage{}
	xml.Unmarshal(message2.Body().([]byte), &testMessage2)

	if message2 == nil {
		t.Error("Expecting message2 not nil")
	}
	if message2 != pipe2Message {
		t.Error("Expecting message2 == pipe1Message")
	}
	if message2.Type() != messages.NORMAL {
		t.Error("Expecting message2.Type() == messages.NORMAL")
	}
	if message2.Header().(Test).testVal != 2 {
		t.Error("Expecting message2.Header().(Test).testVal != '2'")
	}
	if testMessage2.TestAtt != "Pipe2Message" {
		t.Error("Expecting testMessage2.TestAtt != 'Pipe1Message'")
	}
	if message2.Priority() != messages.PRIORITY_HIGH {
		t.Error("Expecting message2.Priority() == messages.PRIORITY_LOW")
	}
}

/**
  Test receiving messages from four pipes using a TeeMerge.
*/
func TestReceiveMessagesFromFourPipesViaTeeMerge(t *testing.T) {
	// create a message to send on pipe 1
	pipe1Message := messages.NewMessage(messages.NORMAL, Test{testVal: 1}, nil, messages.PRIORITY_MED)
	pipe2Message := messages.NewMessage(messages.NORMAL, Test{testVal: 2}, nil, messages.PRIORITY_MED)
	pipe3Message := messages.NewMessage(messages.NORMAL, Test{testVal: 3}, nil, messages.PRIORITY_MED)
	pipe4Message := messages.NewMessage(messages.NORMAL, Test{testVal: 4}, nil, messages.PRIORITY_MED)

	// create pipes 1, 2, 3 and 4
	pipe1 := &plumbing.Pipe{}
	pipe2 := &plumbing.Pipe{}
	pipe3 := &plumbing.Pipe{}
	pipe4 := &plumbing.Pipe{}

	// create merging tee
	teeMerge := &plumbing.TeeMerge{}
	connectedExtraInput1 := teeMerge.ConnectInput(pipe1)
	connectedExtraInput2 := teeMerge.ConnectInput(pipe2)
	connectedExtraInput3 := teeMerge.ConnectInput(pipe3)
	connectedExtraInput4 := teeMerge.ConnectInput(pipe4)

	// create listener
	callback := Callback{}
	listener := &plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod}

	// connect the listener to the tee and write the messages
	connected := teeMerge.Connect(listener)

	// write messages to their respective pipes
	pipe1Written := pipe1.Write(pipe1Message)
	pipe2Written := pipe2.Write(pipe2Message)
	pipe3Written := pipe3.Write(pipe3Message)
	pipe4Written := pipe4.Write(pipe4Message)

	// test assertions
	if pipe1Message == nil {
		t.Error("Expecting pipe1Message not nil")
	}
	if pipe2Message == nil {
		t.Error("Expecting pipe2Message not nil")
	}
	if pipe3Message == nil {
		t.Error("Expecting pipe3Message not nil")
	}
	if pipe4Message == nil {
		t.Error("Expecting pipe4Message not nil")
	}
	if pipe1 == nil {
		t.Error("Expecting pipe1 not nil")
	}
	if pipe2 == nil {
		t.Error("Expecting pipe2 not nil")
	}
	if pipe3 == nil {
		t.Error("Expecting pipe3 not nil")
	}
	if pipe4 == nil {
		t.Error("Expecting pipe4 not nil")
	}
	if teeMerge == nil {
		t.Error("Expecting teeMerge not nil")
	}
	if connected != true {
		t.Error("Expecting connected listener to merging tee")
	}
	if connectedExtraInput1 != true {
		t.Error("Expecting connected extra input pipe1 to merging tee")
	}
	if connectedExtraInput2 != true {
		t.Error("Expecting connected extra input pipe2 to merging tee")
	}
	if connectedExtraInput3 != true {
		t.Error("Expecting connected extra input pipe3 to merging tee")
	}
	if connectedExtraInput4 != true {
		t.Error("Expecting connected extra input pipe4 to merging tee")
	}
	if pipe1Written != true {
		t.Error("Expecting wrote message to pipe 1")
	}
	if pipe2Written != true {
		t.Error("Expecting wrote message to pipe 2")
	}
	if pipe3Written != true {
		t.Error("Expecting wrote message to pipe 3")
	}
	if pipe4Written != true {
		t.Error("Expecting wrote message to pipe 4")
	}

	// test that both messages were received, then test
	// FIFO order by inspecting the messages themselves
	if len(callback.messagesReceived) != 4 {
		t.Error("Expecting received 4 messages")
	}

	// test message 1 assertions
	var message1 interfaces.IPipeMessage
	message1, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:]
	if message1 == nil {
		t.Error("Expecting message1 is IPipeMessage")
	}
	if message1 != pipe1Message {
		t.Error("Expecting message1 === pipe1Message")
	}
	if message1.Type() != messages.NORMAL {
		t.Error("Expecting message1.Type() == messages.NORMAL")
	}
	if message1.Header().(Test).testVal != 1 {
		t.Error("Expecting message1.Header().(*Test).testVal != 1")
	}

	// test message 2 assertions
	var message2 interfaces.IPipeMessage
	message2, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:]
	if message2 == nil {
		t.Error("Expecting message2 is IPipeMessage")
	}
	if message2 != pipe2Message {
		t.Error("Expecting message2 === pipe2Message")
	}
	if message2.Type() != messages.NORMAL {
		t.Error("Expecting message2.Type() == messages.NORMAL")
	}
	if message2.Header().(Test).testVal != 2 {
		t.Error("Expecting message2.Header().(*Test).testVal != 2")
	}

	// test message 3 assertions
	var message3 interfaces.IPipeMessage
	message3, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:]
	if message3 == nil {
		t.Error("Expecting message3 is IPipeMessage")
	}
	if message3 != pipe3Message {
		t.Error("Expecting message3 === pipe2Message")
	}
	if message3.Type() != messages.NORMAL {
		t.Error("Expecting message3.Type() == messages.NORMAL")
	}
	if message3.Header().(Test).testVal != 3 {
		t.Error("Expecting message3.Header().(*Test).testVal != 3")
	}

	// test message 4 assertions
	var message4 interfaces.IPipeMessage
	message4, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:]
	if message4 == nil {
		t.Error("Expecting message4 is IPipeMessage")
	}
	if message4 != pipe4Message {
		t.Error("Expecting message4 === pipe4Message")
	}
	if message4.Type() != messages.NORMAL {
		t.Error("Expecting message4.Type() == messages.NORMAL")
	}
	if message4.Header().(Test).testVal != 4 {
		t.Error("Expecting message4.Header().(*Test).testVal != 4")
	}
}
