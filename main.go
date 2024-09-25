package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/urfave/cli"
)

type Printer func(string)

func listDirectorySize(root string, depth int, verbose bool) {
	var printer Printer
	if verbose {
		printer = timestampPrintln
	} else {
		printer = func(message string) {}
	}

	directories := Subdirectories(root, depth)
	jobs := ToJobs(directories)

	numJobs := len(jobs)
	jobsChannel := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup
	numWorkers := 8
	for wId := 0; wId < numWorkers; wId++ {
		wg.Add(1)
		go worker(jobsChannel, results, &wg, printer)
	}

	for _, job := range jobs {
		jobsChannel <- job
	}
	close(jobsChannel)

	// Prevent dead-lock
	go func() {
		wg.Wait()
		close(results)
	}()

	resultsMap := make(map[uuid.UUID]int64)
	for res := range results {
		resultsMap[res.id] = res.size
		printer(fmt.Sprintf("%x finished. Size=%d", res.id[:4], res.size))
	}

	directories, _ = AddResults(directories, resultsMap)

	Print(directories, 0)
}

var StartTime = time.Now()

func main() {
	var root string
	var depth int
	var verbose bool

	app := &cli.App{
		Name:  "lsl",
		Usage: "windows replica of shell command 'ls -l'",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Value:       ".",
				Usage:       "root path to start from",
				Destination: &root,
			},
			&cli.IntFlag{
				Name:        "depth",
				Value:       1,
				Usage:       "search depth",
				Destination: &depth,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Usage:       "use verbose mode",
				Destination: &verbose,
			},
		},
		Action: func(cCtx *cli.Context) error {
			listDirectorySize(root, depth, verbose)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	// // For debugging
	// root := "C:\\Users\\isak.liljequist\\OneDrive - CGI\\Documents\\Uppdrag\\Rottneros - Vallvik"
	// depth := 2
	// verbose := true
	// ListDirectorySize(root, depth, verbose)
}
