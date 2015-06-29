package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/tokuhirom/go-hsperfdata/attach"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: hsstack pid\n")
		os.Exit(1)
	}
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("invalid pid: %v", err)
	}
	sock, err := attach.New(pid)
	if err != nil {
		log.Fatalf("cannot open unix socket: %s", err)
	}
	err = sock.Execute("threaddump")
	if err != nil {
		log.Fatalf("cannot write to unix socket: %s", err)
	}

	stack, err := sock.ReadString()
	fmt.Printf("%s\n", stack)
}
