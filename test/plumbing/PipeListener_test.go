//
//  PipeListener_test.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import (
	"encoding/xml"
	"github.com/puremvc/puremvc-go-util-pipes/src/messages"
	"github.com/puremvc/puremvc-go-util-pipes/src/plumbing"
	"testing"
)

/*
Test the PipeListener class.
*/

/*
  Test connecting a pipe listener to a pipe.
*/
func TestConnectingToAPipe(t *testing.T) {
	// create pipe and listener
	pipe := &plumbing.Pipe{}
	callback := Callback{}
	listener := &plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod}

	// connect the listener to the pipe
	success := pipe.Connect(listener)

	//test assertions
	if pipe == nil {
		t.Error("Expecting pipe is not nil")
	}
	if success == false {
		t.Error("Expecting successfully connected listener to pipe")
	}
}

/*
  Test receiving a message from a pipe using a PipeListener.
*/
func TestReceiveMessageViaPipeListener(t *testing.T) {
	// create a message
	messageToSend := messages.NewMessage(messages.NORMAL, &Test{testVal: 1}, []byte(`<testMessage testAtt='Hello'/>`), messages.PRIORITY_HIGH)

	// create pipe and listener
	pipe := &plumbing.Pipe{}
	callback := Callback{}
	listener := &plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod}

	// connect the listener to the pipe and write the message
	connected := pipe.Connect(listener)
	written := pipe.Write(messageToSend)

	testMessage := TestMessage{}
	xml.Unmarshal(callback.messagesReceived[0].Body().([]byte), &testMessage)

	// test assertions
	if pipe == nil {
		t.Error("Expecting pipe is not nil")
	}
	if connected != true {
		t.Error("Expecting connected listener to pipe")
	}
	if written != true {
		t.Error("Expecting wrote message to pipe")
	}
	if callback.messagesReceived == nil {
		t.Error("Expecting callback.messageReceived not nil")
	}
	if callback.messagesReceived[0].Type() != messages.NORMAL {
		t.Error("Expecting callback.messageReceived.Type() == messages.NORMAL")
	}
	if callback.messagesReceived[0].Header().(*Test).testVal != 1 {
		t.Error("Expecting callback.messageReceived.Header().(*Test).testVal == 'testval'")
	}
	if testMessage.TestAtt != "Hello" {
		t.Error("Expecting testMessage.TestAtt == 'Hello'")
	}
}
