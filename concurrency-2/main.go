package main

import (
	"fmt"
	"log"
	"net/http"
)

var nextId = make(chan int)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>hello you got %v </h1>", <-nextId) // Reading value out of a channel. Save operation
	// unsafe data race. Increment operation is a read modify write. It would skip numbers.
	// nextId++
}

func counter() {
	for i := 0; ; i++ {
		nextId <- i //Sending numbers. In normal usecase of channel, we can't write to it unless it is ready to be read.
	}
}

// func main() {
// 	http.HandleFunc("/", handler)
// 	err := http.ListenAndServe(":3000", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
    go counter()
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
