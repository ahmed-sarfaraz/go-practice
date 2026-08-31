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
	"time"
)

type htmlType struct {
	html     []byte
	filename string
	err      error
}

func fetch(urls []string, channel chan<- htmlType) {
	var c = &http.Client{Timeout: 5 * time.Second}
	for _, url := range urls {
		go fetchHtml(url, c, channel)
	}
}

func fetchHtml(url string, c *http.Client, channel chan<- htmlType) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		channel <- htmlType{err: err}
		return
	}

	resp, err := c.Do(req)
	if err != nil {
		channel <- htmlType{err: err}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		channel <- htmlType{err: fmt.Errorf("%s: %s", url, resp.Status)}
		return
	}
	body, err := io.ReadAll(resp.Body)
	channel <- htmlType{html: body, err: err, filename: filepath.Base(url)}
}

func main() {
	flag.Parse()
	urls := flag.Args()
	channel := make(chan htmlType, len(urls))

	if len(urls) == 0 {
		log.Fatal("No URLs given to download")
	}

	fetch(urls, channel)

	for range len(urls) {
		html := <-channel

		if html.err != nil {
			log.Println(html.err)
			continue
		}

		dstFile, err := os.Create(html.filename)
		if err != nil {
			fmt.Printf("Error Creating file %s with error: %v\n", html.filename, err)
			continue
		}

		_, err = dstFile.Write(html.html)
		dstFile.Close()

		if err != nil {
			fmt.Printf("Error writing %s: %v\n", html.filename, err)
			continue
		}
	}

}
