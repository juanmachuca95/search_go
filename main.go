package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	fs, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}

	dictionaries := make([]string, 0, len(fs)-1)
	for _, f := range fs {
		if f.Name() != "main.go" {
			dictionaries = append(dictionaries, f.Name())
		}
	}

	// we have to wait for at least one goruntine to find the word.
	wg := sync.WaitGroup{}
	wg.Add(1)

	// channel to receive the dictionary
	out := make(chan string)
	// create context to cancel gorutins
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context, wg *sync.WaitGroup) {
		for {
			select {
			case <-ctx.Done():
				log.Println("Close another goruntines")
				return
			case dictionary := <-out:
				wg.Done()
				log.Println("WINNER DICTIONARY: ", dictionary)
				cancel()
			}
		}
	}(ctx, &wg)

	for _, p := range dictionaries {
		go Search(ctx, "andrew", p, out) // Set path
	}

	wg.Wait()
	log.Println("Program finished")
}

func Search(ctx context.Context, target, path string, out chan string) {
	readFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("cannot open file on path %s - error %v", path, err)
		return
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		// time.Sleep(time.Second * 1) // You can uncomment this linea to see this process in more details

		// add sleep to simulate a to too large file
		log.Println("[+] - ", path, " Palabra: ", fileScanner.Text())
		if fileScanner.Text() == target {
			out <- path
			break
		}
	}
}
