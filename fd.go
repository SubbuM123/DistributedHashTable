package main

import (
	"net/rpc"
	"strconv"
	"time"
)

func (n *Node) PingNeighbor(neighbor NodeInfo, succ bool) {
	if neighbor.Addr == "" {
		return
	}

	client, err := rpc.Dial("tcp", neighbor.Addr)
	if err != nil {
		log_updates(n.id, "NODE FAILURE : Node "+strconv.Itoa(int(neighbor.Id))+" didn't connect")
		if succ {
			if len(n.successor_list) > 1 && n.successor_list[1].Addr != "" {
				oldSuccessor := n.successor
				n.successor = n.successor_list[1]
				n.updateSuccessorList(n.successor)
				if n.successor.Addr != "" && n.successor.Addr != oldSuccessor.Addr {
					if client, err := rpc.Dial("tcp", n.successor.Addr); err == nil {
						defer client.Close()
						var reply bool
						_ = client.Call("Node.RemoteNotifySuccessor", NodeInfo{Id: n.id, Addr: n.addr}, &reply)
					}
				}
			}
		} else {
			log_updates(n.id, "Predecessor failure detected; waiting for ring stabilization")
		}
		return
	}
	defer client.Close()

	var reply bool
	err = client.Call("Node.RemotePing", struct{}{}, &reply)
	if err != nil {
		log_updates(n.id, "NODE FAILURE : Node "+strconv.Itoa(int(neighbor.Id))+" is unreachable")
		if succ {
			if len(n.successor_list) > 1 && n.successor_list[1].Addr != "" {
				oldSuccessor := n.successor
				n.successor = n.successor_list[1]
				n.updateSuccessorList(n.successor)
				if n.successor.Addr != "" && n.successor.Addr != oldSuccessor.Addr {
					if client, err := rpc.Dial("tcp", n.successor.Addr); err == nil {
						defer client.Close()
						var reply bool
						_ = client.Call("Node.RemoteNotifySuccessor", NodeInfo{Id: n.id, Addr: n.addr}, &reply)
					}
				}
			}
		} else {
			log_updates(n.id, "Predecessor unreachable; waiting for ring stabilization")
		}
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
		n.PingNeighbor(n.successor, true)
		n.PingNeighbor(n.predecessor, false)
	}
}

// upon fail of successor
// move successor list
// call update successor list

// send a copy of data to new sucessor
