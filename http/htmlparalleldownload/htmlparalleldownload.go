package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func downloadAndSave(src string, dest string) error {
	req, err := http.NewRequestWithContext(context.TODO(), "GET", src, nil)
	if err != nil {
		return err
	}

	res, err := new(http.Client).Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	dstFile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer dstFile.Close()
	if _, err = io.Copy(dstFile, res.Body); err != nil {
		return err
	}

	return nil
}

func main() {
	var dstDir string
	flag.String("directory", ".", "Location of the files to be stored at")
	flag.Duration("timeout", 2*time.Second, "Timeout for each request")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("No URL's given to download the HTML from")
	}

	dstDir, err := filepath.Abs(dstDir)
	if err != nil {
		log.Fatal("Check the Storage location: %v", err)
	}

	dst := make([]string, len(args))
	for i := range args {
		dst[i] = filepath.Join(dstDir, filepath.Base(args[i]))
	}

	errs := make([]error, len(args))
	wg := new(sync.WaitGroup)
	wg.Add(len(args))
	t := time.Now()
	for i := range args {
		i := i
		go func() {
			defer wg.Done()
			errs[i] = downloadAndSave(args[i], dst[i])
		}()
	}

	wg.Wait()

	fmt.Printf("Time to download concurrently was: %v\n", time.Since(t))

	errCount := 0
	for i := range errs {
		if errs[i] != nil {
			errCount += 1
		}
	}

	os.Exit(errCount)
}
