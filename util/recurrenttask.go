/**
@file          recurrenttask.go
@package       util
@brief         Utility to schedule a recurring task.
@author        Edward Smith
@date          August 2019
@copyright     -©- Copyright © 2019 Edward Smith. All rights reserved. -©-
*/

package util

import (
	"time"
)

// RecurrentTask schedules a recurring task every `interval`. Stop by send `true` to the returned channel.
func RecurrentTask(interval time.Duration, task func()) chan bool {
	stop := make(chan bool)
	go func() {
		for {
			task()
			select {
			case <-time.After(interval):
			case <-stop:
				return
			}
		}
	}()
	return stop
}
