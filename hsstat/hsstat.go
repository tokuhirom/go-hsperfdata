package main

import (
	"flag"
	"fmt"
	"github.com/tokuhirom/go-hsperfdata/hsperfdata"
	"log"
	"os"
)

func newRepository(user string) (*hsperfdata.Repository, error) {
	if user == "" {
		return hsperfdata.New()
	} else {
		return hsperfdata.NewUser(user)
	}
}

func main() {
	version := flag.Bool("v", false, "show version")
	user := flag.String("u", "", "user")
	flag.Parse()

	if *version {
		fmt.Printf("%v\n", hsperfdata.GetVersion())
		return
	}

	if len(os.Args) == 1 {
		fmt.Printf("Usage: hsstat pid\n")
		return
	}

	pid := flag.Args()[0]

	repo, err := newRepository(*user)
	if err != nil {
		log.Fatal(err)
	}

	file := repo.GetFile(pid)
	result, err := file.Read()
	if err != nil {
		log.Fatal("open fail", err)
	}

	for key, value := range result.GetMap() {
		fmt.Printf("%s=%v\n", key, value)
	}
}
