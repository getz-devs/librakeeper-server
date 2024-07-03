package main

import (
	"github.com/getz-devs/librakeeper-server/internal/searcher-agent/rabbit"
)

func main() {
	_, err := rabbit.ScrapISBNFindBook("9785206000344")
	if err != nil {
		panic(err)
	}
}
