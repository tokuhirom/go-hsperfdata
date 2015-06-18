package main

import (
	"flag"
	"fmt"
	"github.com/tokuhirom/go-hsperfdata/hsperfdata"
	"log"
)

func newRepository(user string) (*hsperfdata.HsperfdataRepository, error) {
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

	repo, err := newRepository(*user)
	if err != nil {
		log.Fatal("user", err)
	}
	files, err := repo.GetFiles()
	if err != nil {
		log.Fatal("repo", err)
	}

	for _, f := range files {
		result, err := f.Read()
		if err != nil {
			proc_name := result.GetProcName()
			fmt.Printf("%s %s\n", f.GetPid(), proc_name)
		} else {
			fmt.Printf("%s unknown\n", f.GetPid())
		}
	}
}
