package util

import (
	"errors"
	"io/ioutil"
	"os"
)

func CopyFileInefficientlyOrPanic(source, dest string) {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(dest, input, 0666)
	if err != nil {
		panic(err)
	}
}

func EnsureFile(p string) error {
	_, err := os.Stat(p)
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(p)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}
