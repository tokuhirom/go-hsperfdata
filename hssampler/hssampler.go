package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/tokuhirom/go-hsperfdata/attach"
)

type HsSampler struct {
	pid        int
	re         *regexp.Regexp
	logs       []map[string]int
	limit_logs int
}

func (self *HsSampler) sample() {
	// check attacher state
	attacher, err := attach.New(self.pid)
	if err != nil {
		log.Fatal(err)
	}
	defer attacher.Close()

	stack, err := attacher.RemoteDataDump()
	if err != nil {
		log.Fatal(err)
	}

	row := make(map[string]int)
	for _, method := range self.scan_methods(stack) {
		if method != "java.lang.Object.wait" && method != "java.lang.Thread.run" {
			row[method] += 1
		}
	}
	self.add_log(row)
}

type SampleResult struct {
	method string
	count  int
}

type SampleResults []SampleResult

func (p SampleResults) Len() int {
	return len(p)
}

func (p SampleResults) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p SampleResults) Less(i, j int) bool {
	return p[i].count > p[j].count
}

func print_cell(x_offset int, y int, s string) {
	for x, c := range s {
		termbox.SetCell(x+x_offset, y, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (self *HsSampler) render() {
	stat := make(map[string]int)
	for _, row := range self.logs {
		for k, v := range row {
			stat[k] += v
		}
	}

	results := make(SampleResults, 0)
	for method, cnt := range stat {
		results = append(results, SampleResult{method, cnt})
	}
	sort.Sort(results)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for i, result := range results {
		print_cell(0, i, fmt.Sprintf("%7d %s", result.count, result.method))
	}
	termbox.Flush()
}

func (self *HsSampler) add_log(row map[string]int) {
	self.logs = append(self.logs, row)
	if len(self.logs) > self.limit_logs {
		// unshift logs
		self.logs = self.logs[1:]
	}
}

func (self *HsSampler) scan_methods(stack string) []string {
	methods := make([]string, 0)
	self.re.ReplaceAllStringFunc(stack, func(s string) string {
		method := s[3 : len(s)-1]
		methods = append(methods, method)
		return s
	})
	return methods
}

func NewSampler(pid int, log_limit int) (*HsSampler, error) {
	re, err := regexp.Compile("at ([^(]+)\\(")
	if err != nil {
		return nil, err
	}
	return &HsSampler{pid, re, make([]map[string]int, 0), log_limit}, nil
}

/**
 * hssampler is a sampler for hotspot vm.
 *
 * Usage: hssapmler pid
 */
func main() {
	logfile := flag.String("log", "", "log file")
	log_limit := flag.Int("log_limit", 50, "number of sampler back logs")
	sample_interval := flag.Int("sample_interval", 500, "sample interval [millisecond]")
	refresh_interval := flag.Int("refresh_interval", 1000, "refresh interval [millisecond]")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println("Usage: hssampler pid")
		os.Exit(1)
	}

	pid, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatalf("bad pid: %v", err)
	}

	sampler, err := NewSampler(pid, *log_limit)
	if err != nil {
		log.Fatal(err)
	}

	if *logfile != "" {
		f, err := os.OpenFile(*logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
		log.Println("started hssampler")
	} else {
		log.SetOutput(ioutil.Discard)
	}

	err = termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	sampler.sample()
	sampler.render()
	{
		ticker := time.NewTicker(time.Duration(*sample_interval) * time.Millisecond)
		go func() {
			for range ticker.C {
				sampler.sample()
			}
		}()
	}
	{
		ticker := time.NewTicker(time.Duration(*refresh_interval) * time.Millisecond)
		go func() {
			for range ticker.C {
				sampler.render()
			}
		}()
	}

L:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
			case termbox.KeyCtrlC:
				break L
			}
			sampler.render()
		case termbox.EventResize:
			sampler.render()
		}
	}
}
