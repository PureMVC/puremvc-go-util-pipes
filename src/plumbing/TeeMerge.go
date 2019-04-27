//
//  TeeMerge.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import "github.com/puremvc/puremvc-go-util-pipes/src/interfaces"

/*
Merging Pipe Tee.

Writes the messages from multiple input pipelines into
a single output pipe fitting.
*/
type TeeMerge struct {
	Pipe
}

/*
  Connect an input IPipeFitting.

  NOTE: You can connect as many inputs as you want
  by calling this method repeatedly.

  - parameter input: the IPipeFitting to connect for input.
*/
func (self *TeeMerge) ConnectInput(input interfaces.IPipeFitting) bool {
	return input.Connect(self)
}
