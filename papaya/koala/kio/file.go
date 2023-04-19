package kio

import (
	"os"
)

type KFile struct {
	Path string
}

type KFileImpl interface {
	Cat() string
	IsExist() bool
	IsDir() bool
	IsFile() bool
	IsSymlink() bool
	IsSocket() bool
	IsBlockDev() bool
}

func KFileNew(path string) KFileImpl {

	file := &KFile{
		Path: path,
	}

	return file
}

func (f *KFile) Cat() string {

	data, _ := os.ReadFile(f.Path)
	return string(data)
}

func (f *KFile) IsExist() bool {

	_, err := os.Stat(f.Path)

	if err != nil {

		if !os.IsNotExist(err) {

			panic(err)
		}

		return false
	}

	return true
}

func (f *KFile) IsDir() bool {

	// use Stat, relative path like symlink -> directory, may problem too
	stat, err := os.Stat(f.Path)

	if err != nil {

		if !os.IsNotExist(err) {

			panic(err)
		}

		return false
	}

	return stat.IsDir()
}

func (f *KFile) IsFile() bool {

	stat, err := os.Lstat(f.Path)

	if err != nil {

		if !os.IsNotExist(err) {

			panic(err)
		}

		return false
	}

	return stat.Mode().IsRegular()
}

func (f *KFile) IsSymlink() bool {

	stat, err := os.Lstat(f.Path)

	if err != nil {

		if !os.IsNotExist(err) {

			panic(err)
		}

		return false
	}

	return stat.Mode()&os.ModeSymlink != 0
}

func (f *KFile) IsSocket() bool {

	stat, err := os.Lstat(f.Path)

	if err != nil {

		if !os.IsNotExist(err) {

			panic(err)
		}

		return false
	}

	return stat.Mode()&os.ModeSocket != 0
}

func (f *KFile) IsBlockDev() bool {

	stat, err := os.Lstat(f.Path)

	if err != nil {

		if !os.IsNotExist(err) {

			panic(err)
		}

		return false
	}

	return stat.Mode()&os.ModeDevice != 0
}
