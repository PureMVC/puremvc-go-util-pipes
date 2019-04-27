//
//  Test.go
//  PureMVC Go Multicore Utility - Pipes
//
//  Copyright(c) 2019 Saad Shams <saad.shams@puremvc.org>
//  Your reuse is governed by the Creative Commons Attribution 3.0 License
//

package plumbing

import "github.com/puremvc/puremvc-go-util-pipes/src/interfaces"

type Callback struct {
	messagesReceived []interfaces.IPipeMessage // Array of received messages.
}

func (c *Callback) CallbackMethod(message interfaces.IPipeMessage) {
	c.messagesReceived = append(c.messagesReceived, message)
}

type Test struct {
	testVal int
}

type TestMessage struct {
	TestAtt string xml:"testAtt,attr"
}

type Rect struct {
	Width  float32
	Height float32
}

type Factor struct {
	factor int
}

type Bozo struct {
	level int
	name  string
}

type BozoThreshold struct {
	level int
}
