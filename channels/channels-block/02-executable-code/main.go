package main

import "fmt"

// this would not execute because the code will block at the step of sending in the channel and our reciever
// goroutine has not launched yet
func sender_in_a_different_goroutine_not_executes() {
	my_chan := make(chan int)
	my_chan <- 12
	go func() {
		fmt.Println(<-my_chan)
	}()
}

// The reason it executes because the reciever goroutine is launched before the step of sending channel
func sender_in_a_different_goroutine_executes() {
	my_chan := make(chan int)

	go func() {
		fmt.Println(<-my_chan)
	}()

	my_chan <- 12
}

// The reason it executes because the sender goroutine is launched before the step of consuming channel
func reciever_in_a_different_goroutine_executes() {
	my_chan := make(chan int)

	go func() {
		my_chan <- 12
	}()

	fmt.Println(<-my_chan)
}

func main() {
	sender_in_a_different_goroutine_executes()
	reciever_in_a_different_goroutine_executes()
	//sender_in_a_different_goroutine_not_executes()
}
