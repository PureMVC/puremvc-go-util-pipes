//
//  Message.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package messages

import "github.com/puremvc/puremvc-go-util-pipes/src/interfaces"

const (
	PRIORITY_HIGH = 1                                                      // High priority Messages can be sorted to the front of the queue
	PRIORITY_MED  = 5                                                      // Medium priority Messages are the default
	PRIORITY_LOW  = 10                                                     // Low priority Messages can be sorted to the back of the queue
	NORMAL        = "http://puremvc.org/namespaces/pipes/messages/normal/" // Normal Message type
)

/*
Message Pipe Message.

Messages travelling through a Pipeline can
be filtered, and queued. In a queue, they may
be sorted by priority. Based on type,
they may used as control messages to modify the
behavior of filter or queue fittings connected
to the pipeline into which they are written.
*/
type Message struct {
	_type    string
	header   interface{}
	body     interface{}
	priority int
}

/*
NewMessage Constructor
*/
func NewMessage(_type string, header interface{}, body interface{}, priority int) interfaces.IPipeMessage {
	return &Message{_type: _type, header: header, body: body, priority: priority}
}

/*
Type Get the type of this message
*/
func (self *Message) Type() string {
	return self._type
}

/*
SetType Set the type of this message
*/
func (self *Message) SetType(_type string) {
	self._type = _type
}

/*
Priority Get the priority of this message
*/
func (self *Message) Priority() int {
	return self.priority
}

/*
SetPriority Set the priority of this message
*/
func (self *Message) SetPriority(priority int) {
	self.priority = priority
}

/*
Header Get the header of this message
*/
func (self *Message) Header() interface{} {
	return self.header
}

/*
SetHeader Set the header of this message
*/
func (self *Message) SetHeader(header interface{}) {
	self.header = header
}

/*
Body Get the body of this message
*/
func (self *Message) Body() interface{} {
	return self.body
}

/*
SetBody Set the body of this message
*/
func (self *Message) SetBody(body interface{}) {
	self.body = body
}
