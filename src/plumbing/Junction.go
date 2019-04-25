//
//  Junction.go
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

const (
	INPUT  = "input"
	OUTPUT = "output"
)

/**
Pipe Junction.

Manages Pipes for a Module.

When you register a Pipe with a Junction, it is
declared as being an INPUT pipe or an OUTPUT pipe.

You can retrieve or remove a registered Pipe by name,
check to see if a Pipe with a given name exists, or if
it exists AND is an INPUT or an OUTPUT Pipe.

You can send an `IPipeMessage` on a named INPUT Pipe
or add a `PipeListener` to registered INPUT Pipe.
*/
type Junction struct {
	inputPipes        []string
	outputPipes       []string
	PipesMap          map[string]interfaces.IPipeFitting
	PipesMapMutex     sync.RWMutex
	PipeTypesMap      map[string]string
}

/**
  Register a pipe with the junction.

  Pipes are registered by unique name and type,
  which must be either `Junction.INPUT`
  or `Junction.OUTPUT`.

  NOTE: You cannot have an INPUT pipe and an OUTPUT
  pipe registered with the same name. All pipe names
  must be unique regardless of type.

  - parameter name: name of the Pipe Fitting
  - parameter type: input or output
  - parameter pipe: instance of the `IPipeFitting`

  - returns: Bool true if successfully registered. false if another pipe exists by that name.
*/
func (self *Junction) RegisterPipe(name string, type_ string, pipe interfaces.IPipeFitting) bool {
	self.PipesMapMutex.Lock()
	defer self.PipesMapMutex.Unlock()

	success := true
	if self.PipesMap[name] == nil {
		self.PipesMap[name] = pipe
		self.PipeTypesMap[name] = type_
		switch type_ {
		case INPUT:
			self.inputPipes = append(self.inputPipes, name)
		case OUTPUT:
			self.outputPipes = append(self.outputPipes, name)
		default:
			success = false
		}
	} else {
		success = false
	}
	return success
}

/**
  Does this junction have a pipe by this name?

  - parameter name: the pipe to check for
  - returns: Bool whether as pipe is registered with that name.
*/
func (self *Junction) HasPipe(name string) bool {
	self.PipesMapMutex.RLock()
	defer self.PipesMapMutex.RUnlock()

	return self.PipesMap[name] != nil
}

/**
  Does this junction have an INPUT pipe by this name?

  - parameter name: the pipe to check for
  - returns: Bool whether an INPUT pipe is registered with that name.
*/
func (self *Junction) HasInputPipe(name string) bool {
	self.PipesMapMutex.RLock()
	defer self.PipesMapMutex.RUnlock()

	return self.HasPipe(name) && self.PipeTypesMap[name] == INPUT
}

/**
  Does this junction have an OUTPUT pipe by this name?

  - parameter name: the pipe to check for
  - returns: Bool whether an OUTPUT pipe is registered with that name.
*/
func (self *Junction) HasOutputPipe(name string) bool {
	self.PipesMapMutex.RLock()
	defer self.PipesMapMutex.RUnlock()

	return self.HasPipe(name) && self.PipeTypesMap[name] == OUTPUT
}

/**
  Remove the pipe with this name if it is registered.

  NOTE: You cannot have an INPUT pipe and an OUTPUT
  pipe registered with the same name. All pipe names
  must be unique regardless of type.

  - parameter name: the pipe to remove
*/
func (self *Junction) RemovePipe(name string) {
	self.PipesMapMutex.Lock()
	defer self.PipesMapMutex.Unlock()

	if self.PipesMap[name] != nil {
		type_ := self.PipeTypesMap[name]
		var pipesList []string
		switch type_ {
		case INPUT:
			pipesList = self.inputPipes
		case OUTPUT:
			pipesList = self.outputPipes
		}
		for index, pipeName := range pipesList {
			if pipeName == name {
				pipesList = append(pipesList[0:index], pipesList[index+1:]...)
			}
		}
		delete(self.PipesMap, name)
		delete(self.PipeTypesMap, name)
	}
}

/**
  Retrieve the named pipe.

  - parameter name: the pipe to retrieve
  - returns: IPipeFitting the pipe registered by the given name if it exists
*/
func (self *Junction) RetrievePipe(name string) interfaces.IPipeFitting {
	self.PipesMapMutex.RLock()
	defer self.PipesMapMutex.RUnlock()

	return self.PipesMap[name]
}

/**
  Add a PipeListener to an INPUT pipe.

  NOTE: there can only be one PipeListener per pipe, and the listener function must accept an IPipeMessage as its sole argument.

  - parameter inputPipeName: the INPUT pipe to add a PipeListener to
  - parameter context: the calling context or 'this' object
  - parameter listener: the function on the context to call
*/
func (self *Junction) AddPipeListener(inputPipeName string, context interface{}, listener func(message interfaces.IPipeMessage)) bool {
	success := false
	if self.HasInputPipe(inputPipeName) {
		self.PipesMapMutex.Lock()
		pipe := self.PipesMap[inputPipeName]
		success = pipe.Connect(&PipeListener{Context: context, Listener: listener})
		self.PipesMapMutex.Unlock()
	}
	return success
}

/**
  Send a message on an OUTPUT pipe.

  - parameter outputPipeName: the OUTPUT pipe to send the message on
  - parameter message: the IPipeMessage to send
*/
func (self *Junction) SendMessage(outputPipeName string, message interfaces.IPipeMessage) bool {
	self.PipesMapMutex.RLock()
	defer self.PipesMapMutex.RUnlock()

	success := false
	if self.HasOutputPipe(outputPipeName) {
		pipe := self.PipesMap[outputPipeName]
		success = pipe.Write(message)
	}
	return success
}
