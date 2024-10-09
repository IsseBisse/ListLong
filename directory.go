package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Directory struct {
	path      string
	id        uuid.UUID
	children  []Directory
	filesSize int64
	totalSize int64
}

func NewDirectory(path string, children []Directory) Directory {
	return Directory{path, uuid.New(), children, -1, -1}
}

func ToJob(directory Directory) Job {
	return Job{directory.id, directory.path, len(directory.children) == 0}
}

func ToJobs(directories []Directory) []Job {
	var jobs []Job

	for _, dir := range directories {
		job := ToJob(dir)
		jobs = append(jobs, job)

		if len(dir.children) > 0 {
			subJobs := ToJobs(dir.children)
			jobs = append(jobs, subJobs...)
		}
	}

	// TODO: Look at sorting jobs to get longest running jobs first
	// using some heuristic
	return jobs
}

func AddResults(directories []Directory, results map[uuid.UUID]int64) ([]Directory, int64) {
	var totalSize int64 = 0
	for idx := 0; idx < len(directories); idx++ {
		directories[idx].filesSize = results[directories[idx].id]

		if len(directories[idx].children) > 0 {
			// TODO: Passing down results could probably be handled nicer
			children, childrenSize := AddResults(directories[idx].children, results)
			directories[idx].children = children
			directories[idx].totalSize = directories[idx].filesSize + childrenSize

		} else {
			directories[idx].totalSize = directories[idx].filesSize
		}
		totalSize += directories[idx].totalSize
	}

	return directories, totalSize
}

func Subdirectories(root string, depth int) []Directory {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	var directories []Directory
	for _, e := range entries {
		if e.IsDir() {
			path := filepath.Join(root, e.Name())

			var subDirs []Directory
			if depth > 1 {
				subDirs = Subdirectories(path, depth-1)
			}
			dir := NewDirectory(path, subDirs)
			directories = append(directories, dir)
		}
	}

	return directories
}

var SizePrefix = [...]string{"", "k", "M", "G", "T"}

func ToString(directory Directory) string {
	parts := strings.Split(directory.path, "\\")
	dirName := parts[len(parts)-1]

	size := float32(directory.totalSize)
	idx := 0
	for size > 1024.0 {
		size = size / 1024.0
		idx++
	}
	sizeString := fmt.Sprintf("%.1f %sB", size, SizePrefix[idx])

	return fmt.Sprintf("%s: %s", dirName, sizeString)
}

func Print(directories []Directory, level int) {
	// TODO: Improve this using "tree symbols" like: ├──, └──
	for _, dir := range directories {
		str := ""
		for idx := 0; idx < level; idx++ {
			str += "  "
		}
		str += ToString(dir)
		fmt.Printf("%s\n", str)

		if len(dir.children) > 0 {
			Print(dir.children, level+1)
		}
	}
}
