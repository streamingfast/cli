package cli

import (
	"os"

	"github.com/bobg/go-generics/v2/slices"
)

var globalExitManager = &exitManager{}

// Exit executes the registered exit handlers and then call `os.Exit(code)`.
// This can be used to implement a trap behavior. Of course, you must ensure
// that `os.Exit` are all done through `cli.Exit()` otherwise we will not be invoked.
//
// **Caveats** Panics are not recovered by those exit handlers, right now you need to implement your
// own trapping and then exit through `cli.Exit`.
//
// This library use `cli.Exit(code)` throughout so quitting due to `cli.NoError` or
// `cli.Ensure` will correctly call the exit handlers.
func Exit(code int) {
	globalExitManager.onExit(code)
	os.Exit(code)
}

// ExitHandler registers or unregisters an exit handler. If the `onExit` received
// is nil, unregister the handler with given `id`. Otherwise, register or update an
// existing one.
//
// No collision is checked, so an id overrides any previously existing id. This library
// is meant to be used on final CLI product, you should have low number of exit handler,
// pick your id and keep them short.
func ExitHandler(id string, onExit func(code int)) {
	globalExitManager.updateHandler(id, onExit)
}

type exitManager struct {
	Handlers []exitHandler
}

func (m *exitManager) onExit(code int) {
	for _, handler := range m.Handlers {
		handler.Handler(code)
	}
}

func (m *exitManager) updateHandler(id string, onExit func(code int)) {
	if onExit == nil {
		m.Handlers = slices.Filter(m.Handlers, func(handler exitHandler) bool { return handler.ID != id })
	} else {
		m.Handlers = append(m.Handlers, exitHandler{id, onExit})
	}
}

type exitHandler struct {
	ID      string
	Handler func(code int)
}
