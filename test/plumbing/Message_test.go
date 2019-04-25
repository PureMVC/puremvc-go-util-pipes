//
//  Message_test.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import (
	"encoding/xml"
	"github.com/puremvc/puremvc-go-util-pipes/src/messages"
	"testing"
)

/**
Test the Message class.
*/

/**
  Tests the constructor parameters and getters.
*/
func TestConstructorAndGetters(t *testing.T) {
	var message = messages.NewMessage(messages.NORMAL, &Test{testVal: 1}, []byte(`<testMessage testAtt='Hello'/>`), messages.PRIORITY_HIGH)

	// create a message with complete constructor args
	var msg = TestMessage{}
	xml.Unmarshal(message.Body().([]byte), &msg)

	// test assertions
	if message.Type() != messages.NORMAL {
		t.Error("Expecting message.getType() == messages.NORMAL")
	}
	if message.Header().(*Test).testVal != 1 {
		t.Error("Expecting message.Header().testProp == 'testval'")
	}
	if msg.TestAtt != "Hello" {
		t.Error("Expecting message.TestAtt == 'Hello'")
	}
	if message.Priority() != messages.PRIORITY_HIGH {
		t.Error("Expecting message.Priority() == messages.MESSAGE_PRIORITY_HIGH")
	}
}

/**
  Tests the setters and getters.
*/
func TestSettersAndGetters(t *testing.T) {
	message := messages.NewMessage(messages.NORMAL, nil, nil, messages.PRIORITY_MED)
	message.SetHeader(&Test{testVal: 1})
	message.SetBody([]byte(`<testMessage testAtt='Hello'/>`))
	message.SetPriority(messages.PRIORITY_LOW)

	var msg = TestMessage{}
	xml.Unmarshal(message.Body().([]byte), &msg)

	if message.Type() != messages.NORMAL {
		t.Error("Expecting message.getType() == messages.NORMAL")
	}
	if message.Header().(*Test).testVal != 1 {
		t.Error("Expecting message.Header().testProp == 'testval'")
	}
	if msg.TestAtt != "Hello" {
		t.Error("Expecting message.TestAtt == 'Hello'")
	}
	if message.Priority() != messages.PRIORITY_LOW {
		t.Error("Expecting message.Priority() == messages.PRIORITY_LOW")
	}
}
