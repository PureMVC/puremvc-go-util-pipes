//
//  Filter.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import (
	"github.com/puremvc/puremvc-go-util-pipes/src/interfaces"
	"github.com/puremvc/puremvc-go-util-pipes/src/messages"
)

/*
Filter Pipe Filter.

Filters may modify the contents of messages before writing them to
their output pipe fitting. They may also have their parameters and
filter function passed to them by control message, as well as having
their Bypass/Filter operation mode toggled via control message.
*/
type Filter struct {
	Pipe
	Name   string
	Filter func(message interfaces.IPipeMessage, params interface{}) bool
	Params interface{}
	Mode   string
}

/*
Write Handle the incoming message.

If message type is normal, filter the message (unless in BYPASS mode)
and write the result to the output pipe fitting if the filter
operation is successful.

The messages.SET_PARAMS message type tells the Filter
that the message class is FilterControlMessage, which it
casts the message to in order to retrieve the filter parameters
object if the message is addressed to this filter.

The messages.SET_FILTER message type tells the Filter
that the message class is FilterControlMessage,
which it casts the message to in order to retrieve the filter function.

The messages.BYPASS message type tells the Filter
that it should go into Bypass mode operation, passing all normal
messages through unfiltered.

The messages.FILTER message type tells the Filter
that it should go into Filtering mode operation, filtering all
normal messages before writing out. This is the default
mode of operation and so this message type need only be sent to
cancel a previous BYPASS message.

The Filter only acts on the control message if it is targeted
to this named filter instance. Otherwise, it writes through to the
output.

- parameter message: IPipeMessage to write on the output

- returns: Boolean True if the filter process does not throw an error and subsequent operations
in the pipeline succeeds.
*/
func (self *Filter) Write(message interfaces.IPipeMessage) bool {
	success := true

	switch message.Type() {
	case messages.NORMAL: // Filter normal messages
		if self.Mode == messages.FILTER {
			if self.ApplyFilter(message) {
				success = self.Output.Write(message)
			} else {
				success = false
			}
		} else {
			success = self.Output.Write(message)
		}
	case messages.SET_PARAMS: // Accept parameters from control message
		if self.IsTarget(message) {
			self.Params = message.(*messages.FilterControlMessage).Params()
		} else {
			success = self.Output.Write(message)
		}

	case messages.SET_FILTER: // Accept filter function from control message
		if self.IsTarget(message) {
			self.Filter = message.(*messages.FilterControlMessage).Filter()
		} else {
			success = self.Output.Write(message)
		}
		// Toggle between Filter or Bypass operational modes
	case messages.BYPASS:
		fallthrough
	case messages.FILTER:
		if self.IsTarget(message) {
			self.Mode = message.(*messages.FilterControlMessage).Type()
		} else {
			success = self.Output.Write(message)
		}
	default: // Write control messages for other fittings through
		success = self.Output.Write(message)
	}

	return success
}

// IsTarget Is the message directed at this filter instance?
func (self *Filter) IsTarget(message interfaces.IPipeMessage) bool {
	return message.(*messages.FilterControlMessage).Name() == self.Name
}

// ApplyFilter Filter the message.
func (self *Filter) ApplyFilter(message interfaces.IPipeMessage) bool {
	return self.Filter(message, self.Params)
}
