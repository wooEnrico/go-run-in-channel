package goRunInChannel

import (
	"errors"
	"log"
	"runtime/debug"
	"strings"
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

// defaultPanicHandler is the default panic handler that logs the panic and the parameter.
var defaultPanicHandler PanicHandler = PanicHandlerFunc(func(recovered any, param any) {
	// 格式化堆栈信息
	stackTrace := string(debug.Stack())
	// 添加缩进和分隔线
	formattedStack := strings.ReplaceAll(stackTrace, "\n", "\n    ")

	log.Printf("========== PANIC RECOVERED ==========\n"+
		"ERROR: %v\n"+
		"PARAM: %v\n"+
		"STACK TRACE:\n    %s\n"+
		"======================================",
		recovered, param, formattedStack)
})

// GoRunChannel manages concurrent execution of tasks with a limit on parallelism.
type GoRunChannel[T any] struct {
	channel      chan struct{}
	waitGroup    sync.WaitGroup
	panicHandler PanicHandler
	closeOnce    sync.Once
	closed       bool
}

// NewGoRunChannel creates a new GoRunChannel with a specified parallelism limit.
func NewGoRunChannel[T any](parallel int) *GoRunChannel[T] {
	if parallel <= 0 {
		parallel = 1
	}
	return &GoRunChannel[T]{
		channel:      make(chan struct{}, parallel),
		panicHandler: defaultPanicHandler,
		closed:       false,
	}
}

// NewGoRunChannelWithPanicHandler creates a new GoRunChannel with a specified parallelism limit and panic handler.
func NewGoRunChannelWithPanicHandler[T any](parallel int, panicHandler PanicHandler) *GoRunChannel[T] {
	if parallel <= 0 {
		parallel = 1
	}
	if panicHandler == nil {
		panicHandler = defaultPanicHandler
	}
	return &GoRunChannel[T]{
		channel:      make(chan struct{}, parallel),
		panicHandler: panicHandler,
		closed:       false,
	}
}

// WaitAndClose waits for all the goroutines to finish and close.
func (g *GoRunChannel[T]) WaitAndClose() {
	g.waitGroup.Wait()
	g.closeOnce.Do(func() {
		g.closed = true
		close(g.channel)
	})
}

// Run runs the provided runnable function with the given parameter.
func (g *GoRunChannel[T]) Run(runnable func(param T), param T) error {
	return g.RunWithRecover(runnable, param, g.panicHandler)
}

// RunWithRecover runs the provided runnable function with panic recovery.
func (g *GoRunChannel[T]) RunWithRecover(runnable func(param T), param T, handler PanicHandler) error {
	if runnable == nil {
		return errors.New("runnable function is nil")
	}
	if g.closed {
		return errors.New("channel is closed")
	}

	g.waitGroup.Add(1)
	g.channel <- struct{}{} // Acquire a slot

	go func() {
		defer func() {
			<-g.channel        // Release the slot
			g.waitGroup.Done() // Notify that this goroutine is done
		}()

		// Recover from panic to ensure Done is called
		defer func() {
			if r := recover(); r != nil && handler != nil {
				handler.Handle(r, param) // Call the panic handler
			}
		}()

		runnable(param)
	}()

	return nil
}
