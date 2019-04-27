//
//  IPipeAware.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package interfaces

/*
Pipe Aware interface.

Can be implemented by any PureMVC Core that wishes
to communicate with other Cores using the Pipes
utility.
*/
type IPipeAware interface {
	/*
	  Connect input Pipe Fitting.

	  - parameter name: name of the input pipe
	  - parameter pipe: input Pipe Fitting
	*/
	AcceptInputPipe(name string, pipe IPipeFitting)

	/*
	  Connect output Pipe Fitting.

	  - parameter name: name of the input pipe
	  - parameter pipe: output Pipe Fitting
	*/
	AcceptOutputPipe(name string, pipe IPipeFitting)
}
