package main

import (
	"fmt"
	"github.com/tokuhirom/go-hsperfdata/hsperfdata"
	"log"
	"os"
)

func main() {
	repository, err := hsperfdata.New()
	if err != nil {
		log.Fatal(err)
	}

	pid := os.Args[1]
	file := repository.GetFile(pid)
	ch, err := file.ReadHsperfdata()
	if err != nil {
		log.Fatal("open fail", err)
	}

	for entry := range ch {
		fmt.Printf("%s=%v\n", entry.Key, entry.Value)
	}
}
