/*
Package backoff provides retry backoff algorithm.
*/
package backoff

import "time"

// Exp provides exponential back off.
type Exp struct {
	Min time.Duration
	Max time.Duration

	count uint32
}

// Wait sleeps using exponential back off.
func (exp *Exp) Wait() {
	d := exp.min() * (1 << exp.count)
	if m := exp.max(); d > m {
		d = m
	}
	time.Sleep(d)
	if exp.count < 31 {
		exp.count++
	}
}

// Reset resets exponential count.
func (exp *Exp) Reset() {
	exp.count = 0
}

func (exp *Exp) min() time.Duration {
	if exp.Min <= 0 {
		return time.Millisecond
	}
	return exp.Min
}

func (exp *Exp) max() time.Duration {
	if exp.Max <= 0 {
		return time.Second
	}
	return exp.Max
}
