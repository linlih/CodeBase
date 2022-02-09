package main

import "fmt"

type file struct {
	name string
}

func (f *file) print(identation string) {
	fmt.Println(identation + f.name)
}

func (f *file) clone() inode {
	return &file{name: f.name + "_clone"}
}
