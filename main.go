// mini-k8s: main.go
package main

import (
	"dude/master"
	"dude/node"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [master|node1|node2]")
		return
	}

	switch os.Args[1] {
	case "master":
		go master.StartAutoScaler() // run autoscaler
		master.StartMaster()
	case "node1":
		node.RunServer("9001")
	case "node2":
		node.RunServer("9002")
	default:
		fmt.Println("Unknown role")
	}
}

/*
todo: add capability to add pod dinamically through node worker endpoint by json payload or read json config
todo: replace simulateCPU using realCPU usage for each node --> required that each node is single vm/has cpu
*/
