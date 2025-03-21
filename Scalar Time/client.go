package main

import (
	"fmt"
	"net/rpc"
)

type Process struct {
	Name  string
	Clock int
}

type SenderReceiver struct {
	Sender, Receiver *Process
}

type Matrix struct {
	M1, M2 [][]int
}

type ReplyMatrix struct {
	Mat [][]int
}

func main() {
	client, _ := rpc.Dial("tcp", "localhost:1234")
	defer client.Close()

	mat1 := [][]int{
		{0, 0, 1},
		{1, 0, 0},
		{0, 1, 0},
	};

	mat2 := [][]int{
		{0, 0, 1},
		{1, 1, 1},
		{2, 1, 1},
	};

	P1 := &Process{Name: "P1", Clock: 0}
	P2 := &Process{Name: "P2", Clock: 0}
	P3 := &Process{Name: "P3", Clock: 0}

	var reply Process

	client.Call("ClockService.Internal", P1, &reply)
	P1.Clock = reply.Clock
	fmt.Println("P1 Internal Event", P1.Clock)

	client.Call("ClockService.Internal", P2, &reply)
	P2.Clock = reply.Clock
	fmt.Println("P2 Internal Event", P2.Clock)

	client.Call("ClockService.Internal", P3, &reply)
	P3.Clock = reply.Clock
	fmt.Println("P3 Internal Event", P3.Clock)

	var replyMat ReplyMatrix;
	client.Call("ClockService.MatrixAdd", &Matrix{M1: mat1, M2: mat2}, &replyMat);
	fmt.Println(replyMat.Mat);

	err := client.Call("ClockService.Send", P1, &reply)
	if err != nil {
		fmt.Println("RPC error:", err)
		return
	}
	P1.Clock = reply.Clock

	client.Call("ClockService.Receive", &SenderReceiver{P1, P2}, &reply)
	P2.Clock = reply.Clock
	fmt.Println("P1 -> P2", P1.Clock, P2.Clock)

	client.Call("ClockService.Send", P2, &reply)
	P2.Clock = reply.Clock

	client.Call("ClockService.Receive", &SenderReceiver{P2, P3}, &reply)
	P3.Clock = reply.Clock
	fmt.Println("P2 -> P3", P2.Clock, P3.Clock)
}
