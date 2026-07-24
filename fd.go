package main

import ("net/rpc"
		"strconv"
		"time")

func (n *Node) PingNeighbor(neighbor NodeInfo) {
	if neighbor.Addr == "" {
        return
    }
	client, err := rpc.Dial("tcp", neighbor.Addr)
	if err != nil {
		log_updates(n.id, "NODE FAILURE : Node " + strconv.Itoa(int(neighbor.Id)) + " didn't connect")
		// need to call stabilization here
		return
	}
	defer client.Close()

	var reply bool
	err = client.Call("Node.RemotePing", struct{}{}, &reply)
	if err != nil {
		log_updates(n.id, "NODE FAILURE : Node " + strconv.Itoa(int(neighbor.Id)) + " is unreachable")
		// need to call stabilization here
		return
	}
}

func (n *Node) RemotePing(_ struct{}, reply *bool) error {
	*reply = true
	return nil
}

func (n *Node) FailureDetector() {
	ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        n.PingNeighbor(n.successor)
		n.PingNeighbor(n.predecessor)
    }
}

// upon fail of successor
// move successor list
// call update successor list

// send a copy of data to new sucessor