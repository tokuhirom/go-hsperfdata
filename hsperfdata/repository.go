package hsperfdata

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type Repository struct {
	dir string
}

func New() (*Repository, error) {
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

func NewUser(userName string) (*Repository, error) {
	dir := filepath.Join(os.TempDir(), "hsperfdata_"+userName)
	return &Repository{dir}, nil
}

func (repository *Repository) GetFile(pid string) File {
	return File{filepath.Join(repository.dir, pid)}
}

func (repository *Repository) GetFiles() ([]File, error) {
	files, err := ioutil.ReadDir(repository.dir)
	if err != nil {
		return nil, err
	}
	retval := make([]File, len(files))
	for i, f := range files {
		retval[i] = File{filepath.Join(repository.dir, f.Name())}
	}

	return retval, nil
}
