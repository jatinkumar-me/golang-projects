package main

import (
	"log"
	"time"
)

// Select allows any "ready" alternative to proceed among
// - a channel we can read from
// - a channel we can write to
// - a deafult action that's always ready

func main() {
	chans := []chan int{
		make(chan int),
		make(chan int),
	}

	for i := range chans {
		go func(i int, ch chan<- int) {
			for {
				time.Sleep(time.Duration(i) * time.Second)
				ch <- i
			}
		}(i+1, chans[i])
	}

	// for i := 0; i < 12; i++ {
 //        // Prints whichever comes first
	// 	select {
	// 	case m0 := <-chans[0]:
	// 		log.Println("received0", m0)
	// 	case m1 := <-chans[1]:
	// 		log.Println("received1", m1)
	// 	}
	// }

    for i := 0; i < 12; i++ {
        m0 := <- chans[0]
        log.Println("received0", m0)

        m1 := <- chans[1]
        log.Println("received1", m1)
	}
}
