package main

import (
	"fmt"

	"github.com/sbinet/go-python"
)

func init() {
	err := python.Initialize()
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	m := python.PyImport_ImportModule("sys")
	if m == nil {
		fmt.Println("import error")
		return
	}
	path := m.GetAttrString("path")
	if path == nil {
		fmt.Println("get path error")
		return
	}
	
	currentDir := python.PyString_FromString("")
	python.PyList_Insert(path, 0, currentDir)

	m = python.PyImport_ImportModule("python_word_count")
	if m == nil {
		fmt.Println("import error")
		return
	}
	word_Count := m.GetAttrString("word_count")
	if word_Count == nil {
		fmt.Println("get word_Count error")
		return
	}
	word_Count.CallFunction()
	
}