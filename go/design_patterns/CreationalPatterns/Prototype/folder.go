package main

import "fmt"

type folder struct {
	children []inode
	name     string
}

func (f *folder) print(identation string) {
	fmt.Println(identation + f.name)
	for _, i := range f.children {
		i.print(identation + identation)
	}
}

func (f *folder) clone() inode {
	cloneFolder := &folder{name: f.name + "_clone"}
	var tempChildren []inode
	for _, i := range f.children {
		copy := i.clone()
		tempChildren = append(tempChildren, copy)
	}
	cloneFolder.children = tempChildren
	return cloneFolder
}
