package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tokuhirom/go-hsperfdata/attach"
	"github.com/tokuhirom/go-hsperfdata/hsperfdata"
)

func main() {
	if len(os.Args) > 1 {
		pid, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("invalid pid: '%v'", os.Args[1])
		}

		sock, err := attach.New(pid)
		if err != nil {
			log.Fatal(err)
		}

		result, err := sock.Jcmd(strings.Join(os.Args[2:], " "))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	} else {
		repo, err := hsperfdata.New()
		if err != nil {
			log.Fatal("user", err)
		}
		files, err := repo.GetFiles()
		if err != nil {
			log.Fatal("repo", err)
		}

		for _, f := range files {
			result, err := f.Read()
			if err == nil {
				proc_name := result.GetProcName()
				fmt.Printf("%s %s\n", f.GetPid(), proc_name)
			} else {
				fmt.Printf("%s unknown\n", f.GetPid())
			}
		}
	}
}
