package main

import (
	"testing"
	"time"
)

func TestGoPool(t *testing.T) {
	pool := NewGoPool(2)

	// Start two jobs
	pool.StartJob()
	pool.StartJob()

	// Third job should block until one of the first two jobs is done
	go func() {
		pool.StartJob()
		t.Log("Third job started")
	}()

	// Wait for the third job to start
	time.Sleep(100 * time.Millisecond)

	// Finish one job
	pool.JobDone()

	// Wait for the third job to start
	time.Sleep(100 * time.Millisecond)

	// Finish the other job
	pool.JobDone()

	// Wait for the third job to start
	time.Sleep(100 * time.Millisecond)

	// Finish the third job
	pool.JobDone()
}
