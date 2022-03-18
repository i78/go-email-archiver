package repository

import "io/fs"

type Filesystem interface {
	MkDir(directory string) error
	Ls(directory string) ([]fs.FileInfo, error)
	Read(filename string) ([]byte, error)
	Save(filename string, content []byte) error
}
