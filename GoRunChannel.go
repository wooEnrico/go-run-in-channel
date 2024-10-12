package goRunInChannel

import (
	"sync"
)

type GoRunChannel struct {
	channel   chan struct{}
	waitGroup sync.WaitGroup
}

// NewGoRunChannel creates a new GoRunChannel with the parallel parameter
func NewGoRunChannel(parallel int) *GoRunChannel {
	return &GoRunChannel{
		channel: make(chan struct{}, parallel),
	}
}

// Wait waits for all the goroutines to finish
func (GoRunChannel *GoRunChannel) Wait() {
	GoRunChannel.waitGroup.Wait()
}

// Run runs the runnable function with the param parameter
func (GoRunChannel *GoRunChannel) Run(runnable func(param interface{}), param interface{}) {
	GoRunChannel.waitGroup.Add(1)
	GoRunChannel.channel <- struct{}{}
	go running(GoRunChannel, runnable, param)
}

// running is a goroutine that runs the runnable function
func running(GoRunChannel *GoRunChannel, runnable func(param interface{}), param interface{}) {
	defer GoRunChannel.waitGroup.Done()
	runnable(param)
	<-GoRunChannel.channel
}
