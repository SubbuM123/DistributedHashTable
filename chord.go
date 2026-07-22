package main

import (
	"errors"
	"net/rpc"
	"strconv"
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

	return n.FindSuccessorRPC(id)
}

func (n *Node) FindSuccessorRPC(id uint64) NodeInfo {
	client, err := rpc.Dial("tcp", n.successor.Addr)
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
