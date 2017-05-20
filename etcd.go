package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

type Backend struct {
	Name string
	Host string
	Port string
}

func GetBackends(client *etcd.Client,
	service string,
	backendName string) ([]Backend, error) {

	resp, err := client.Get(service, false, true)
	if err != nil {
		log.Println("Error when reading etcd: ", err)
		return nil, err
	}

	backends := make([]Backend, len(resp.Node.Nodes))

	for i, node := range resp.Node.Nodes {
		key := (*node).Key
		address := strings.Split(key[strings.LastIndex(key, "/")+1:], ":")

		backend := Backend{
			Name: fmt.Sprintf("back-%v", i),
			Host: address[0],
			Port: address[1],
		}

		backends[i] = backend
	}

	return backends, nil
}
