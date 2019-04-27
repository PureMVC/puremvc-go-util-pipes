//
//  Queue_test.go
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
Test the Queue class.
*/

/*
  Test connecting input and output pipes to a queue.
*/
func TestConnectingIOPipesQueue(t *testing.T) {
	// create output pipes 1
	pipe1 := &plumbing.Pipe{}
	pipe2 := &plumbing.Pipe{}

	// create queue
	queue := &plumbing.Queue{Pipe: *pipe2, Mode: messages.SORT}

	// connect input fitting
	connectedInput := pipe1.Connect(queue)

	// connect output fitting
	connectedOutput := queue.Connect(pipe2)

	// test assertions
	if pipe1 == nil {
		t.Error("Expecting pipe1 not nil")
	}
	if pipe2 == nil {
		t.Error("Expecting pipe2 not nil")
	}
	if queue == nil {
		t.Error("Expecting queue not nil")
	}
	if connectedInput != true {
		t.Error("Expecting connected input")
	}
	if connectedOutput != true {
		t.Error("Expecting connected output")
	}
}

/*
  Test writing multiple messages to the Queue followed by a Flush message.

  Creates messages to send to the queue.
  Creates queue, attaching an anonymous listener to its output.
  Writes messages to the queue. Tests that no messages have been
  received yet (they've been enqueued). Sends FLUSH message. Tests
  that messages were receieved, and in the order sent (FIFO).
*/
func TestWritingMultipleMessagesAndFlush(t *testing.T) {
	// create messages to send to the queue
	message1 := messages.NewMessage(messages.NORMAL, Test{testVal: 1}, nil, messages.PRIORITY_MED)
	message2 := messages.NewMessage(messages.NORMAL, Test{testVal: 2}, nil, messages.PRIORITY_MED)
	message3 := messages.NewMessage(messages.NORMAL, Test{testVal: 3}, nil, messages.PRIORITY_MED)

	// create queue control flush message
	flushMessage := messages.NewQueueControlMessage(messages.FLUSH)

	// create queue, attaching an anonymous listener to its output
	callback := Callback{}
	queue := &plumbing.Queue{Pipe: plumbing.Pipe{Output: &plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod}}, Mode: messages.SORT}

	// write messages to the queue
	message1written := queue.Write(message1)
	message2written := queue.Write(message2)
	message3written := queue.Write(message3)

	// test assertions
	if message1 == nil {
		t.Error("Expecting message1 not nil")
	}
	if message2 == nil {
		t.Error("Expecting message2 not nil")
	}
	if message3 == nil {
		t.Error("Expecting message3 not nil")
	}
	if flushMessage == nil {
		t.Error("Expecting flushMessage not nil")
	}
	if queue == nil {
		t.Error("Expecting queue not nil")
	}
	if message1written != true {
		t.Error("Expecting wrote message1 to queue")
	}
	if message2written != true {
		t.Error("Expecting wrote message2 to queue")
	}
	if message3written != true {
		t.Error("Expecting wrote message3 to queue")
	}

	// write flush control message to the queue
	flushWritten := queue.Write(flushMessage)

	if flushWritten != true {
		t.Error("Expecting wrote flush message to queue")
	}

	// test that all messages were received, then test
	// FIFO order by inspecting the messages themselves
	if len(callback.messagesReceived) != 3 {
		t.Error("Expecting received 3 messages")
	}

	// test message 1 assertions
	var received1 interfaces.IPipeMessage
	received1, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:] // shift
	if received1 == nil {
		t.Error("Expecting received1 not nil")
	}
	if received1 != message1 {
		t.Error("Expecting received1 == message1")
	}

	// test message 2 assertions
	var received2 interfaces.IPipeMessage
	received2, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:] // shift
	if received2 == nil {
		t.Error("Expecting received2 not nil")
	}
	if received2 != message2 {
		t.Error("Expecting received2 == message2")
	}

	// test message 3 assertions
	var received3 interfaces.IPipeMessage
	received3, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:] // shift
	if received3 == nil {
		t.Error("Expecting received3 not nil")
	}
	if received3 != message3 {
		t.Error("Expecting received3 == message3")
	}
}

