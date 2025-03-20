package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type DependencyGraph [][]int

type Query struct {
	Initiator int
	From      int
	To        int
}

type Reply struct {
	Initiator int
	From      int
}

// Server holds the dependency graph, per-process data and the global deadlock flag.
type Server struct {
	dependencyGraph DependencyGraph
	num             []int
	wait            []bool
	detected        bool // global deadlock flag
}

// Initialize sets up the dependency graph and clears local states.
func (s *Server) Initialize(graph DependencyGraph, reply *bool) error {
	s.dependencyGraph = graph
	n := len(graph)
	s.num = make([]int, n)
	s.wait = make([]bool, n)
	s.detected = false
	*reply = false	
	return nil
}

// StartQuery follows the algorithm: if a process is not waiting, mark it and send queries
// to all processes it depends on; if already waiting, send a reply back.
func (s *Server) StartQuery(q Query, reply *bool) error {
	if !s.wait[q.To] {
		s.num[q.To] = s.countDependencies(q.To);
		s.wait[q.To] = true;
		for i, dep := range s.dependencyGraph[q.To] {
			if dep == 1 {
				go func(i int) {
					client, err := rpc.Dial("tcp", "localhost:1234");
					if err != nil {
						fmt.Println("Error dialing:", err);
						return;
					}
					defer client.Close();
					fmt.Println("[LOG]: Sending Query to", i, "for Deadlock Detection. Initiator:", q.Initiator);
					var dummy bool;
					// We ignore the reply from asynchronous calls because they update the server’s global flag.
					client.Call("Server.StartQuery", Query{Initiator: q.Initiator, From: q.To, To: i}, &dummy);
				}(i); 
			}
		}
	} else {
		client, err := rpc.Dial("tcp", "localhost:1234");
		if err != nil {
			fmt.Println("Error dialing:", err);
			return err;
		}
		defer client.Close();
		fmt.Println("[LOG]: Sending Reply From", q.To);
		var dummy bool;
		client.Call("Server.ReceiveReply", Reply{Initiator: q.Initiator, From: q.To}, &dummy);
	}
	return nil;
}


// ReceiveReply decrements the waiting counter. If a process finishes (num becomes 0),
// if it is the initiator, deadlock is detected; otherwise, the reply is propagated.
func (s *Server) ReceiveReply(r Reply, reply *bool) error {
	if s.wait[r.From] {
		s.num[r.From] -= 1;
		if s.num[r.From] == 0 {
			if r.Initiator == r.From {
				fmt.Println("Deadlock detected!")
				s.detected = true;
			} else {
				for i, dep := range s.dependencyGraph[r.From] {
					if dep == 1 {
						client, err := rpc.Dial("tcp", "localhost:1234");
						if err != nil {
							fmt.Println("Error dialing:", err);
							continue;
						}
						defer client.Close();
						var dummy bool;
						client.Call("Server.ReceiveReply", Reply{Initiator: r.Initiator, From: i}, &dummy);
					}
				}
			}
		}
	}
	*reply = s.detected;
	return nil;
}

// countDependencies simply counts the number of dependencies for a given process.
func (s *Server) countDependencies(p int) int {
	count := 0;
	for _, dep := range s.dependencyGraph[p] {
		if dep == 1 {
			count += 1;
		}
	}
	return count;
}

// GetDeadlockStatus returns the global deadlock flag.
func (s *Server) GetDeadlockStatus(args struct{}, reply *bool) error {
	*reply = s.detected;
	return nil;
}

func main() {
	server := new(Server);
	rpc.Register(server);

	listener, err := net.Listen("tcp", ":1234");
	if err != nil {
		fmt.Println("Error starting server:", err);
		return;
	}
	defer listener.Close();

	fmt.Println("[LOG]: Server for Deadlock Detection Running at Port 1234...");
	for {
		conn, err := listener.Accept();
		if err != nil {
			fmt.Println("Error accepting connection:", err);
			continue;
		}
		go rpc.ServeConn(conn);
	}
}