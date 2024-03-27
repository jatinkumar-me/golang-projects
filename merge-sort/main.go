package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
    fmt.Println("Number of CPU", runtime.NumCPU())
    fmt.Println("GOMAXPROC", runtime.GOMAXPROCS(1))

	arr := generateRandomSlice(10_000_000, -10, 1000)
	arr2 := make([]int, len(arr))
	arr3 := make([]int, len(arr))
	copy(arr2, arr)
	copy(arr3, arr)

    runWithTime(arr, mergeSort)
    // runWithTime(arr2, mergeSortConcur)
    // runWithTime(arr3, mergeSortConcurV2)
}

func runWithTime(arr []int, cb func([]int)) {
	startTime := time.Now()
	cb(arr)
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Time taken: %d\n", elapsedTime.Milliseconds())
}

func merge(s []int, mid int) {
	leftEnd := mid - 1
	rightStart := mid
	rightEnd := len(s) - 1

	temp := make([]int, len(s))
	i, j, k := 0, rightStart, 0

	for i <= leftEnd && j <= rightEnd {
		if s[i] <= s[j] {
			temp[k] = s[i]
			i++
		} else {
			temp[k] = s[j]
			j++
		}
		k++
	}

	for i <= leftEnd {
		temp[k] = s[i]
		i++
		k++
	}

	for j <= rightEnd {
		temp[k] = s[j]
		j++
		k++
	}

	for i, val := range temp {
		s[i] = val
	}
}

// Sequential approach
func mergeSort(s []int) {
	if len(s) <= 1 {
		return
	}
	middle := len(s) / 2
	mergeSort(s[:middle])
	mergeSort(s[middle:])
	merge(s, middle)
}

// Concurrent approach
func mergeSortConcur(s []int) {
	if len(s) <= 1 {
		return
	}
	middle := len(s) / 2

	mergeSort(s[:middle])

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		mergeSort(s[middle:])
	}()
	wg.Wait()

	merge(s, middle)
}

const MAX_CONCUR_LEN = 1 << 8

// Concurrent approach
func mergeSortConcurV2(s []int) {
	len := len(s)
	if len <= 1 {
		return
	}

	if len < MAX_CONCUR_LEN {
		// Proceed sequentially
		mergeSort(s)
		return
	}

	middle := len / 2

	mergeSort(s[:middle])

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		mergeSort(s[middle:])
	}()
	wg.Wait()

	merge(s, middle)
}

// generates a slice of random ints
func generateRandomSlice(length, min, max int) []int {
	if length <= 0 || min > max {
		return nil
	}

	result := make([]int, length)
	rangeSize := max - min + 1

	for i := 0; i < length; i++ {
		result[i] = rand.Intn(rangeSize) + min
	}

	return result
}
