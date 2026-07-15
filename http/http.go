package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Data struct {
	UserId, Id int
	Title      string
	Completed  bool
}

func getMethod(url string) *http.Request {
	ctx := context.TODO()
	const method string = "GET"
	var body io.Reader = nil
	req, err := http.NewRequestWithContext(ctx, method, url, body)

	if err != nil {
		log.Fatal(err)
	}

	return req
}

func getMethodWithParameters() {
	ctx := context.TODO()
	const method string = "GET"
	const path string = "https://scryfall.com/search"
	var body io.Reader = nil
	v := make(url.Values)
	v.Add("q", `"of Emrakul"`)
	v.Add("order", "released")
	v.Add("direction", "asc")
	dst := path + "?" + v.Encode()

	req, err := http.NewRequestWithContext(ctx, method, dst, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Write(os.Stdout)
}

func writeResponseToLocalFile() {
	log.SetPrefix("writeResponseToLocalFile: ")
	dir := flag.String("dir", ".", "where the file will be created")
	timeout := flag.Duration("timeout", 30*time.Second, "timeout for download")
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		log.Fatal("usage: ./http [-timeout duration] url filename")
	}

	url, filename := args[0], args[1]
	c := http.Client{Timeout: *timeout}
	req := getMethod(url)
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(fmt.Errorf("request: %v", err))
	}

	if res.StatusCode != http.StatusOK {
		log.Fatal(fmt.Errorf("request with code: %d", res.StatusCode))
	}
	defer res.Body.Close()
	dstPath := filepath.Join(*dir, filename)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		log.Fatal(fmt.Errorf("Error creating file %v", err))
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, res.Body); err != nil {
		log.Fatal(fmt.Errorf("copying response to file: %v", err))
	}
}

func main() {
	writeResponseToLocalFile()
}
