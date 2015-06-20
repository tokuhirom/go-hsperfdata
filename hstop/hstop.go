package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
	termbox "github.com/nsf/termbox-go"
	"github.com/tokuhirom/go-hsperfdata/hsperfdata"
	"github.com/tokuhirom/go-hsperfdata/hstop/core"
	"github.com/tokuhirom/go-hsperfdata/hstop/fields"
	"github.com/tokuhirom/go-hsperfdata/hstop/state"
)

type MachineTopRenderer struct {
	fields []core.Field
}

func print_cell(x_offset int, y int, s string) {
	for x, c := range s {
		termbox.SetCell(x+x_offset, y, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func shorten_main_class_name(class_name string, length int) string {
	if len(class_name) > 15 {
		return class_name[len(class_name)-length:]
	} else {
		return class_name
	}
}

func (self *MachineTopRenderer) print_headers() int {
	y := 0

	loadavg, err := linuxproc.ReadLoadAvg("/proc/loadavg")
	if err == nil {
		print_cell(0, y, fmt.Sprintf("Load Average: %5v %5v %5v", loadavg.Last1Min, loadavg.Last5Min, loadavg.Last15Min))
	} else {
		print_cell(0, y, fmt.Sprintf("Load Average: N/A"))
	}

	y++
	y++

	{
		x := 0
		for _, field := range self.fields {
			buf := fmt.Sprintf("%"+strconv.Itoa(field.GetWidth())+"s ", field.GetTitle())
			print_cell(x, y, buf)
			x += len(buf)
		}
	}
	y++

	// T# : The total number of running threads

	return y
}

func (self *MachineTopRenderer) draw_all(repo *hsperfdata.Repository) {
	log.Println("draw_all")
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// print header
	header_lines := self.print_headers()

	files, err := repo.GetFiles()
	if err != nil {
		log.Fatal(err)
	}
	for y, file := range files {
		info, err := file.Read()
		if err != nil {
			// process may existed after getting file list.
			continue
		}

		state, err := state.New(file.GetPid(), info)
		if err != nil {
			log.Printf("%v", err)
		}

		{
			x := 0
			for _, field := range self.fields {
				buf := field.Render(state)
				print_cell(x, y+header_lines, buf)
				x += len(buf) + 1
			}
		}
	}
	termbox.Flush()
}

func main() {
	logfile := flag.String("log", "", "log file")
	flag.Parse()

	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	if *logfile != "" {
		f, err := os.OpenFile(*logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
		log.Println("started hstop")
	} else {
		log.SetOutput(ioutil.Discard)
	}

	repo, err := hsperfdata.New()
	if err != nil {
		log.Fatal(err)
	}

	renderer := &MachineTopRenderer{
		[]core.Field{
			&fields.PidField{},
			&fields.MainClassField{},
			&fields.NiceField{},
			&fields.VsizeField{},
			&fields.RssField{},
			&fields.ThreadsField{},
		},
	}

	renderer.draw_all(repo)
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			renderer.draw_all(repo)
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
			renderer.draw_all(repo)
		case termbox.EventResize:
			renderer.draw_all(repo)
		}
	}
}
