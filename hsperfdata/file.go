package hsperfdata

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type HsperfdataFile struct {
	filename string
}

func (file *HsperfdataFile) GetPid() string {
	return filepath.Base(file.filename)
}

func (datafile *HsperfdataFile) Read() (*Result, error) {
	file, err := os.Open(datafile.filename)
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

	header := PerfdataHeader{}
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

	result := &Result{make(map[string]interface{})}

	start_offset := prologue.EntryOffset

	parsed_entries := int32(0)

	for parsed_entries < prologue.NumEntries {
		entry := DataEntryHeader2{}
		buffer := bytes.NewBuffer(data[start_offset:])
		err = binary.Read(buffer, header.GetEndian(), &entry)
		if err != nil {
			return nil, fmt.Errorf("Cannot read binary: %v", err)
		}

		name_start := int(start_offset) + int(entry.NameOffset)
		name_end := bytes.Index(data[name_start:], []byte{'\x00'})
		if name_end == -1 {
			return nil, fmt.Errorf("invalid binary: %v", err)
		}
		name := string(data[name_start : int(name_start)+name_end])
		data_start := start_offset + entry.DataOffset

		if entry.VectorLength == 0 {
			if entry.DataType != TYPE_LONG {
				return nil, fmt.Errorf("Unexpected monitor type: %v", entry.DataType)
			}
			i := int64(0)
			buffer := bytes.NewBuffer(data[data_start : data_start+8])
			err = binary.Read(buffer, header.GetEndian(), &i)
			if err != nil {
				return nil, fmt.Errorf("Cannot read binary: %v", err)
			}

			result.data[name] = fmt.Sprintf("%v", i)
		} else {
			if entry.DataType != TYPE_BYTE || entry.DataUnits != UNITS_STRING || (entry.DataVar != VARIABILITY_CONSTANT && entry.DataVar != VARIABILITY_VARIABLE) {
				return nil, fmt.Errorf("Unexpected vector monitor: DataType:%c,DataUnits:%v,DataVar:%v", entry.DataType, entry.DataUnits, entry.DataVar)
			}
			value := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(string(data[data_start:data_start+entry.VectorLength]), "\r"), "\n"), "\x00")

			result.data[name] = value
		}

		start_offset += entry.EntryLength
		parsed_entries++
	}

	return result, nil
}

type Result struct {
	data map[string]interface{}
}

func (self *Result) GetMap() map[string]interface{} {
	return self.data
}

func (self *Result) GetProcName() string {
	javaCommand := self.data["sun.rt.javaCommand"]
	if javaCommand != nil {
		if str, ok := javaCommand.(string); ok {
			splitted := strings.SplitN(str, " ", 2)
			return splitted[0]
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func (self *Result) GetString(key string) string {
	v := self.data[key]
	if v != nil {
		if str, ok := v.(string); ok {
			return str
		} else {
			return ""
		}
	} else {
		return ""
	}
}
