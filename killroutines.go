/*
 * killroutines.go
 * Library to safely kill multiple goroutines
 * by J. Stuart McMurray
 * Created 20141127
 * Last modified 20141127
 *
 * Copyright (c) 2014 J. Stuart McMurray
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

// Package killroutines provides a one-time signal to signal multiple
// goroutines at once, by waiting for a channel close.  The benefit of using
// killroutines over a channel is killroutines provides built-in protection
// against closing a channel twice and the ability to test if the channel is
// already closed.
//
// The general pattern of usage is as follows.  A pointer is made that points
// to a struct of type K, which is passed to each goroutine to be signaled.
//     done := killroutines.New()
//     go worker(in chan int, out chan int, done *killroutines.K)
// Each goroutine waits on the channel returned by the Done() to close
//     select {
//     case i := <-in:
//             /* Do some work */
//
//             if allShouldExit() {
//                     done.Signal()
//             }
//
//     case <-done:
//             /* Cleanup */
//             return
//     }
package killroutines

import (
	"sync"
)

/* TODO: Examples */

/* K is an opaque struct passed to goroutines to allow for a one-time signal */
type K struct {
	c      chan int     /* Closed when all are to die */
	m      sync.RWMutex /* Held before changing c */
	closed bool         /* Whether C is closed or not */
}

/* New allocates and returns a pointer to a new K, which is is a ready-to-use state */
func New() *K {
	k := &K{}
	k.c = make(chan int)
	return k
}

//Signal safely closes the channel returned by Chan.  It is not an error for
//multiple goroutines to call Signal.  It returns true if it actually closed the
//channel (as opposed to it already being closed).  This is the means by which
//goroutines are signaled.
func (k *K) Signal() {
	/* Get a lock */
	k.m.Lock()
	defer k.m.Unlock()
	/* If it's open, close the channel, note it */
	if !k.closed {
		k.closed = true
		close(k.c)
	}
}

/* Chan returns a channel that will be closed when Signal is called. */
func (k *K) Chan() <-chan int {
	return k.c
}

//Closed returns true if the channel returned by Chan has been closed by a
//call to Signal.
func (k *K) Closed() bool {
	/* Get a lock */
	k.m.RLock()
	defer k.m.RUnlock()
	/* Return the closedness */
	return k.closed
}
