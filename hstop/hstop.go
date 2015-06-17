package main

import (
	"fmt"
	"log"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/tokuhirom/go-hsperfdata/hsperfdata"
)

func print_cell(x_offset int, y int, s string) {
	for x, c := range s {
		termbox.SetCell(x+x_offset, y, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func draw_all(repo *hsperfdata.HsperfdataRepository) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	files, err := repo.GetFiles()
	if err != nil {
		log.Fatal(err)
	}
	for y, file := range files {
		procName, err := file.GetProcName()
		if err != nil {
			// process may existed after getting file list.
			continue
		}

		buf := fmt.Sprintf("%5s %15s", file.GetPid(), procName)

		print_cell(0, y, buf)
	}
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	repo, err := hsperfdata.New()
	if err != nil {
		log.Fatal(err)
	}

	draw_all(repo)
	ticker := time.NewTicker(time.Duration(1))
	go func() {
		for range ticker.C {
			draw_all(repo)
		}
	}()

L:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
			case termbox.KeyCtrlC:
				break L
			}
			draw_all(repo)
		case termbox.EventResize:
			draw_all(repo)
		}
	}
}
