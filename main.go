package main

import (
	// "fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

func main() {
	// start with the following format "./node bs id ip" -> "./node 127.0.0.1:8033 1 127.0.0.1:8001"
	// for starting node, bs can be START
	bs := os.Args[1]
	id := os.Args[2]
	log_addr := os.Args[3]
	node_id, _ := strconv.ParseUint(id, 10, 64)

	n := Node{
		hashtable: make(map[string]string),
		id:        uint64(node_id),
		addr:      log_addr,
		// neighbors: make(map[uint64]string),
	}

	log_updates(n.id, "Joined the system, yet to be placed")

	rpc.Register(&n)

	// Start listening
	listener, err := net.Listen("tcp", log_addr)
	if err != nil {
		log.Fatal(err)
	}

	// Accept RPC connections in the background
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	if err := n.JoinSystem(bs); err != nil {
		log_updates(n.id, "Join failed: "+err.Error())
	}

	n.event_manage()
	// switch {}
}
