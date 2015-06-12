package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Header struct {
	// 0xcafec0c0
	Magic uint32
	// big_endian == 0, little_endian == 1
	ByteOrder byte
	Major     byte
	Minor     byte
	// Reserved  byte
}

// see sun/management/counter/perf/Prologue.java
type Prologue struct {
	Accessible   byte
	Used         int32
	Overflow     int32
	ModTimestamp int64
	EntryOffset  int32
	NumEntries   int32
}

type DataEntryHeader2 struct {
	EntryLength  int32
	NameOffset   int32
	VectorLength int32
	DataType     byte
	Flags        byte
	DataUnits    byte
	DataVar      byte
	DataOffset   int32
}

type DataEntry struct {
	key   string
	value string
}

const (
	TYPE_LONG            = 'J'
	TYPE_BYTE            = 'B'
	UNITS_STRING         = 5
	VARIABILITY_CONSTANT = 1
	VARIABILITY_VARIABLE = 3
)

func main() {
	if os.Args[1] == "ps" {
		Jps()
	} else if os.Args[1] == "stat" {
		dir, err := GetHsperfdataDirectoryForCurrentUser()
		if err != nil {
			log.Fatal("user", err)
		}
		pid := os.Args[2]
		ch, err := ReadHsperfdata(filepath.Join(dir, pid))
		if err != nil {
			log.Fatal("open fail", err)
		}
		for entry := range ch {
			fmt.Printf("%s=%v\n", entry.key, entry.value)
		}
	}
}

func Jps() {
	dir, err := GetHsperfdataDirectoryForCurrentUser()
	if err != nil {
		log.Fatal("user", err)
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("readdir", err)
	}
	for _, f := range files {
		proc_name, err := GetProcName(f.Name())
		if err != nil {
			log.Fatal("procname", err)
		}
		fmt.Printf("%s %s\n", f.Name(), proc_name)
	}
}

func GetProcName(pid string) (string, error) {
	dir, err := GetHsperfdataDirectoryForCurrentUser()
	if err != nil {
		log.Fatal("user", err)
	}
	ch, err := ReadHsperfdata(filepath.Join(dir, string(pid)))
	if err != nil {
		log.Fatal("open fail", err)
	}
	for entry := range ch {
		if entry.key == "sun.rt.javaCommand" {
			splitted := strings.SplitN(entry.value, " ", 2)
			return splitted[0], nil
		}
	}
	return "", errors.New("there's no sun.rt.javaCommand")
}

func GetHsperfdataDirectoryForCurrentUser() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return GetHsperfdataDirectory(user.Username), nil
}

func GetHsperfdataDirectory(user string) string {
	return filepath.Join(os.TempDir(), "hsperfdata_"+user)
}

func ReadHsperfdata(filename string) (chan DataEntry, error) {
	ch := make(chan DataEntry)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	var size int64 = info.Size()
	data := make([]byte, size)

	buffer1 := bufio.NewReader(file)
	_, err = buffer1.Read(data)
	if err != nil {
		return nil, err
	}

	header := Header{}
	{
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &header)
		if err != nil {
			return nil, err
		}
		if header.Magic != 0xcafec0c0 {
			return nil, errors.New("illegal magic")
		}
		if !(header.Major == 2 && header.Minor == 0) {
			return nil, errors.New("unsupported version")
		}
	}

	prologue := Prologue{}
	{
		buffer := bytes.NewBuffer(data[7:])
		err = binary.Read(buffer, header.GetEndian(), &prologue)
		if prologue.Accessible != 1 {
			return nil, errors.New("not accessible")
		}
	}

	go func() {
		start_offset := prologue.EntryOffset

		parsed_entries := int32(0)

		for parsed_entries < prologue.NumEntries {
			entry := DataEntryHeader2{}
			buffer := bytes.NewBuffer(data[start_offset:])
			err = binary.Read(buffer, header.GetEndian(), &entry)
			if err != nil {
				log.Fatal("binary.read", err)
			}

			name_start := int(start_offset) + int(entry.NameOffset)
			name_end := bytes.Index(data[name_start:], []byte{'\x00'})
			if name_end == -1 {
				log.Fatal("invalid format", nil)
			}
			name := string(data[name_start : int(name_start)+name_end])
			data_start := start_offset + entry.DataOffset

			if entry.VectorLength == 0 {
				if entry.DataType != TYPE_LONG {
					log.Fatal("Unexpected monitor type", nil)
				}
				i := int64(0)
				buffer := bytes.NewBuffer(data[data_start : data_start+8])
				err = binary.Read(buffer, header.GetEndian(), &i)
				if err != nil {
					log.Fatal("binary.read", err)
				}
				ch <- DataEntry{name, fmt.Sprintf("%v", i)}
			} else {
				if entry.DataType != TYPE_BYTE || entry.DataUnits != UNITS_STRING || (entry.DataVar != VARIABILITY_CONSTANT && entry.DataVar != VARIABILITY_VARIABLE) {
					log.Fatal(fmt.Sprintf("Unexpected vector monitor: DataType:%c,DataUnits:%v,DataVar:%v", entry.DataType, entry.DataUnits, entry.DataVar), nil)
				}
				value := strings.TrimSuffix(string(data[data_start:data_start+entry.VectorLength]), "\r\n\x00")

				ch <- DataEntry{name, value}
			}

			start_offset += entry.EntryLength
			parsed_entries++
		}

		close(ch)
	}()

	return ch, err
}

func (header *Header) GetEndian() binary.ByteOrder {
	if header.ByteOrder == 0 {
		return binary.BigEndian
	} else {
		return binary.LittleEndian
	}
}
