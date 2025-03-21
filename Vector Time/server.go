package main

import (
	"fmt"
	"net"
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

type ClockService struct{}

func (cs *ClockService) Send(sender *Process, reply *Process) error {
	sender.Clock[sender.Index] += 1
	*reply = *sender
	return nil
}

func (cs *ClockService) Receive(processes *SenderReceiver, reply *Process) error {
	*reply = *processes.Receiver
	// if processes.Sender.Clock > processes.Receiver.Clock {
	// 	reply.Clock = processes.Sender.Clock
	// }

	for i := 0; i < len(reply.Clock); i++ {
		if reply.Clock[i] < processes.Sender.Clock[i] {
			reply.Clock[i] = processes.Sender.Clock[i];
		}
	}
	reply.Clock[reply.Index] += 1
	return nil
}

func (cs *ClockService) Internal(sender *Process, reply *Process) error {
	sender.Clock[sender.Index] += 1
	*reply = *sender
	return nil
}

func main() {
	api := new(ClockService)
	rpc.Register(api)
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("[LOG]: Scalar Time RPC Server is Running on Port 1234...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
