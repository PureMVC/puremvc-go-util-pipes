//
//  FilterControlMessage.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package messages

import "github.com/puremvc/puremvc-go-util-pipes/src/interfaces"

const (
	SET_PARAMS = "http://puremvc.org/namespaces/pipes/messages/normal/filter-control/setParams" // Set filter parameters.
	SET_FILTER = "http://puremvc.org/namespaces/pipes/messages/normal/filter-control/setFilter" // Set filter function.
	BYPASS     = "http://puremvc.org/namespaces/pipes/messages/normal/filter-control/bypass"    // Toggle to filter bypass mode.
	FILTER     = "http://puremvc.org/namespaces/pipes/messages/normal/filter-control/filter"    // Toggle to filtering mode. (default behavior).
)

/**
Filter Control Message.

A special message type for controlling the behavior of a Filter.

The `messages.SET_PARAMS` message type tells the Filter
to retrieve the filter parameters object.

The `messages.SET_FILTER` message type tells the Filter
to retrieve the filter function.

The `messages.BYPASS` message type tells the Filter
that it should go into Bypass mode operation, passing all normal
messages through unfiltered.

The `messages.FILTER` message type tells the Filter
that it should go into Filtering mode operation, filtering all
normal normal messages before writing out. This is the default
mode of operation and so this message type need only be sent to
cancel a previous `FilterControlMessage.BYPASS` message.

The Filter only acts on a control message if it is targeted
to this named filter instance. Otherwise it writes the
message through to its output unchanged.
*/
type FilterControlMessage struct {
	Message
	name   string
	filter func(interfaces.IPipeMessage, interface{}) bool
	params interface{}
}

// Constructor
func NewFilterControlMessage(type_ string, name string, filter func(interfaces.IPipeMessage, interface{}) bool, params interface{}) *FilterControlMessage {
	return &FilterControlMessage{Message: Message{_type: type_}, name: name, filter: filter, params: params}
}

/**
  Set the target filter name.
*/
func (self *FilterControlMessage) SetName(name string) {
	self.name = name
}

/**
  Get the target filter name.
*/
func (self *FilterControlMessage) Name() string {
	return self.name
}

/**
  Set the filter function.
*/
func (self *FilterControlMessage) SetFilter(filter func(interfaces.IPipeMessage, interface{}) bool) {
	self.filter = filter
}

/**
  Get the filter function.
*/
func (self *FilterControlMessage) Filter() func(interfaces.IPipeMessage, interface{}) bool {
	return self.filter
}

/**
  Set the parameters object.
*/
func (self *FilterControlMessage) SetParams(params interface{}) {
	self.params = params
}

/**
  Get the parameters object.
*/
func (self *FilterControlMessage) Params() interface{} {
	return self.params
}
