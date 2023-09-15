package Zookeeper

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func Test() {
	// ZooKeeper server addresses
	zkServers := []string{"localhost:2181"}

	// Connect to ZooKeeper
	conn, _, err := zk.Connect(zkServers, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create a new ZkRegistry instance
	registry := NewZkRegistry("services", conn)

	// Register a node
	path := "127.0.0.1:50052"
	data := []byte("my-service-data")

	err = registry.Register(path, data)
	if err != nil {
		panic(err)
	}

	fmt.Println("Node registered successfully!")

	// Wait for a while to see the registered node in ZooKeeper
	time.Sleep(time.Second * 5)

	// Unregister the node
	err = registry.Unregister(path)
	if err != nil {
		panic(err)
	}

	fmt.Println("Node unregistered successfully!")
}
