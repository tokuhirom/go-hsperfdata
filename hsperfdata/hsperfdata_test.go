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

	result, err := proc.Read()
	if err != nil {
		t.Error(err)
	}
	procName := result.GetProcName()
	if procName != "org.jetbrains.jps.cmdline.Launcher" {
		t.Errorf("GetProcname failed: %v", procName)
	}
	path := result.GetString("sun.property.sun.boot.library.path")
	expected := "/Library/Java/JavaVirtualMachines/jdk1.8.0_31.jdk/Contents/Home/jre/lib"
	if path != expected {
		t.Errorf("sun.property.sun.boot.library.path miss match: '%#v'(len:%d) != '%#v'(len:%d)", path, len(path), expected, len(expected))
	}

	files, err := repo.GetFiles()
	if err != nil {
		t.Error(err)
	}
	if len(files) != 2 {
		t.Errorf("GetFiles return 2 files, but it returns %d files", len(files))
	}
}