/*
  Test the Sort-by-Priority and FIFO modes.

  Creates messages to send to the queue, priorities unsorted.
  Creates queue, attaching an anonymous listener to its output.
  Sends SORT message to start sort-by-priority order mode.
  Writes messages to the queue. Sends FLUSH message, tests
  that messages were receieved in order of priority, not how
  they were sent.

  Then sends a FIFO message to switch the queue back to
  default FIFO behavior, sends messages again, flushes again,
  tests that the messages were recieved and in the order they
  were originally sent.
*/
func TestSortByPriorityAndFIFO(t *testing.T) {
	// create messages to send to the queue
	message1 := messages.NewMessage(messages.NORMAL, nil, nil, messages.PRIORITY_MED)
	message2 := messages.NewMessage(messages.NORMAL, nil, nil, messages.PRIORITY_LOW)
	message3 := messages.NewMessage(messages.NORMAL, nil, nil, messages.PRIORITY_HIGH)

	// create queue, attaching an anonymous listener to its output
	callback := Callback{}
	queue := plumbing.Queue{Pipe: plumbing.Pipe{Output: &plumbing.PipeListener{Context: callback, Listener: callback.CallbackMethod}}}

	// begin sort-by-priority order mode
	sortWritten := queue.Write(messages.NewQueueControlMessage(messages.SORT))

	// write messages to the queue
	message1Written := queue.Write(message1)
	message2Written := queue.Write(message2)
	message3Written := queue.Write(message3)

	// test assertions
	if sortWritten != true {
		t.Error("Expecting wrote sort message to queue")
	}
	if message1Written != true {
		t.Error("Expecting wrote message1 to queue")
	}
	if message2Written != true {
		t.Error("Expecting wrote message2 to queue")
	}
	if message3Written != true {
		t.Error("Expecting wrote message3 to queue")
	}

	// flush the queue
	flushWritten := queue.Write(messages.NewQueueControlMessage(messages.FLUSH))

	if flushWritten != true {
		t.Error("Expecting worte flush message to queue")
	}

	// test that 3 messages were received
	if len(callback.messagesReceived) != 3 {
		t.Error("Expecting received 3 messages")
	}

	// get the messages
	var received1 interfaces.IPipeMessage
	received1, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:] // shift
	var received2 interfaces.IPipeMessage
	received2, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:] // shift
	var received3 interfaces.IPipeMessage
	received3, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:] // shift

	// test that the message order is sorted
	if received1.Priority() > received2.Priority() {
		t.Error("Expecting received1 is higher priority than received2")
	}
	if received2.Priority() > received3.Priority() {
		t.Error("Expecting received2 is higher priority than received3")
	}
	if received1 != message3 {
		t.Error("Expecting received1 == message3")
	}
	if received2 != message1 {
		t.Error("Expecting received2 == message1")
	}
	if received3 != message2 {
		t.Error("Expecting received3 == message2")
	}

	// begin FIFO order mode
	fifoWritten := queue.Write(messages.NewQueueControlMessage(messages.FIFO))

	// write messages to the queue
	message1WrittenAgain := queue.Write(message1)
	message2WrittenAgain := queue.Write(message2)
	message3WrittenAgain := queue.Write(message3)

	// flush the queue
	flushWrittenAgain := queue.Write(messages.NewQueueControlMessage(messages.FLUSH))

	// test assertions
	if fifoWritten != true {
		t.Error("Expecting wrote fifo message to queue")
	}
	if message1WrittenAgain != true {
		t.Error("Expecting wrote message1 to queue again")
	}
	if message2WrittenAgain != true {
		t.Error("Expecting wrote message2 to queue again")
	}
	if message3WrittenAgain != true {
		t.Error("Expecting wrote message3 to queue again")
	}
	if flushWrittenAgain != true {
		t.Error("Expecting wrote flush message to queue again")
	}

	// test that 3 messages were received
	if len(callback.messagesReceived) != 3 {
		t.Error("Expecting received 3 messages")
	}

	// get the messages
	var received1Again interfaces.IPipeMessage
	received1Again, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:] // shift
	var received2Again interfaces.IPipeMessage
	received2Again, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:] // shift
	var received3Again interfaces.IPipeMessage
	received3Again, callback.messagesReceived = callback.messagesReceived[0], callback.messagesReceived[1:] // shift

	// test message order is FIFO
	if received1Again != message1 {
		t.Error("Expecting received1Again == message1")
	}
	if received2Again != message2 {
		t.Error("Expecting received2Again == message2")
	}
	if received3Again != message3 {
		t.Error("Expecting received3Again == message3")
	}
	if received1Again.Priority() != messages.PRIORITY_MED {
		t.Error("Expecting received1Again is priority med")
	}
	if received2Again.Priority() != messages.PRIORITY_LOW {
		t.Error("Expecting received2Again is priority low")
	}
	if received3Again.Priority() != messages.PRIORITY_HIGH {
		t.Error("Expecting received3Again is priority high")
	}
}