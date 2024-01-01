package main

import "fmt"

// Concept:- Channels block: It means that until the sender and reciever both are not ready the code will not run

func sender_not_ready() {
	my_channel := make(chan int)
	fmt.Println(<-my_channel)
	my_channel <- 45
}

func reciever_not_ready() {
	my_channel := make(chan int)
	my_channel <- 45
	fmt.Println(<-my_channel)
}

func main() {
	fmt.Println("Hello Channel block concept!!")
	// receiver_not_ready()
	// sender_not_ready()
}
