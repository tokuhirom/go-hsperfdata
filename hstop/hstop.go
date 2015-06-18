package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
	termbox "github.com/nsf/termbox-go"
	"github.com/tokuhirom/go-hsperfdata/hsperfdata"
)

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

func print_headers() int {
	y := 0

	loadavg, err := linuxproc.ReadLoadAvg("/proc/loadavg")
	if err == nil {
		print_cell(0, y, fmt.Sprintf("Load Average: %5v %5v %5v", loadavg.Last1Min, loadavg.Last5Min, loadavg.Last15Min))
	} else {
		print_cell(0, y, fmt.Sprintf("Load Average: N/A"))
	}

	y++
	y++

	// T# : The total number of running threads
	buf := fmt.Sprintf("%5s %15s %3s %8s %8s %3s", "PID", "MAIN-CLASS", "NI", "VSIZE", "RSS", "T#")
	print_cell(0, y, buf)

	y++

	return y
}

func size(n uint64) string {
	if n > 1000*1000 {
		return fmt.Sprintf("%.1fM", (float64(n) / (1000 * 1000)))
	} else if n > 1000 {
		return fmt.Sprintf("%.1fk", (float64(n) / 1000))
	} else {
		return fmt.Sprintf("%v", n)
	}
}

func draw_all(repo *hsperfdata.Repository) {
	log.Println("draw_all")
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// print header
	header_lines := print_headers()

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

		pid := file.GetPid()

		buf := fmt.Sprintf("%5s %15s", file.GetPid(), shorten_main_class_name(info.GetProcName(), 15))
		stat, err := linuxproc.ReadProcessStat("/proc/" + pid + "/stat")
		if err == nil {
			buf += fmt.Sprintf(" %3d %8s %8s", stat.Nice, size(stat.Vsize), size(uint64(stat.Rss)))
		} else {
			buf += fmt.Sprintf(" %3s %8s %8s", "N/A", "N/A", "N/A")
		}
		buf += fmt.Sprintf(" %3s", info.GetString("java.threads.live"))

		print_cell(0, y+header_lines, buf)
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

	draw_all(repo)
	ticker := time.NewTicker(1 * time.Second)
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
