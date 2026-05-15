package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Data struct {
	UserId, Id int
	Title      string
	Completed  bool
}

func main() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	data := string(body)
	fmt.Println(string(body))
	var test Data
	json.Unmarshal([]byte(data), &test)
	fmt.Printf("%v\n", test)
	http.HandleFunc("/", fooHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
