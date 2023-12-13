//
//  Queue.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import (
	"github.com/puremvc/puremvc-go-util-pipes/src/interfaces"
	"github.com/puremvc/puremvc-go-util-pipes/src/messages"
	"sort"
	"sync"
)

/*
Queue Pipe Queue.

The Queue always stores inbound messages until you send it
a FLUSH control message, at which point it writes it's buffer
to the output pipe fitting. The Queue can be sent a SORT
control message to go into sort-by-priority mode or a FIFO
control message to cancel sort mode and return the
default mode of operation, FIFO.

NOTE: There can effectively be only one Queue on a given
pipeline, since the first Queue acts on any queue control
message. Multiple queues in one pipeline are of dubious
use, and so having to name them would make their operation
more complex than need be.
*/
type Queue struct {
	Pipe
	Mode          string
	Messages      []interfaces.IPipeMessage
	MessagesMutex sync.Mutex
}

/**
 * Write Handle the incoming message.
 *
 * Normal messages are enqueued.
 *
 * The FLUSH message type tells the Queue to write all
 * stored messages to the output PipeFitting, then
 * return to normal enqueuing operation.
 *
 * The SORT message type tells the Queue to sort all
 * <I>subsequent</I> incoming messages by priority. If there
 * are un-flushed messages in the queue, they will not be
 * sorted unless a new message is sent before the next FLUSH.
 * Sorting-by-priority behavior continues even after a FLUSH,
 * and can be turned off by sending a FIFO message, which is
 * the default behavior for enqueue/dequeue.
 */
func (self *Queue) Write(message interfaces.IPipeMessage) bool {
	success := true

	switch message.Type() {
	case messages.NORMAL: // Store normal messages
		self.Store(message)

	case messages.FLUSH: // Flush the queue
		success = self.Flush()
		// Put Queue into Priority Sort or FIFO mode
		// Subsequent messages written to the queue
		// will be affected. Sorted messages cannot
		// be put back into FIFO order!
	case messages.SORT:
		fallthrough
	case messages.FIFO:
		self.Mode = message.Type()
	}
	return success
}

/*
Store a message.

- parameter message: the IPipeMessage to enqueue.
*/
func (self *Queue) Store(message interfaces.IPipeMessage) {
	self.MessagesMutex.Lock()
	defer self.MessagesMutex.Unlock()

	self.Messages = append(self.Messages, message)
	if self.Mode == messages.SORT {
		sort.Sort(SortByPriority(self.Messages))
	}
}

/*
SortByPriority Sort the Messages by priority.
*/
type SortByPriority []interfaces.IPipeMessage

func (s SortByPriority) Len() int {
	return len(s)
}
func (s SortByPriority) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortByPriority) Less(i, j int) bool {
	return s[i].Priority() < s[j].Priority()
}

/*
Flush the queue.

NOTE: This empties the queue.

- returns: Bool true if all messages written successfully.
*/
func (self *Queue) Flush() bool {
	self.MessagesMutex.Lock()
	defer self.MessagesMutex.Unlock()

	success := true
	if len(self.Messages) > 0 {
		for message := self.Messages[0]; message != nil; {
			ok := self.Pipe.Write(message)
			if ok == false {
				success = false
			}

			self.Messages = self.Messages[1:]

			if len(self.Messages) >= 1 {
				message = self.Messages[0]
			} else {
				message = nil
			}
		}
	}
	return success
}
