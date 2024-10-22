package common

import (
	"fmt"
	"os"
)

type File struct {
	FileName string
	Content  []byte
}

func NewFile(fileName string, content []byte) *File {
	return &File{
		FileName: fileName,
		Content:  content,
	}
}

func (f *File) SaveFile() error {
	file, err := os.Create(f.FileName)
	if err != nil {
		return fmt.Errorf("error creating file report: %v", err)
	}
	defer file.Close()

	_, err = file.Write(f.Content)
	if err != nil {
		return fmt.Errorf("error insert new content: %v", err)
	}
	return nil
}
