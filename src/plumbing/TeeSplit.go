//
//  TeeSplit.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import (
	"github.com/puremvc/puremvc-go-util-pipes/src/interfaces"
	"sync"
)

/*
TeeSplit Splitting Pipe Tee.

Writes input messages to multiple output pipe fittings.
*/
type TeeSplit struct {
	outputs      []interfaces.IPipeFitting
	outputsMutex sync.RWMutex // Mutex for messagesQueue
}

/*
Connect the output IPipeFitting.

NOTE: You can connect as many outputs as you want
by calling this method repeatedly.

- parameter output: the IPipeFitting to connect for output.
*/
func (self *TeeSplit) Connect(output interfaces.IPipeFitting) bool {
	self.outputsMutex.Lock()
	defer self.outputsMutex.Unlock()

	self.outputs = append(self.outputs, output)
	return true
}

/*
Disconnect the most recently connected output fitting. (LIFO)

To disconnect all outputs, you must call this
method repeatedly untill it returns nil.

- parameter output: the IPipeFitting to connect for output.
*/
func (self *TeeSplit) Disconnect() interfaces.IPipeFitting {
	self.outputsMutex.Lock()
	defer self.outputsMutex.Unlock()

	disconnectedFitting := self.outputs[len(self.outputs)-1]
	self.outputs = self.outputs[:len(self.outputs)-1]
	return disconnectedFitting
}

/*
DisconnectFitting Disconnect a given output fitting.

If the fitting passed in is connected
as an output of this TeeSplit, then
it is disconnected and the reference returned.

If the fitting passed in is not connected as an
output of this TeeSplit, then nil
is returned.

- parameter output: the IPipeFitting to connect for output.
*/
func (self *TeeSplit) DisconnectFitting(target interfaces.IPipeFitting) interfaces.IPipeFitting {
	self.outputsMutex.Lock()
	defer self.outputsMutex.Unlock()

	var removed interfaces.IPipeFitting

	for index, pipe := range self.outputs {
		if pipe == target {
			removed = self.outputs[index]
			self.outputs = append(self.outputs[:index], self.outputs[index+1:]...)
		}
	}
	return removed
}

/*
Write the message to all connected outputs.

Returns false if any output returns false,
but all outputs are written to regardless.

- parameter message: the message to write

- returns: Boolean whether any connected outputs failed
*/
func (self *TeeSplit) Write(message interfaces.IPipeMessage) bool {
	self.outputsMutex.RLock()
	defer self.outputsMutex.RUnlock()

	success := true
	for _, pipe := range self.outputs {
		if pipe.Write(message) == false {
			success = false
		}
	}
	return success
}
