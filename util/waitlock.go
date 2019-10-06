/**
@file          waitlock.go
@package       util
@brief         A broadcast wait lock with timeout.
@author        Edward Smith
@date          August 2019
@copyright     -©- Copyright © 2019 Edward Smith. All rights reserved. -©-
*/

package util

import (
	"errors"
	"time"
)

var (
	// ErrTimeout indicates that waiting for a CI job has timed out.
	ErrTimeout = errors.New("timed out")
)

// WaitLock is a lock timeout lock that multiple processes can wait for a broadcast signal.
type WaitLock struct {
	c        chan error
	happened bool
}

// NewWaitLock returns a new WaitLock.
func NewWaitLock() *WaitLock {
	w := new(WaitLock)
	w.c = make(chan error)
	return w
}

// Wait waits until the timeout or the lock is signal.
func (w *WaitLock) Wait(d time.Duration) error {
	timeout := make(chan error)
	go func() {
		time.Sleep(d)
		timeout <- ErrTimeout
	}()
	var err error
	select {
	case err = <-timeout:
	case err = <-w.c:
	}
	return err
}

// Signal signals all processes waiting on WaitLock to resume.
func (w *WaitLock) Signal() {
	close(w.c)
}
