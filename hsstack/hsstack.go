package main

import (
	"fmt"
	"github.com/tokuhirom/go-hsperfdata/attach"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: hsstack pid\n")
		os.Exit(1)
	}
	pid := os.Args[1]
	sock, err := attach.New(pid)
	if err != nil {
		log.Fatalf("cannot open unix socket: %s", err)
	}
	err = sock.Execute("threaddump")
	if err != nil {
		log.Fatalf("cannot write to unix socket: %s", err)
	}

	stack, err := sock.ReadString()
	fmt.Printf("%s", stack)
}
