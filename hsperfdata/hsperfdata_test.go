package hsperfdata

import (
	"testing"
)

func TestNew(t *testing.T) {
	repo, err := New()
	if err != nil {
		t.Error(err)
	}

	repo.dir = "../test-data/hsperfdata_tokuhirom/"
	proc := repo.GetFile("21916")
	if proc.GetPid() != "21916" {
		t.Error("GetPid failed")
	}
	procName, err := proc.GetProcName()
	if err != nil {
		t.Error(err)
	}
	if procName != "org.jetbrains.jps.cmdline.Launcher" {
		t.Errorf("GetProcname failed: %v", procName)
	}

	ch, err := proc.ReadHsperfdata()
	if err != nil {
		t.Error(err)
	}
	m := make(map[string]interface{})
	for entry := range ch {
		m[entry.Key] = entry.Value
	}
	if str, ok := m["sun.property.sun.boot.library.path"].(string); ok {
		expected := "/Library/Java/JavaVirtualMachines/jdk1.8.0_31.jdk/Contents/Home/jre/lib"
		if str != expected {
			t.Errorf("sun.property.sun.boot.library.path miss match: '%#v'(len:%d) != '%#v'(len:%d)", str, len(str), expected, len(expected))
		}
	} else {
		t.Error("'sun.property.sun.boot.library.path' is not a string")
	}

	files, err := repo.GetFiles()
	if err != nil {
		t.Error(err)
	}
	if len(files) != 2 {
		t.Errorf("GetFiles return 2 files, but it returns %d files", len(files))
	}
}
