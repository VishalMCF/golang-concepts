package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Hello Wal Core concepts by dissecting a library....")

	// how to print root directory
	path, err := filepath.Abs("mylog")

	// how to get current working directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("error -> ", err)
		return
	}

	fmt.Println("Current working directory -> ", dir)
	if err != nil {
		return
	}

	fmt.Println("curret root directiry -> ", path)
}
