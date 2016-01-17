package waitop

import (
	"errors"
	"sync"
)

type AsyncOp func() error

type WaitOp interface {
	Do(AsyncOp) (interface{}, error)
	Fulfill(interface{})
	Reject(error)
	Close()
}

type waitOp struct {
	c   *sync.Cond
	s   state
	v   interface{} // return value
	err error
}

var (
	ErrAlreadyDoing = errors.New("already doing")
	ErrTerminated = errors.New("terminated")
)

type state int

const (
	idle state = iota
	doing
	done
)

func New() WaitOp {
	return &waitOp{
		c: sync.NewCond(new(sync.Mutex)),
	}
}

func (w *waitOp) Do(f AsyncOp) (interface{}, error) {
	// mark as doing.
	w.c.L.Lock()
	if w.s != idle {
		w.c.L.Unlock()
		return nil, ErrAlreadyDoing
	}
	w.s, w.v, w.err = doing, nil, nil
	w.c.L.Unlock()
	// start async operation.
	err := f()
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

func (w *waitOp) done(v interface{}, err error) {
	w.c.L.Lock()
	if w.s == doing {
		w.s = done
		w.v = v
		w.err = err
		w.c.Signal()
	}
	w.c.L.Unlock()
}

func (w *waitOp) Fulfill(retval interface{}) {
	w.done(retval, nil)
}

func (w *waitOp) Reject(err error) {
	w.done(nil, err)
}

func (w *waitOp) Close() {
	w.done(nil, ErrTerminated)
}
