package main

import (
	"fmt"
	"net/rpc"
	"time"
)

type Query struct {
	Initiator int
	From      int
	To        int
}

func main() {
	// You can switch between the two graphs below.
	waitGraph1 := [][]int{
		{0, 0, 1, 0, 0},
		{1, 0, 0, 1, 0},
		{0, 1, 0, 0, 1},
		{0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
	} // Example with deadlock

	// waitGraph2 := [][]int{
	// 	{0, 1, 0, 0, 0},
	// 	{0, 0, 1, 0, 0},
	// 	{0, 0, 0, 1, 0},
	// 	{0, 0, 0, 0, 1},
	// 	{0, 0, 0, 0, 0},
	// } // Example without deadlock

	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	var isDeadlock bool
	err = client.Call("Server.Initialize", waitGraph1, &isDeadlock)
	if err != nil {
		fmt.Println("Error initializing:", err)
		return
	}

	// Initiate query from process 0.
	err = client.Call("Server.StartQuery", Query{Initiator: 0, From: 0, To: 1}, &isDeadlock)
	if err != nil {
		fmt.Println("Error starting query:", err)
		return
	}

	// Wait for asynchronous processing to complete.
	time.Sleep(3 * time.Second)

	// Ask the server for the final deadlock status.
	err = client.Call("Server.GetDeadlockStatus", struct{}{}, &isDeadlock)
	if err != nil {
		fmt.Println("Error getting deadlock status:", err)
		return
	}

	if isDeadlock {
		fmt.Println("Deadlock Detected !!")
	} else {
		fmt.Println("No Deadlock")
	}
}