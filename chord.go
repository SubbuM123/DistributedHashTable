package main

import (
	"errors"
	// "fmt"
	"net/rpc"
	"strconv"
	"time"
)

func (n *Node) FindOwner(hash uint64) string {
	if n.Owns(hash) {
		return n.addr
	} else {
		succ := n.FindSuccessor(hash)
		return succ.Addr
	}
}

func (n *Node) Owns(hash uint64) bool {
	if n.successor.Addr == "" || n.successor.Addr == n.addr {
		return true
	}

	if n.predecessor.Id == 0 || n.predecessor.Addr == "" {
		return true
	}

	if n.predecessor.Id == n.id && n.successor.Id == n.id {
		return true
	}

	if n.predecessor.Id < n.id {
		return hash > n.predecessor.Id && hash <= n.id
	}

	return hash > n.predecessor.Id || hash <= n.id
}

// func (n *Node) Join()
func (n *Node) FindSuccessor(id uint64) NodeInfo {
	if n.successor.Addr == "" || n.successor.Addr == n.addr {
		return NodeInfo{Id: n.id, Addr: n.addr}
	}

	if n.Owns(id) {
		log_updates(n.id, "Sent successor to "+strconv.Itoa(int(id))+" (was self)")
		return NodeInfo{Id: n.id, Addr: n.addr}
	}

	if n.id < n.successor.Id && id > n.id && id <= n.successor.Id {
		log_updates(n.id, "Sent successor to "+strconv.Itoa(int(id))+" (was successor)")
		return n.successor
	}

	if n.id > n.successor.Id && (id > n.id || id <= n.successor.Id) {
		log_updates(n.id, "Sent successor to "+strconv.Itoa(int(id))+" (wrapped)")
		return n.successor
	}
	f := n.FindFinger(id)
	return n.FindSuccessorRPC(f, id)
}

func (n *Node) FindSuccessorRPC(target NodeInfo, id uint64) NodeInfo {
	if target.Addr == "" || target.Addr == n.addr {
		return NodeInfo{Id: n.id, Addr: n.addr}
	}

	client, err := rpc.Dial("tcp", target.Addr)
	if err != nil {
		return NodeInfo{Id: n.id, Addr: n.addr}
	}
	defer client.Close()

	var reply NodeInfo
	err = client.Call("Node.RemoteFindSuccessor", id, &reply)
	if err != nil {
		return NodeInfo{Id: n.id, Addr: n.addr}
	}

	return reply
}

func (n *Node) RemoteFindSuccessor(id uint64, reply *NodeInfo) error {
	succ := n.FindSuccessor(id)
	if succ.Addr == "" {
		return errors.New("no successor found")
	}

	*reply = succ
	return nil
}

func (n *Node) StartFingerTable() {
	for i := 0; i < 8; i++ {
		start := (n.id + (1 << i)) % 256
		n.fingers[i].Num = start
		n.fingers[i].FingerNode = n.FindSuccessor(start)
	}
}

func (n *Node) MakeFingerTable() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		n.StartFingerTable()
	}
}

func (n *Node) FindFinger(hash uint64) NodeInfo {
	if len(n.fingers) == 0 {
		return NodeInfo{Id: n.id, Addr: n.addr}
	}

	for i := len(n.fingers) - 1; i >= 0; i-- {
		if n.fingers[i].FingerNode.Id == 0 || n.fingers[i].FingerNode.Addr == "" {
			continue
		}

		if betweenFinger(n.id, n.fingers[i].FingerNode.Id, hash) {
			return n.fingers[i].FingerNode
		}
	}

	return n.successor
}

func betweenFinger(start, target, end uint64) bool {
	if start == end {
		return true
	}
	if start < end {
		return target > start && target < end
	}
	return target > start || target < end
}
