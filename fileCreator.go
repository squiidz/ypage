package ypage

import (
	"fmt"
	"io/ioutil"
	"os"
)

type FileData struct {
	Path     string
	Filename string
	Right    os.FileMode
	Data     []byte
}

func NewFile(p string, n string, r os.FileMode) *FileData {
	f := &FileData{
		Path:     p,
		Filename: n,
		Right:    r,
	}
	err := ioutil.WriteFile(f.Path+"/"+f.Filename, nil, f.Right)
	if err != nil {
		fmt.Println("FILE NOT CREATED")
		return nil
	}
	return f
}

func (f *FileData) Insert(data []byte) error {
	file, err := os.OpenFile(f.Path+"/"+f.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, f.Right)
	if err != nil {
		fmt.Println("CANNOT OPEN FILE")
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("NOT ABLE TO WRITE TO FILE:", err)
		return err
	}
	return nil
}
