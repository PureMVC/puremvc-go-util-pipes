//
//  Pipe.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import "github.com/puremvc/puremvc-go-util-pipes/src/interfaces"

/*
Pipe Pipe.

This is the most basic IPipeFitting,
simply allowing the connection of an output
fitting and writing of a message to that output.
*/
type Pipe struct {
	Output interfaces.IPipeFitting
}

/*
Connect another PipeFitting to the output.

PipeFittings connect to and write to other
PipeFittings in a one-way, synchronous chain.

- parameter output: IPipeFitting the output fitting to connect.

- returns: Bool true if no other fitting was already connected.
*/
func (self *Pipe) Connect(output interfaces.IPipeFitting) bool {
	success := false
	if self.Output == nil {
		self.Output = output
		success = true
	}
	return success
}

/*
Disconnect the Pipe Fitting connected to the output.

This disconnects the output fitting, returning a
reference to it. If you were splicing another fitting
into a pipeline, you need to keep (at least briefly)
a reference to both sides of the pipeline in order to
connect them to the input and output of whatever
fitting that you're splicing in.

- returns: IPipeFitting the now disconnected output fitting
*/
func (self *Pipe) Disconnect() interfaces.IPipeFitting {
	disconnectedFitting := self.Output
	self.Output = nil
	return disconnectedFitting
}

/*
Write the message to the connected output.

- parameter message: the message to write

- returns: Bool whether any connected down-pipe outputs failed
*/
func (self *Pipe) Write(message interfaces.IPipeMessage) bool {
	return self.Output.Write(message)
}
