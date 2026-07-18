package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	server := http.Server{Addr: ":8080", Handler: TextHandler("Hello world\r\n")}
	go server.ListenAndServe()
	req, err := http.NewRequestWithContext(context.TODO(), "GET", "http://localhost:8080", nil)
	if err != nil {
		log.Fatal("Error sending request")
	}
	resp, err := new(http.Client).Do(req)

	if err != nil {
		log.Fatal("Error reading response")
	}

	defer resp.Body.Close()

	resp.Write(os.Stdout)

}

type TextHandler string

var _ http.Handler = TextHandler("")

func (t TextHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) { w.Write([]byte(t)) }
