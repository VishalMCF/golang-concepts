package main

import (
	"fmt"
	"github.com/tidwall/wal"
)

func main() {
	fmt.Println("Hello Wal Library exploration day....")

	// open a new log file
	log, _ := wal.Open("mylog", nil)

	// preparing batch for writes
	batch := new(wal.Batch)
	batch.Write(1, []byte("first entry of wal written in batch mode"))
	batch.Write(2, []byte("second entry of wal written in batch mode"))
	batch.Write(3, []byte("third entry of wal written in batch mode"))
	batch.Write(4, []byte("fourth entry of wal written in batch mode"))

	// writing to the log files
	err := log.WriteBatch(batch)
	if err != nil {
		return
	}

	// read an entry
	data, _ := log.Read(1)
	fmt.Println("data read for index {", 1, "} is :-> ", string(data))

	// close the log
	_ = log.Close()
}
