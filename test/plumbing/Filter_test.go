//
//  Filter_test.go
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

/*
Test the Filter class.
*/

/*
Test connecting input and output pipes to a filter as well as disconnecting the output.
*/
func TestConnectingAndDisconnectingIOPipesFilter(t *testing.T) {
	// create output pipes 1
	pipe1 := &plumbing.Pipe{}
	pipe2 := &plumbing.Pipe{}

	//create filter
	filter := &plumbing.Filter{Name: "TestFilter"}

	// connect input fitting
	connectedInput := pipe1.Connect(filter)

	// connect output fitting
	connectedOutput := filter.Connect(pipe2)

	// test assertions
	if pipe1 == nil {
		t.Error("Expecting pipe1 is not nil")
	}
	if pipe2 == nil {
		t.Error("Expecting pipe2 is not nil")
	}
	if filter == nil {
		t.Error("Expecting filter not nil")
	}
	if connectedInput != true {
		t.Error("Expecting connected input")
	}
	if connectedOutput != true {
		t.Error("Expecting connected output")
	}

	// disconnect pipe 2 from filter
	disconnectedPipe := filter.Disconnect()
	if disconnectedPipe != pipe2 {
		t.Error("Expecting disconnected pipe2 from filter")
	}
}

/*
Test applying filter to a normal message.
*/
func TestFilteringNormalMessage(t *testing.T) {
	// create messages to send to the queue
	message := messages.NewMessage(messages.NORMAL, &Rect{Width: 10, Height: 2}, nil, messages.PRIORITY_MED)

	// create filter, attach an anonymous listener to the filter output to receive the message,
	// pass in an anonymous function an parameter object
	callback := Callback{}
	filter := &plumbing.Filter{
		Name: "scale",
		Pipe: plumbing.Pipe{Output: &plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod}},
		Filter: func(message interfaces.IPipeMessage, params interface{}) bool {
			message.Header().(*Rect).Width *= float32(params.(Factor).factor)
			message.Header().(*Rect).Height *= float32(params.(Factor).factor)
			return true
		},
		Params: Factor{factor: 10},
		Mode:   messages.FILTER}

	// write messages to the filter
	written := filter.Write(message)

	// test assertions
	if message == nil {
		t.Error("Expecting message not nil")
	}
	if filter == nil {
		t.Error("Expecting filter not nil")
	}
	if written != true {
		t.Error("Expecting wrote message to filter")
	}
	if len(callback.messagesReceived) != 1 {
		t.Error("Expecting received 1 messages")
	}

	// test filtered message assertions
	received := callback.messagesReceived[0]
	if received == nil {
		t.Error("Expecting received not nil")
	}
	if received != message {
		t.Error("Expecting received == message")
	}
	if received.Header().(*Rect).Width != 100 {
		t.Error("Expecting received.Header().(*Rect).Width == 100")
	}
	if received.Header().(*Rect).Height != 20 {
		t.Error("Expecting received.Header().(*Rect).Height == 20")
	}
}

