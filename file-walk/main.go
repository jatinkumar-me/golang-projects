package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// Build a program to find duplicate files based on their content.
// For this we have to use secure hash, because names/dates may differ.
type pair struct {
	hash string
	path string
}

type fileList []string

type results map[string]fileList

// SEQUENTIAL APPROACH ================================================================

func hashFile(path string) pair {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	return pair{fmt.Sprintf("%x", hash.Sum(nil)), path}
}

func searchTree(dir string) (results, error) {
	hashes := make(results)

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			hashPair := hashFile(path)
			hashes[hashPair.hash] = append(hashes[hashPair.hash], hashPair.path)
		}

		return nil
	})

	return hashes, err
}

func doFileWalk(dir string) {
	if hashes, err := fileWalkConcur(os.Args[1]); err == nil {
		for hash, files := range hashes {
			if len(files) > 1 {
				fmt.Println(hash[len(hash)-7:], len(files))

				for _, file := range files {
					fmt.Println("  ", file)
				}
			}
		}
	}
}

// CONCURRENT APPROACH ================================================================

// Use a fixed pool of goroutines and a collector and channels

func collectHashes(pairs <-chan pair, result chan<- results) {
	// fmt.Println("collecting hashes")
	hashes := make(results)

	for p := range pairs {
		// fmt.Println("received pair", p)
		hashes[p.hash] = append(hashes[p.hash], p.path)
	}

	result <- hashes
}

func processFiles(paths <-chan string, pairs chan<- pair, done chan<- bool) {
	for path := range paths {
		pairs <- hashFile(path)
	}

	done <- true
}

func fileWalkConcur(dir string) (results, error) {
	workers := 2 * runtime.GOMAXPROCS(0)
	fmt.Println("Main program started with", workers, "workers")

	paths := make(chan string)
	pairs := make(chan pair)
	done := make(chan bool)
	result := make(chan results)

	for i := 0; i < workers; i++ {
		go processFiles(paths, pairs, done)
	}

	go collectHashes(pairs, result)

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsRegular() && info.Size() > 0 {
			paths <- path
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	close(paths)

	for i := 0; i < workers; i++ {
		<-done
	}

	close(pairs)

	hashes := <-result
	return hashes, nil
}

// CONCURRENT APPROACH 2 =============================================================
// multi-threaded walk of the directory tree;

func fileWalkConcur2(dir string) (results, error) {
	workers := 2 * runtime.GOMAXPROCS(0)

	paths := make(chan string)
	pairs := make(chan pair)
	done := make(chan bool)
	result := make(chan results)

	for i := 0; i < workers; i++ {
		go processFiles(paths, pairs, done)
	}

	go collectHashes(pairs, result)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	err := walkDir(dir, paths, wg)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	close(paths)

	for i := 0; i < workers; i++ {
		<-done
	}

	close(pairs)

	hashes := <-result
	return hashes, nil
}

func walkDir(dir string, paths chan<- string, wg *sync.WaitGroup) error {
	defer wg.Done()
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsDir() && path != dir {
			wg.Add(1)
			go walkDir(path, paths, wg)
			return filepath.SkipDir
		}

		if info.Mode().IsRegular() && info.Size() > 0 {
			paths <- path
		}

		return nil
	})
	return err
}

// CONCURRENT APPROACH 2 =============================================================

func fileWalkConcur3(dir string) (results, error) {
	workers := 2 * runtime.GOMAXPROCS(0)

	wg := new(sync.WaitGroup)

	pairs := make(chan pair, workers)
	result := make(chan results)
	limits := make(chan bool, workers)

	for i := 0; i < workers; i++ {
		go processFiles(paths, pairs, done)
	}

	go collectHashes(pairs, result)

	wg.Add(1)
	err := walkDir2(dir, pairs, wg, limits)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	close(paths)

	for i := 0; i < workers; i++ {
		<-done
	}

	close(pairs)

	hashes := <-result
	return hashes, nil
}

func walkDir2(dir string, pairs chan<- string, wg *sync.WaitGroup) error {
	defer wg.Done()
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsDir() && path != dir {
			wg.Add(1)
			go walkDir2(path, paths, wg)
			return filepath.SkipDir
		}

		if info.Mode().IsRegular() && info.Size() > 0 {
			paths <- path
		}

		return nil
	})
}

// ===================================================================================

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing parameter, directory name is required!")
	}

	dir := os.Args[1]
	doFileWalk(dir)
}
