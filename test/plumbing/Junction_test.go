//
//  Junction_test.go
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
Test the Junction class.
*/

/**
  Test registering an INPUT pipe to a junction.

  Tests that the INPUT pipe is successfully registered and
  that the hasPipe and hasInputPipe methods work. Then tests
  that the pipe can be retrieved by name.

  Finally, it removes the registered INPUT pipe and tests
  that all the previous assertions about it's registration
  and accessability via the Junction are no longer true.
*/
func TestRegisterRetrieveAndRemoveInputPipe(t *testing.T) {
	// create pipe connected to this test with a pipelistener
	pipe := &plumbing.Pipe{}

	// create junction
	junction := &plumbing.Junction{PipesMap: map[string]interfaces.IPipeFitting{}, PipeTypesMap: map[string]string{}}

	// register the pipe with the junction, giving it a name and direction
	registered := junction.RegisterPipe("testInputPipe", plumbing.INPUT, pipe)

	// test assertions
	if pipe == nil {
		t.Error("Expecting pipe not nil")
	}
	if junction == nil {
		t.Error("Expecting junction not nil")
	}
	if registered != true {
		t.Error("Expecting success registering pipe")
	}

	// assertions about junction methods once input  pipe is registered
	if junction.HasPipe("testInputPipe") != true {
		t.Error("Expecting junction has pipe")
	}
	if junction.HasInputPipe("testInputPipe") != true {
		t.Error("Expecting junction has pipe registered as an INPUT type")
	}
	if junction.RetrievePipe("testInputPipe") != pipe {
		t.Error("Expecting pipe retrieved from junction")
	}

	// now remove the pipe and be sure that it is no longer there (same assertions should be false)
	junction.RemovePipe("testInputPipe")
	if junction.HasPipe("testInputPipe") != false {
		t.Error("Expecting junction has no pipe")
	}
	if junction.HasInputPipe("testInputPipe") != false {
		t.Error("Expecting junction has no input pipe")
	}
	if junction.RetrievePipe("testInputPipe") != nil {
		t.Error("Expecting nil retrieved from junction")
	}
}

/**
  Test registering an OUTPUT pipe to a junction.

  Tests that the OUTPUT pipe is successfully registered and
  that the hasPipe and hasOutputPipe methods work. Then tests
  that the pipe can be retrieved by name.

  Finally, it removes the registered OUTPUT pipe and tests
  that all the previous assertions about it's registration
  and accessability via the Junction are no longer true.
*/
func TestRegisterRetrieveAndRemoveOutputPipe(t *testing.T) {
	// create pipe connected to this test with a pipelistener
	pipe := &plumbing.Pipe{}

	// create junction
	junction := &plumbing.Junction{PipesMap: make(map[string]interfaces.IPipeFitting), PipeTypesMap: make(map[string]string)}

	// register the pipe with the junction, giving it a name and direction
	registered := junction.RegisterPipe("testOutputPipe", plumbing.OUTPUT, pipe)

	// test assertions
	if pipe == nil {
		t.Error("Expecting pipe is Pipe")
	}
	if junction == nil {
		t.Error("Expecting junction is Junction")
	}
	if registered != true {
		t.Error("Expecting success registering pipe")
	}

	// assertions about junction methods once output pipe is registered
	if junction.HasPipe("testOutputPipe") != true {
		t.Error("Expecting junction has pipe")
	}
	if junction.HasOutputPipe("testOutputPipe") != true {
		t.Error("Expecting junction has pipe registered as an OUTPUT pipe")
	}
	if junction.RetrievePipe("testOutputPipe") != pipe {
		t.Error("Expecting pipe retrieved from junction")
	}

	// now remove the pipe and be sure that it is no longer there (same assertions should be false)
	junction.RemovePipe("testOutputPipe")
	if junction.HasPipe("testOutputPipe") != false {
		t.Error("Expecting junction has no pipe")
	}
	if junction.HasOutputPipe("testOutputPipe") != false {
		t.Error("Expecting junction has pipe not registered as an OUTPUT pipe")
	}
	if junction.RetrievePipe("testOutputPipe") != nil {
		t.Error("Expecting nil retrieved from junction")
	}
}

