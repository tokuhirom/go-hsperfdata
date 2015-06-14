package main

import (
	"flag"
	"fmt"
	"github.com/tokuhirom/go-hsperfdata/hsperfdata"
	"log"
	"os"
)

func main() {
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Printf("%v\n", hsperfdata.GetVersion())
		return
	}

	if len(os.Args) == 1 {
		fmt.Printf("Usage: hsstat pid\n")
		return
	}

	pid := os.Args[1]

	repository, err := hsperfdata.New()
	if err != nil {
		log.Fatal(err)
	}

	file := repository.GetFile(pid)
	ch, err := file.ReadHsperfdata()
	if err != nil {
		log.Fatal("open fail", err)
	}

	for entry := range ch {
		fmt.Printf("%s=%v\n", entry.Key, entry.Value)
	}
}
