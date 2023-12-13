//
//  QueueControlMessage.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package messages

const (
	FLUSH string = "http://puremvc.org/namespaces/pipes/messages/normal/queue/flush" // Flush the queue.
	SORT  string = "http://puremvc.org/namespaces/pipes/messages/normal/queue/sort"  // Toggle to sort-by-priority operation mode.
	FIFO  string = "http://puremvc.org/namespaces/pipes/messages/normal/queue/fifo"  // Toggle to FIFO operation mode (default behavior)
)

/*
QueueControlMessage Queue Control Message.

A special message for controlling the behavior of a Queue.

When written to a pipeline containing a Queue, the type
of the message is interpreted and acted upon by the Queue.

Unlike filters, multiple serially connected queues aren't
very useful and so they do not require a name. If multiple
queues are connected serially, the message will be acted
upon by the first queue only.
*/
type QueueControlMessage struct {
	Message
}

/*
NewQueueControlMessage Constructor
*/
func NewQueueControlMessage(_type string) *QueueControlMessage {
	return &QueueControlMessage{Message: Message{_type: _type, header: nil, body: nil, priority: PRIORITY_MED}}
}
