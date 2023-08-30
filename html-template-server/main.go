package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

type todo struct {
	UserID    int    `json:"userID"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

const base = "https://jsonplaceholder.typicode.com/"

var form = `
<h1>Todo #{{.ID}}</h1>
<div>{{printf "User %d" .UserID}}</div>
<div>{{printf "%s (completed: %t)" .Title .Completed}}</div>`

func handler(writer http.ResponseWriter, request *http.Request) {
	var item todo

	response, err := http.Get(base + request.URL.Path[1:])
	if response.StatusCode != http.StatusOK {
		http.NotFound(writer, request)
		return
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&item)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	template := template.New("my-template")

	template.Parse(form)
	template.Execute(writer, item)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
