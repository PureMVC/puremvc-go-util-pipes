//
//  IPipeMessage.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package interfaces

/*
Pipe Message Interface.

IPipeMessages are objects written intoto a Pipeline,
composed of IPipeFittings. The message is passed from
one fitting to the next in syncrhonous fashion.

Depending on type, messages may be handled differently by the
fittings.
*/
type IPipeMessage interface {
	Type() string                 // Get type of this message
	SetType(_type string)         // Set type of this message
	Priority() int                // Get priority of this message
	SetPriority(priority int)     // Set priority of this message
	Header() interface{}          // Get header of this message
	SetHeader(header interface{}) // Set header of this message
	Body() interface{}            // Get body of this message
	SetBody(body interface{})     // Set body of this message
}
