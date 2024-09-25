package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func size(dir string, recursive bool) (int64, error) {
	var size int64
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		} else if path != dir && !recursive {
			return filepath.SkipDir
		}
		return nil
	})
	return size, err
}

func timestampPrintln(message string) {
	now := time.Now()
	duration := now.Sub(StartTime).Round(time.Millisecond).String()
	fmt.Printf("%s: %s\n", duration, message)
}
