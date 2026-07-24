package main

import (
	"os"
	"fmt"
)

func log_updates(node_id uint64, update string) {
	file, err := os.OpenFile(
		"logs/system_log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	fmt.Fprintf(
		file,
		"Node %d > %s \n",
		node_id,
		update,
	)
}

func log_data(node_id uint64, update string) {
	file, err := os.OpenFile(
		"logs/data_log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	fmt.Fprintf(
		file,
		"Node %d > %s \n",
		node_id,
		update,
	)
}

func (n *Node) BackupReplica() {
	file, err := os.OpenFile(
		"logs/backup.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	fmt.Fprintf(
		file,
		"\n Node %d BACKUP \n",
		n.predecessor.Id,
	)

	for key, value := range n.pred_replica {
		fmt.Fprintf(
		file,
		"%s : %s \n",
		key,
		value,)
	}
	
}