/**
@file          recurrenttask_test.go
@package       util
@brief         Tests the recurrenttask functions.
@author        Edward Smith
@date          July 2019
@copyright     -©- Copyright © 2019 Edward Smith. All rights reserved. -©-
*/

package util

import (
	"testing"
	"time"
)

func TestRecurrentTask(t *testing.T) {
	var pings string
	ping := func() { pings += "#" }
	stop := RecurrentTask(100*time.Millisecond, ping)
	time.Sleep(500 * time.Millisecond)
	if pings != "#####" {
		t.Errorf("Expected '#####' got '%s'.", pings)
	}
	stop <- true
	time.Sleep(200 * time.Millisecond)
	if pings != "#####" {
		t.Errorf("Expected '#####' got '%s'.", pings)
	}
}
