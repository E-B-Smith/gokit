/**
@file          waitlock_test.go
@package       util
@brief         Tests wait functions.
@author        Edward Smith
@date          July 2019
@copyright     -©- Copyright © 2019 Edward Smith. All rights reserved. -©-
*/

package util

import (
	"sync"
	"testing"
	"time"
)

func TestWaitLockTimeout(t *testing.T) {
	w := NewWaitLock()
	n := time.Now()
	e := w.Wait(time.Second)
	d := time.Since(n)
	if e != ErrTimeout {
		t.Errorf("Expected ErrTimeout. Got %v.", e)
	}
	if d < time.Second || d > time.Second*2 {
		t.Errorf("Expected one second timeout. Got %f.", d.Seconds())
	}
}

func TestWaitLockSignal(t *testing.T) {
	w := NewWaitLock()
	go func() {
		time.Sleep(time.Second)
		w.Signal()
	}()
	n := time.Now()
	e := w.Wait(time.Second * 2)
	d := time.Since(n)
	if e != nil {
		t.Errorf("Expected nil got %v.", e)
	}
	if d >= 2*time.Second || d < time.Second {
		t.Errorf("Expected one second but got %f.", d.Seconds())
	}
}

func TestWaitLockRace(t *testing.T) {
	w := NewWaitLock()
	w.Signal()
	n := time.Now()
	e := w.Wait(time.Second * 2)
	d := time.Since(n)
	if e != nil {
		t.Errorf("Expected nil got %v.", e)
	}
	if d > time.Millisecond*10 || d < 0 {
		t.Errorf("Expected < 10ms but got %d.", d)
	}
}

func TestMultiWait(t *testing.T) {
	var waitgroup = new(sync.WaitGroup)
	waitgroup.Add(2)
	w := NewWaitLock()
	go func() {
		defer waitgroup.Done()
		w.Wait(time.Second * 2)
		// println("Done 1")
	}()
	go func() {
		defer waitgroup.Done()
		w.Wait(time.Second * 2)
		// println("Done 2")
	}()
	n := time.Now()
	time.Sleep(1 * time.Second)
	w.Signal()
	waitgroup.Wait()
	d := time.Since(n)
	if d >= 2*time.Second || d < time.Second {
		t.Errorf("Expected one second but got %d.", d)
	}
}
