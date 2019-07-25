package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func parseFile(fileName string, words chan string, wgWord *sync.WaitGroup) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// https://www.dotnetperls.com/file-go
	// http://networkbit.ch/read-text-file-in-go/
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		words <- word
	}
	wgWord.Done()
}

func counterWords(words chan string, numWord map[string]int) {
	for word := range words {
		if _, ok := numWord[word]; ok {
			numWord[word]++
		} else {
			numWord[word] = 1
		}
	}
}

func main() {
	dir := "./check"
	words := make(chan string)
	numWord := make(map[string]int)
	var wgWord sync.WaitGroup

	if len(os.Args[1:]) > 0 {
		dir = os.Args[1]
	}

	// read files name on directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		wgWord.Add(1)
		go parseFile(dir+"/"+file.Name(), words, &wgWord)
	}

	go counterWords(words, numWord)
	wgWord.Wait()

	fmt.Println(numWord)
}
