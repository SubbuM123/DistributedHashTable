# Distributed Hashtable in Go

This project implements a distributed hash table (DHT) in Go using a Chord-inspired ring architecture. It provides a simple key-value store that can be distributed across multiple nodes, with support for node joining, successor discovery, finger tables, replication, and failure detection.

The system is designed to distribute data across a dynamic set of participating nodes. Each node is responsible for a portion of the key space based on a consistent hash, and requests are routed through the ring to the appropriate owner. The implementation uses Go's built-in RPC package for communication between nodes and supports basic CRUD operations over the distributed store.

## Requirements

To run the project, install these:

- Go 1.26 or newer
- Git
- A terminal or command prompt

## Clone the Repository

```bash
git clone https://github.com/SubbuM123/DistributedHashTable.git
cd "Distributed Hashtable"
```

## Build the Project

From the repository root, build the binary:

```bash
go build -o dht.exe .
```

## Run the System

This project is started by launching one or more nodes. Each node must be given:

- a bootstrap address (this node will introduce the new one to others)
- a node ID
- a listening address

### Start the first node

```bash
./dht.exe START 1 127.0.0.1:8001
```

### Start a second node

Open a second terminal and run:

```bash
./dht.exe 127.0.0.1:8001 2 127.0.0.1:8002
```

### Start additional nodes

Additional nodes can be started in the same manner by providing the address of an existing node as the bootstrap target.

## Using the Interactive CLI

Once a node is running, you can interact with it through the terminal using the following commands:

- `PUT <key> <value>`: store a value under a key
- `GET <key>`: retrieve a stored value
- `DELETE <key>`: remove a key from the store
- `ls`: list locally stored keys
- `lsrep`: list replicated keys
- `s`: display the current successor list
- `f`: display finger table entries
- `n`: display current successor and predecessor information
- `EXIT`: shut down the node

Example:

```text
PUT user alice
GET user
DELETE user
```

### Maintenance and Status

Various logs are maintained throughout the system. logs/data_log tracks CRUD operations. logs/system_log tracks node joins, failures, and ring updates. Each node also updates a file with its content to serve as a disk backup in case of failure

## Selected Features

### Finger table support

Nodes maintain a finger table that helps accelerate routing decisions by providing shortcuts to nodes that are closer to a target key. This improves lookup efficiency across the distributed structure.

### RPC-based communication

All node-to-node communication is implemented using Go's remote procedure call mechanism. This provides a simple and reliable way for nodes to exchange metadata, route requests, and coordinate operations in a distributed environment.

### Replication

The implementation includes basic replica support so that data can be replicated to successor nodes. This improves fault tolerance and helps preserve data availability when a node becomes unreachable.

### Failure detection

Nodes periodically probe their successor and predecessor to detect failures. If a node is found to be unreachable, the system can transition to an alternate successor and continue operation with minimal disruption. Lost data is recovered and redistributed via replicas. Nodes who re-join the system will regain their data


## Project Structure

- `main.go`: entry point and node startup logic
- `node.go`: node state, join logic, and CLI handling
- `chord.go`: Chord-style successor and finger table logic
- `crud.go`: distributed CRUD operations
- `rpc.go`: RPC handlers for remote requests
- `send.go`: client-side RPC helpers
- `fd.go`: failure detection logic
- `replica.go`: replication behavior
- `logger.go`: logging utilities

## Notes

This is a learning-oriented distributed systems implementation and is intended to demonstrate core DHT concepts rather than provide a production-grade distributed database. It is a useful foundation for understanding peer-to-peer storage, consistent hashing, ring-based routing, and distributed fault tolerance.