/**
  Test adding a PipeListener to an Input Pipe.

  Registers an INPUT Pipe with a Junction, then tests
  the Junction's addPipeListener method, connecting
  the output of the pipe back into to the test. If this
  is successful, it sends a message down the pipe and
  checks to see that it was received.
*/
func TestAddingPipeListenerToAnInputPipe(t *testing.T) {
	// create pipe
	pipe := &plumbing.Pipe{}

	// create junction
	junction := &plumbing.Junction{PipesMap: map[string]interfaces.IPipeFitting{}, PipeTypesMap: map[string]string{}}

	// create test message
	message := messages.NewMessage(messages.NORMAL, Test{testVal: 1}, nil, messages.PRIORITY_MED)

	// register the pipe with the junction, giving it a name and direction
	registered := junction.RegisterPipe("testInputPipe", plumbing.INPUT, pipe)

	// add the pipelistener using the junction method
	callback := Callback{}
	listenerAdded := junction.AddPipeListener("testInputPipe", callback, callback.CallbackMethod)

	// send the message using our reference to the pipe,
	// it should show up in messageReceived property via the pipeListener
	sent := pipe.Write(message)

	// test assertions
	if pipe == nil {
		t.Error("Expecting pipe not nil")
	}
	if junction == nil {
		t.Error("Expecting junction not nil")
	}
	if registered != true {
		t.Error("Expecting registered pipe")
	}
	if listenerAdded != true {
		t.Error("Expecting added pipeListener")
	}
	if sent != true {
		t.Error("Expecting successful write to pipe")
	}
	if len(callback.messagesReceived) != 1 {
		t.Error("Expecting 1 message received")
	}
	if callback.messagesReceived[0] != message {
		t.Error("Expecting received message was same instance sent")
	}
}

/**
  Test using sendMessage on an OUTPUT pipe.

  Creates a Pipe, Junction and Message.
  Adds the PipeListener to the Pipe.
  Adds the Pipe to the Junction as an OUTPUT pipe.
  uses the Junction's sendMessage method to send
  the Message, then checks that it was received.
*/
func TestSendMessageOnAnOutputPipe(t *testing.T) {
	// create pipe
	pipe := &plumbing.Pipe{}

	// add a PipeListener manually
	callback := Callback{}
	listenerAdded := pipe.Connect(&plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod})

	// create junction
	junction := &plumbing.Junction{PipesMap: map[string]interfaces.IPipeFitting{}, PipeTypesMap: map[string]string{}}

	// create test message
	message := messages.NewMessage(messages.NORMAL, Test{testVal: 1}, nil, messages.PRIORITY_MED)

	// register the pipe with the junction, giving it a name and direction
	registered := junction.RegisterPipe("testOutputPipe", plumbing.OUTPUT, pipe)

	// send the message using the Junction's method
	// it should show up in messageReceived property via the pipeListener
	sent := junction.SendMessage("testOutputPipe", message)

	// test assertions
	//XCTAssertNotNil(pipe as! Pipe, "Expecting pipe is Pipe")
	if pipe == nil {
		t.Error("Expecting pipe not nil")
	}
	if junction == nil {
		t.Error("Expecting junction not nil")
	}
	if registered != true {
		t.Error("Expecting registered pipe")
	}
	if listenerAdded != true {
		t.Error("Expecting added pipeListener")
	}
	if sent != true {
		t.Error("Expecting message sent")
	}
	if len(callback.messagesReceived) != 1 {
		t.Error("Expecting 1 message received")
	}
	if callback.messagesReceived[0] != message {
		t.Error("Expecting received message was same instance sent")
	}

}
