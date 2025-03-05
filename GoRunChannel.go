package goRunInChannel

import (
	"log"
	"sync"
)

// PanicHandler defines an interface for handling panics in goroutines.
type PanicHandler interface {
	Handle(recovered any, param any)
}

// PanicHandlerFunc is a function type that implements the PanicHandler interface.
type PanicHandlerFunc func(recovered any, param any)

// Handle calls the PanicHandlerFunc itself.
func (f PanicHandlerFunc) Handle(recovered any, param any) {
	f(recovered, param)
}

// GoRunChannel manages concurrent execution of tasks with a limit on parallelism.
type GoRunChannel[T any] struct {
	channel   chan struct{}
	waitGroup sync.WaitGroup
}

// NewGoRunChannel creates a new GoRunChannel with a specified parallelism limit.
func NewGoRunChannel[T any](parallel int) *GoRunChannel[T] {
	return &GoRunChannel[T]{
		channel: make(chan struct{}, parallel),
	}
}

// Wait waits for all the goroutines to finish.
func (g *GoRunChannel[T]) Wait() {
	g.waitGroup.Wait()
}

// Run runs the provided runnable function with the given parameter.
func (g *GoRunChannel[T]) Run(runnable func(param T), param T) {
	g.RunWithRecover(runnable, param, PanicHandlerFunc(func(recovered any, param any) {
		log.Printf("Panic recovered: %v Param: %v\n", recovered, param)
	}))
}

// RunWithRecover runs the provided runnable function with panic recovery.
func (g *GoRunChannel[T]) RunWithRecover(runnable func(param T), param T, handler PanicHandler) {
	g.waitGroup.Add(1)
	g.channel <- struct{}{} // Acquire a slot

	go func() {
		defer func() {
			<-g.channel        // Release the slot
			g.waitGroup.Done() // Notify that this goroutine is done
		}()

		// Recover from panic to ensure Done is called
		defer func() {
			if r := recover(); r != nil {
				handler.Handle(r, param) // Call the panic handler
			}
		}()

		runnable(param)
	}()
}