/*
Test setting filter to bypass mode, writing, then setting back to filter mode and writing.
*/
func TestBypassAndFilterModeToggle(t *testing.T) {
	// create messages to send to the queue
	message := messages.NewMessage(messages.NORMAL, &Rect{Width: 10, Height: 2}, nil, messages.PRIORITY_MED)

	// create filter, attach an anonymous listener to the filter output to receive the message,
	// pass in an anonymous function an parameter object
	callback := Callback{}
	filter := &plumbing.Filter{
		Name: "scale",
		Pipe: plumbing.Pipe{Output: &plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod}},
		Filter: func(message interfaces.IPipeMessage, params interface{}) bool {
			message.Header().(*Rect).Width *= float32(params.(Factor).factor)
			message.Header().(*Rect).Height *= float32(params.(Factor).factor)
			return true
		},
		Params: Factor{factor: 10},
		Mode:   messages.FILTER}

	// create bypass control message
	bypassMessage := messages.NewFilterControlMessage(messages.BYPASS, "scale", nil, nil)

	// write bypass control message to the filter
	bypassWritten := filter.Write(bypassMessage)

	// write normal message to the filter
	written1 := filter.Write(message)

	// test assertions
	if message == nil {
		t.Error("Expecting message not nil")
	}
	if filter == nil {
		t.Error("Expecting filter not nil")
	}
	if bypassWritten != true {
		t.Error("Expecting wrote bypass message to filter")
	}
	if written1 != true {
		t.Error("Expecting wrote normal message to filter")
	}
	if len(callback.messagesReceived) != 1 {
		t.Error("Expecting received 1 messages")
	}

	// test filtered message assertions (no change to message)
	var received1 interfaces.IPipeMessage
	received1, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:]
	if received1 == nil {
		t.Error("Expecting received1 not nil")
	}
	if received1 != message {
		t.Error("Expecting received1 == message")
	}
	if received1.Header().(*Rect).Width != 10 {
		t.Error("Expecting received1.Header().(*Rect).Width == 10")
	}
	if received1.Header().(*Rect).Height != 2 {
		t.Error("Expecting received1.Header().(*Rect).Height == 2")
	}

	// create filter control message
	filterMessage := messages.NewFilterControlMessage(messages.FILTER, "scale", nil, nil)

	// write bypass control message to the filter
	filterWritten := filter.Write(filterMessage)

	//let write normal message to the filter again
	written2 := filter.Write(message)

	// test assertions
	if filterWritten != true {
		t.Error("Expecting wrote filter message to filter")
	}
	if written2 != true {
		t.Error("Expecting wrote normal message to filter")
	}
	if len(callback.messagesReceived) != 1 {
		t.Error("Expecting wrote 1 messages")
	}

	// test filtered message assertions (message filtered)
	var received2 interfaces.IPipeMessage
	received2, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:]
	if received2 == nil {
		t.Error("Expecting received2 not nil")
	}
	if received2 != message {
		t.Error("Expecting received2 == message")
	}
	if received2.Header().(*Rect).Width != 100 {
		t.Error("Expecting received2.Header().(*Rect).Width != 100")
	}
	if received2.Header().(*Rect).Height != 20 {
		t.Error("Expecting received2.Header().(*Rect).Height != 20")
	}
}

/*
Test setting filter parameters by sending control message.
*/
func TestSetParamsByControlMessage(t *testing.T) {
	// create messages to send to the queue
	message := messages.NewMessage(messages.NORMAL, &Rect{Width: 10, Height: 2}, nil, messages.PRIORITY_MED)

	// create filter, attach an anonymous listener to the filter output to receive the message,
	// pass in an anonymous function an parameter object
	callback := Callback{}
	filter := &plumbing.Filter{
		Name: "scale",
		Pipe: plumbing.Pipe{Output: &plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod}},
		Filter: func(message interfaces.IPipeMessage, params interface{}) bool {
			message.Header().(*Rect).Width *= float32(params.(Factor).factor)
			message.Header().(*Rect).Height *= float32(params.(Factor).factor)
			return true
		},
		Params: Factor{factor: 10},
		Mode:   messages.FILTER}

	// create setParams control message
	setParamsMessage := messages.NewFilterControlMessage(messages.SET_PARAMS, "scale", nil, Factor{5})

	// write filter control message to the filter
	setParamsWritten := filter.Write(setParamsMessage)

	// write normal message to the filter
	written := filter.Write(message)

	// test assertions
	if message == nil {
		t.Error("Expecting message is nil")
	}
	if filter == nil {
		t.Error("Expecting filter is nil")
	}
	if setParamsWritten != true {
		t.Error("Expecting wrote set_params message")
	}
	if written != true {
		t.Error("Expecting wrote normal message to filter")
	}
	if len(callback.messagesReceived) != 1 {
		t.Error("Expecting received 1 messages")
	}

	// test filtered message assertions (message filtered with overridden parameters)
	received := callback.messagesReceived[0]
	if received == nil {
		t.Error("Expecting received not nil")
	}
	if received != message {
		t.Error("Expecting received == message")
	}
	if received.Header().(*Rect).Width != 50 {
		t.Error("Expecting received.Header().(*Rect).Width == 50")
	}
	if received.Header().(*Rect).Height != 10 {
		t.Error("Expecting received.Header().(*Rect).Height == 10")
	}
}

