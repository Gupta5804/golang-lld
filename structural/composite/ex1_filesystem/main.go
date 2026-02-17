package main

import (
	"strings"
)

// component
type FileSystemNode interface {
	Search(keyword string) []FileSystemNode
	GetName() string
	GetSize() int
}

// Leaf
type File struct {
	name string
	size int
}

func (f *File) Search(keyword string) []FileSystemNode {
	result := []FileSystemNode{}
	if strings.Contains(f.name, keyword) {
		result = append(result, f)
	}
	return result
}
func (f *File) GetName() string {
	return f.name
}
func (f *File) GetSize() int {
	return f.size
}

// Composite
type Folder struct {
	name       string
	components []FileSystemNode
}

func (f *Folder) Search(keyword string) []FileSystemNode {
	results := []FileSystemNode{}
	for _, component := range f.components {
		results = append(results, component.Search(keyword)...)
	}
	return results
}
func (f *Folder) GetName() string {
	return f.name
}

func (f *Folder) Add(c FileSystemNode) {
	f.components = append(f.components, c)
}
func (f *Folder) GetSize() int{
	size := 0
	for _,component := range f.components{
		size += component.GetSize()
	}
	return size
}
