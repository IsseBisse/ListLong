package main

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Job struct {
	id                uuid.UUID
	path              string
	isBottomDirectory bool
}

type Result struct {
	id   uuid.UUID
	size int64
}

func Worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, printer Printer) {
	defer wg.Done()
	for job := range jobs {
		printer(fmt.Sprintf("%x started", job.id[:4]))

		size, _ := size(job.path, job.isBottomDirectory)
		results <- Result{job.id, size}
	}
	printer("worker done!")
}
