package waitop

import (
	"errors"
	"sync"
)

var (
	// ErrAlreadyDoing indicates WaitOp is executing already.
	ErrAlreadyDoing = errors.New("already doing")

	// ErrTerminated indicates WaitOp have been closed without Fulfill() or
	// Rejected().
	ErrTerminated = errors.New("terminated")

	// ErrNotDoing indicates asynchronous operation is not started yet.
	ErrNotDoing = errors.New("not started")
)

// AsyncOp starts asynchronous operation.
// When AsyncOp is finished Fulfill() or Reject() should be called.
type AsyncOp func() error

type state int

const (
	idle state = iota
	doing
	done
)

// WaitOp provides wait for exclusive asynchronous operation.
type WaitOp struct {
	c   *sync.Cond
	s   state
	v   interface{} // return value
	err error
}

// New creates a WaitOp instance.
func New() *WaitOp {
	return &WaitOp{
		c: sync.NewCond(new(sync.Mutex)),
	}
}

// Do starts asynchronous operation if it's not started yet.
// ErrAlreadyDoing and ErrTerminated will be retured.
func (w *WaitOp) Do(f AsyncOp) (interface{}, error) {
	// mark as doing.
	w.c.L.Lock()
	if w.s != idle {
		w.c.L.Unlock()
		return nil, ErrAlreadyDoing
	}
	w.s, w.v, w.err = doing, nil, nil
	w.c.L.Unlock()
	// start async operation.
	var err error
	if f != nil {
		err = f()
	}
	if err != nil {
		w.c.L.Lock()
		w.s, w.v, w.err = idle, nil, nil
		w.c.L.Unlock()
		return nil, err
	}
	// wait until done.
	w.c.L.Lock()
	for w.s != done {
		w.c.Wait()
	}
	v, err := w.v, w.err
	w.s, w.v, w.err = idle, nil, nil
	w.c.L.Unlock()
	return v, err
}

// TODO: support DoWithContext()

func (w *WaitOp) done(v interface{}, err error) error {
	rerr := ErrNotDoing
	w.c.L.Lock()
	if w.s == doing {
		rerr = nil
		w.s = done
		w.v = v
		w.err = err
		w.c.Signal()
	}
	w.c.L.Unlock()
	return rerr
}

// Fulfill should be called with result value if asynchronous operation is
// finished successfully.
// When WaitOp is not started yet, this returns ErrNotDoing.
func (w *WaitOp) Fulfill(retval interface{}) error {
	return w.done(retval, nil)
}

// Reject should be called if asynchronous operation is failed.
// When WaitOp is not started yet, this returns ErrNotDoing.
func (w *WaitOp) Reject(err error) error {
	return w.done(nil, err)
}

// Close closes WaitOp.  When executing WaitOp is closed, it returns
// ErrTerminated.
func (w *WaitOp) Close() error {
	w.done(nil, ErrTerminated)
	return nil
}
