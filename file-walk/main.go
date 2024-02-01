package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// Build a program to find duplicate files based on their content.
// For this we have to use secure hash, because names/dates may differ.
type pair struct {
	hash string
	path string
}

type fileList []string

type results map[string]fileList

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

func fileWalkSeq(dir string) {
	if hashes, err := searchTree(os.Args[1]); err == nil {
		for hash, files := range hashes {
			if len(files) > 0 {
				fmt.Println(hash[len(hash)-7:], len(files))

				for _, file := range files {
					fmt.Println("  ", file)
				}
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing parameter, directory name is required!")
	}

	dir := os.Args[1]
	fileWalkSeq(dir)
}
