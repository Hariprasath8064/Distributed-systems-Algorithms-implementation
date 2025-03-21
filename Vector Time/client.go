package main

import (
	"fmt"
	"net/rpc"
)

type Process struct {
	Name  string
	Clock []int
	Index int
}

type SenderReceiver struct {
	Sender, Receiver *Process
}

func main() {
	client, _ := rpc.Dial("tcp", "localhost:1234")
	defer client.Close()

	P1 := &Process{Name: "P1", Clock: []int{0, 0, 0}, Index: 0}
	P2 := &Process{Name: "P2", Clock: []int{0, 0, 0}, Index: 1}
	P3 := &Process{Name: "P3", Clock: []int{0, 0, 0}, Index: 2}

	var reply Process

	client.Call("ClockService.Internal", P1, &reply)
	for i := 0; i < len(P1.Clock); i++ {
		P1.Clock[i] = reply.Clock[i]
	}
	fmt.Println("P1 Internal Event", P1.Clock)

	client.Call("ClockService.Internal", P2, &reply)
	for i := 0; i < len(P2.Clock); i++ {
		P2.Clock[i] = reply.Clock[i]
	}
	fmt.Println("P2 Internal Event", P2.Clock)

	client.Call("ClockService.Internal", P3, &reply)
	for i := 0; i < len(P3.Clock); i++ {
		P3.Clock[i] = reply.Clock[i]
	}
	fmt.Println("P3 Internal Event", P3.Clock)

	err := client.Call("ClockService.Send", P1, &reply)
	if err != nil {
		fmt.Println("RPC error:", err)
		return
	}
	P1.Clock = reply.Clock

	client.Call("ClockService.Receive", &SenderReceiver{P1, P2}, &reply)
	for i := 0; i < len(P2.Clock); i++ {
		P2.Clock[i] = reply.Clock[i]
	}
	fmt.Println("P1 -> P2", P1.Clock, P2.Clock)
}
