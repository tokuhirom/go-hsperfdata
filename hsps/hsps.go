package main

import (
	"flag"
	"fmt"
	"github.com/tokuhirom/go-hsperfdata/hsperfdata"
	"log"
)

func main() {
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Printf("%v\n", hsperfdata.GetVersion())
		return
	}

	repo, err := hsperfdata.New()
	if err != nil {
		log.Fatal("user", err)
	}
	files, err := repo.GetFiles()
	if err != nil {
		log.Fatal("repo", err)
	}

	for _, f := range files {
		proc_name, err := f.GetProcName()
		if err != nil {
			log.Fatal("procname", err)
		}
		fmt.Printf("%s %s\n", f.GetPid(), proc_name)
	}
}
