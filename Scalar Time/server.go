package main

import (
	"fmt"
	"net"
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

type ClockService struct{}

func (cs *ClockService) Send(sender *Process, reply *Process) error {
	sender.Clock += 1
	*reply = *sender
	return nil
}

func (cs *ClockService) Receive(processes *SenderReceiver, reply *Process) error {
	*reply = *processes.Receiver
	if processes.Sender.Clock > processes.Receiver.Clock {
		reply.Clock = processes.Sender.Clock
	}
	reply.Clock += 1
	return nil
}

func (cs *ClockService) Internal(sender *Process, reply *Process) error {
	sender.Clock += 1
	*reply = *sender
	return nil
}

func (cs *ClockService) MatrixAdd(matrices *Matrix, reply *ReplyMatrix) error {
	n := len(matrices.M1);
	m := len(matrices.M1[0]);

	reply.Mat = make([][]int, n);
	for i := 0; i < n; i++ {
		reply.Mat[i] = make([]int, m);
		for j := 0; j < m; j++ {
			reply.Mat[i][j] += (matrices.M1[i][j] + matrices.M2[i][j]);
		}
	}
	return nil;
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
