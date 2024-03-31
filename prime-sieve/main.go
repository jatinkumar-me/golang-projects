// A concurrent programming for generating prime numbers using goroutines and channels
package main

import "fmt"

// generator sends numbers from 2 to limit to src channel
func generator(limit int, src chan<- int) {
	for i := 2; i < limit; i++ {
		src <- i
	}

	close(src)
}

func filter(src <-chan int, dest chan<- int, prime int) {
	for i := range src {
		if i%prime != 0 {
			dest <- i
		}
	}

	close(dest)
}

func sieve(limit int) {
	src := make(chan int)
	go generator(limit, src)

	for {
		prime, ok := <-src
		if !ok {
			break
		}

		newSrc := make(chan int)
		go filter(src, newSrc, prime)

		src = newSrc

		fmt.Print(prime, " ")
	}
}

func main() {
	sieve(100)
}
