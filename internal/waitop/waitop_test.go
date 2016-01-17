package waitop

import (
	"errors"
	"testing"
	"time"
)

var (
	errRejected = errors.New("rejected")
	errDoFail = errors.New("do failed")
)

func TestFulfill_goroutine(t *testing.T) {
	op := New()
	go func() {
		time.Sleep(time.Millisecond * 10)
		op.Fulfill("foo")
	}()
	r, err := op.Do(func() error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := r.(string); !ok || v != "foo" {
		t.Errorf("unexpected result: %v", r)
	}
}

func TestFulfill_AsyncOp(t *testing.T) {
	op := New()
	r, err := op.Do(func() error {
		op.Fulfill("foo")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := r.(string); !ok || v != "foo" {
		t.Errorf("unexpected result: %v", r)
	}
}

func TestFulfill_AsyncOp_goroutine(t *testing.T) {
	op := New()
	r, err := op.Do(func() error {
		go op.Fulfill("foo")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := r.(string); !ok || v != "foo" {
		t.Errorf("unexpected result: %v", r)
	}
}

func TestReject_goroutine(t *testing.T) {
	op := New()
	go func() {
		time.Sleep(time.Millisecond * 10)
		op.Reject(errRejected)
	}()
	r, err := op.Do(func() error {
		return nil
	})
	if err != errRejected {
		t.Errorf("unexpected error: %v", err)
	}
	if r != nil {
		t.Errorf("unexpected result: %v", r)
	}
}

func TestReject_AsyncOp(t *testing.T) {
	op := New()
	r, err := op.Do(func() error {
		op.Reject(errRejected)
		return nil
	})
	if err != errRejected {
		t.Errorf("unexpected error: %v", err)
	}
	if r != nil {
		t.Errorf("unexpected result: %v", r)
	}
}

func TestReject_AsyncOp_goroutine(t *testing.T) {
	op := New()
	r, err := op.Do(func() error {
		go op.Reject(errRejected)
		return nil
	})
	if err != errRejected {
		t.Errorf("unexpected error: %v", err)
	}
	if r != nil {
		t.Errorf("unexpected result: %v", r)
	}
}

func TestTerminated_goroutine(t *testing.T) {
	op := New()
	go func() {
		time.Sleep(time.Millisecond * 10)
		op.Close()
	}()
	r, err := op.Do(func() error {
		return nil
	})
	if err != ErrTerminated {
		t.Errorf("unexpected error: %v", err)
	}
	if r != nil {
		t.Errorf("unexpected result: %v", r)
	}
}

func TestTerminated_AsyncOp(t *testing.T) {
	op := New()
	r, err := op.Do(func() error {
		op.Close()
		return nil
	})
	if err != ErrTerminated {
		t.Errorf("unexpected error: %v", err)
	}
	if r != nil {
		t.Errorf("unexpected result: %v", r)
	}
}

func TestTerminated_AsyncOp_goroutine(t *testing.T) {
	op := New()
	r, err := op.Do(func() error {
		go op.Close()
		return nil
	})
	if err != ErrTerminated {
		t.Errorf("unexpected error: %v", err)
	}
	if r != nil {
		t.Errorf("unexpected result: %v", r)
	}
}

func TestDoing(t *testing.T) {
	op := New()
	go func() {
		op.Do(func() error {
			return nil
		})
	}()
	time.Sleep(time.Millisecond * 10)
	r, err := op.Do(func() error {
		return nil
	})
	if err != ErrAlreadyDoing {
		t.Errorf("unexpected error: %v", err)
	}
	if r != nil {
		t.Errorf("unexpected result: %v", r)
	}
	op.Close()
}

func TestNotDoing(t *testing.T) {
	op := New()
	err := op.Fulfill("foo")
	if err != ErrNotDoing {
		t.Errorf("unexpected error: %v", err)
	}
	err = op.Reject(errRejected)
	if err != ErrNotDoing {
		t.Errorf("unexpected error: %v", err)
	}
	op.Close()
}

func TestDo_Fail(t *testing.T) {
	op := New()
	r, err := op.Do(func() error {
		return errDoFail
	})
	if err != errDoFail {
		t.Errorf("unexpected error: %v", err)
	}
	if r != nil {
		t.Errorf("unexpected result: %v", r)
	}
	err = op.Fulfill("foo")
	if err != ErrNotDoing {
		t.Errorf("unexpected error: %v", err)
	}
}
