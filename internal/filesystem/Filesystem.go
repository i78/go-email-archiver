package filesystem

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LocalFilesystem struct {
}

func (l LocalFilesystem) Ls(directory string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(directory)
}

func (l LocalFilesystem) Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func ensureDir(fileName string) error {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			return merr
		}
	}
	return nil
}

func (l LocalFilesystem) MkDir(directory string) error {
	return ensureDir(directory)
}

func (l LocalFilesystem) Save(filename string, content []byte) error {
	ensureDir(filename)
	return ioutil.WriteFile(filename, content, 0600)
}

func Exists(path string) bool {
	if _, statError := os.Stat(path); !os.IsNotExist(statError) {
		return true
	}
	return false
}
