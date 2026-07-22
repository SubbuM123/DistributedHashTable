package main

import (
	"errors"
	"fmt"
	"net/rpc"
)

func (n *Node) SendPut(addr string, key string, value string) error {
	// n_id := curr.id
	// addr := curr.addr
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer client.Close()

	args := RPCArgs{
		Key:   key,
		Value: value,
	}

	var reply bool

	err = client.Call("Node.RemotePut", args, &reply)
	if err != nil {
		return err
	}

	if !reply {
		return errors.New("put failed")
	}

	return nil
}

func (n *Node) SendGet(addr string, key string) (string, error) {
	// addr := curr.addr
	fmt.Println("sent an rpc")
	// fmt.Println(addr)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		//fmt.Println("failed to open conn")
		return "null", err
	}
	defer client.Close()

	args := RPCArgs{
		Key:   key,
		Value: "",
	}

	var reply RPCReply

	err = client.Call("Node.RemoteGet", args, &reply)
	if err != nil {
		//fmt.Println("failed to rpc")
		return "null", err
	}

	if reply.Exists {
		return reply.Value, nil
	} else {
		//fmt.Println("failed to find key")
		return "null", errors.New("key not found")
	}
}

func (n *Node) SendDelete(addr string, key string) error {
    fmt.Println("addrs is " + addr)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer client.Close()

	args := RPCArgs{
		Key:   key,
		Value: "",
	}

	var reply bool

	err = client.Call("Node.RemoteDelete", args, &reply)
	if err != nil {
		return err
	}

	if reply {
		return nil
	}

	return errors.New("key not found")
}

// func (n *Node) SendDelete(n_id uint64, hash uint64) error {
// 	addr := n.neighbors[n_id]
// 	client, err := rpc.Dial("tcp", addr)
//     if err != nil {
//         return err
//     }
//     defer client.Close()

//     args := RPCArgs{
//         Hash:  hash,
// 		Value: "",
//     }

//     var reply bool

//     err = client.Call("Node.RemoteDelete", args, &reply)
//     if err != nil {
//         return err
//     }

//     if !reply {
//         return errors.New("get failed")
//     }

//     return nil
// }
