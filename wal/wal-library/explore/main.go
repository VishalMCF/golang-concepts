package main

import (
	"fmt"
	"github.com/tidwall/wal"
)

func main() {
	fmt.Println("Hello Wal Library exploration day....")

	// open a new log file
	log, _ := wal.Open("mylog", nil)

	// write some entries
	_ = log.Write(1, []byte("My first WAL entry"))
	_ = log.Write(2, []byte("My second WAL entry"))
	_ = log.Write(3, []byte("My third WAL entry"))

	// read an entry
	data, _ := log.Read(1)
	fmt.Println("data read for index {", 1, "} is :-> ", string(data))

	// close the log
	_ = log.Close()
}
