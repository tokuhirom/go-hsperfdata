package main

import (
	"../hsperfdata"
	"fmt"
	"log"
)

func main() {
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
