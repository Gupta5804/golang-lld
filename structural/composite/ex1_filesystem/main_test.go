package main

import (
	"testing"
)

func TestFolder_Search_Recursion(t *testing.T) {
	// Setting up the tree
	file1 := &File{name: "root_file.txt"}
	file2 := &File{name: "nested_file.txt"}
	file3 := &File{name: "nested_image.png"}

	folder1 := &Folder{name: "folder1"}
	folder1.Add(file2)
	folder1.Add(file3)

	rootFolder := &Folder{name: "rootFolder"}
	rootFolder.Add(folder1)
	rootFolder.Add(file1)

	results := rootFolder.Search("nested")

	if len(results) != 2 {
		t.Errorf("Expected 2 results , got %d", len(results))
	}

	foundFile2 := false
	for _, res := range results {
		if res.GetName() == "nested_file.txt" {
			foundFile2 = true
		}
	}

	if !foundFile2 {
		t.Error("Expected to find 'nested_file.txt', but it was missing from results")
	}
}

func TestFolder_GetSize_Recursion(t *testing.T) {
	file1 := &File{name: "file1", size: 10}
	file2 := &File{name: "file2", size: 20}

	folder1 := &Folder{name: "folder1"}
	root := &Folder{name: "root"}

	folder1.Add(file2)
	root.Add(folder1)
	root.Add(file1)

	if size := root.GetSize(); size != 30 {
		t.Errorf("Expected the size 30 , got %d", size)
	}
}
