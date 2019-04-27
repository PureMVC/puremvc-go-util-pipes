//
//  PipeListener.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import "github.com/puremvc/puremvc-go-util-pipes/src/interfaces"

/*
Pipe Listener.

Allows a class that does not implement IPipeFitting to
be the final recipient of the messages in a pipeline.
*/
type PipeListener struct {
	Context  interface{}
	Listener func(message interfaces.IPipeMessage)
}

/*
  Can't connect anything beyond this.
 */
func (self *PipeListener) Connect(output interfaces.IPipeFitting) bool {
	return false
}

/*
  Can't disconnect since you can't connect, either.
 */
func (self *PipeListener) Disconnect() interfaces.IPipeFitting {
	return nil
}

/*
  Write the message to the listener
 */
func (self *PipeListener) Write(message interfaces.IPipeMessage) bool {
	self.Listener(message)
	return true
}