/*
Test setting filter function by sending control message.
*/
func TestSetFilterByControlMessage(t *testing.T) {
	// create messages to send to the queue
	message := messages.NewMessage(messages.NORMAL, &Rect{Width: 10, Height: 2}, nil, messages.PRIORITY_MED)

	// create filter, attach an anonymous listener to the filter output to receive the message,
	// pass in an anonymous function and an anonymous parameter object
	callback := Callback{}
	filter := &plumbing.Filter{
		Name: "scale",
		Pipe: plumbing.Pipe{Output: &plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod}},
		Filter: func(message interfaces.IPipeMessage, params interface{}) bool {
			message.Header().(*Rect).Width *= float32(params.(Factor).factor)
			message.Header().(*Rect).Height *= float32(params.(Factor).factor)
			return true
		},
		Params: Factor{factor: 10},
		Mode:   messages.FILTER}

	// create setFilter control message
	setFilterMessage := messages.NewFilterControlMessage(messages.SET_FILTER, "scale",
		func(message interfaces.IPipeMessage, params interface{}) bool {
			message.Header().(*Rect).Width /= float32(params.(Factor).factor)
			message.Header().(*Rect).Height /= float32(params.(Factor).factor)
			return true
		}, nil)

	// write filter control message to the filter
	setFilterWritten := filter.Write(setFilterMessage)

	// write normal message to the filter
	written := filter.Write(message)

	// test assertions
	if message == nil {
		t.Error("Expecting message not nil")
	}
	if filter == nil {
		t.Error("Expecting filter not nil")
	}
	if setFilterWritten != true {
		t.Error("Expecting wrote message to filter")
	}
	if written != true {
		t.Error("Expecting wrote normal message to filter")
	}
	if len(callback.messagesReceived) != 1 {
		t.Error("Expecting received 1 messages")
	}

	// test filtered message assertions (message filtered with overridden filter function)
	received := callback.messagesReceived[0]
	if received == nil {
		t.Error("Expecting received not nil")
	}
	if received != message {
		t.Error("Expecting received == message")
	}
	if received.Header().(*Rect).Width != 1 {
		t.Error("Expecting received.Header().(*Rect).Width == 1")
	}
	if received.Header().(*Rect).Height != .2 {
		t.Error("Expecting received.Header().(*Rect).Height == .2")
	}
}

/*
Test using a filter function to stop propagation of a message.

The way to stop propagation of a message from within a filter
is to throw an error from the filter function. This test creates
two NORMAL messages, each with Rectangle objects that contain
a bozoLevel property. One has this property set to
10, the other to 3.

Creates a Filter, named 'bozoFilter' with an anonymous pipe listener
feeding the output back into this test. The filter funciton is an
anonymous function that throws an error if the message's bozoLevel
property is greater than the filter parameter bozoThreshold.
the anonymous filter parameters object has a bozoThreshold
value of 5.

The messages are written to the filter and it is shown that the
message with the bozoLevel of 10 is not written, while
the message with the bozoLevel of 3 is.
*/
func TestUseFilterToStopAMessage(t *testing.T) {
	// create messages to send to the queue
	message1 := messages.NewMessage(messages.NORMAL, &Bozo{level: 10, name: "Dastardly Dan"}, nil, messages.PRIORITY_MED)
	message2 := messages.NewMessage(messages.NORMAL, &Bozo{level: 3, name: "Dudely Doright"}, nil, messages.PRIORITY_MED)

	// create filter, attach an anonymous listener to the filter output to receive the message,
	// pass in an anonymous function and an anonymous parameter object
	callback := Callback{}
	filter := &plumbing.Filter{Name: "bozoFilter",
		Pipe: plumbing.Pipe{Output: &plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod}},
		Filter: func(message interfaces.IPipeMessage, params interface{}) bool {
			if message.Header().(*Bozo).level > params.(BozoThreshold).level {
				return false
			} else {
				return true
			}
		},
		Params: BozoThreshold{level: 5},
		Mode:   messages.FILTER}

	// write normal message to the filter
	written1 := filter.Write(message1)
	written2 := filter.Write(message2)

	// test assertions
	if message1 == nil {
		t.Error("Expecting message1 not nil")
	}
	if message2 == nil {
		t.Error("Expecting message2 not nil")
	}
	if filter == nil {
		t.Error("Expecting filter not nil")
	}
	if written1 != false {
		t.Error("Expecting failed to write bad message")
	}
	if written2 != true {
		t.Error("Expecting wrote good message")
	}
	if len(callback.messagesReceived) != 1 {
		t.Error("Expecting received 1 messages")
	}

	// test filtered message assertions (message with good auth token passed)
	received := callback.messagesReceived[0]
	if received == nil {
		t.Error("Expecting received is Message")
	}
	if received != message2 {
		t.Error("Expecting received == message2")
	}
}
