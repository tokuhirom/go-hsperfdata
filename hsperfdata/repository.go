package hsperfdata

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type HsperfdataRepository struct {
	dir string
}

func New() (*HsperfdataRepository, error) {
	var user string
	if runtime.GOOS == "windows" {
		user = os.Getenv("USERNAME")
	} else {
		user = os.Getenv("USER")
	}
	if user == "" {
		return nil, fmt.Errorf("error: Environment variable USER not set")
	}

	return NewUser(user)
}

func NewUser(userName string) (*HsperfdataRepository, error) {
	dir := filepath.Join(os.TempDir(), "hsperfdata_"+userName)
	return &HsperfdataRepository{dir}, nil
}

func (repository *HsperfdataRepository) GetFile(pid string) HsperfdataFile {
	return HsperfdataFile{filepath.Join(repository.dir, pid)}
}

func (repository *HsperfdataRepository) GetFiles() ([]HsperfdataFile, error) {
	files, err := ioutil.ReadDir(repository.dir)
	if err != nil {
		return nil, err
	}
	retval := make([]HsperfdataFile, len(files))
	for i, f := range files {
		retval[i] = HsperfdataFile{filepath.Join(repository.dir, f.Name())}
	}

	return retval, nil
}
