//
//  JunctionMediator.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import (
	puremvc "github.com/puremvc/puremvc-go-multicore-framework/src/interfaces"
	"github.com/puremvc/puremvc-go-multicore-framework/src/patterns/mediator"
	"github.com/puremvc/puremvc-go-util-pipes/src/interfaces"
)

const (
	ACCEPT_INPUT_PIPE  = "acceptInputPipe"
	ACCEPT_OUTPUT_PIPE = "acceptOutputPipe"
)

/**
Junction Mediator.

A base class for handling the Pipe Junction in an IPipeAware
Core.
*/
type JunctionMediator struct {
	mediator.Mediator
}

/**
  List Notification Interests.

  Returns the notification interests for this base class.
  Override in subclass and call `super.listNotificationInterests`
  to get this list, then add any sublcass interests to
  the array before returning.
*/
func (self *JunctionMediator) ListNotificationInterests() []string {
	return []string{
		ACCEPT_INPUT_PIPE,
		ACCEPT_OUTPUT_PIPE}
}

/**
  Handle Notification.

  This provides the handling for common junction activities. It
  accepts input and output pipes in response to `IPipeAware`
  interface calls.

  Override in subclass, and call `super.handleNotification`
  if none of the subclass-specific notification names are matched.
*/
func (self *JunctionMediator) HandleNotification(notification puremvc.INotification) {
	switch notification.Name() {
	// accept an input pipe
	// register the pipe and if successful
	// set this mediator as its listener
	case ACCEPT_INPUT_PIPE:
		inputPipeName := notification.Type()
		inputPipe := notification.Body().(interfaces.IPipeFitting)
		if self.Junction().RegisterPipe(inputPipeName, INPUT, inputPipe) {
			self.Junction().AddPipeListener(inputPipeName, self, self.HandlePipeMessage)
		}
	case ACCEPT_OUTPUT_PIPE: // accept an output pipe
		outputPipeName := notification.Type()
		outputPipe := notification.Body().(interfaces.IPipeFitting)
		self.Junction().RegisterPipe(outputPipeName, OUTPUT, outputPipe)
	}
}

/**
  Handle incoming pipe messages.

  Override in subclass and handle messages appropriately for the module.
*/
func (self *JunctionMediator) HandlePipeMessage(message interfaces.IPipeMessage) {

}

/**
  The Junction for this Module.
*/
func (self *JunctionMediator) Junction() *Junction {
	return self.ViewComponent.(*Junction)
}
