package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const url = "https://jsonplaceholder.typicode.com"

type todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	response, err := http.Get(url + "/todos/1")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	defer response.Body.Close()

	var item todo

	err = json.NewDecoder(response.Body).Decode(&item)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	fmt.Printf("%#v", item)
}
