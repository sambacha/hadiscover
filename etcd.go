package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

// Service is a named collection of backend Instances
//
// For example, and API service may have multiple Instances that it can
// send requests to.
type Service struct {
	Name      string
	Instances []Instance
}

// Instance represents an instance of a service.
//
// For example, a container exposing an API service could be an instance.
type Instance struct {
	Name string
	Host string
	Port string
}

// GetServices reads the current state of Services and Instances from the
// key in etcd.
func GetServices(client *etcd.Client, key string) ([]Service, error) {
	// read data from the key, recursively
	resp, err := client.Get(key, false, true)
	if err != nil {
		log.Println("Error when reading etcd: ", err)
		return nil, err
	}

	services := []Service{}

	// iterate over the nodes in the key we are watching
	for _, n := range resp.Node.Nodes {
		// get the parts of the name
		parts := strings.Split(n.Key, "/")

		// name of the service is the last part, assuming that the first part
		// is the prefix we are listening to
		name := parts[len(parts)-1]

		// add the backend instances to the service
		backends := make([]Instance, len(n.Nodes))

		for i, b := range n.Nodes {
			address := strings.Split(b.Value, ":")

			backends[i] = Instance{
				Name: fmt.Sprintf("%v-%v", name, i),
				Host: address[0],
				Port: address[1],
			}
		}

		service := Service{
			Name:      name,
			Instances: backends,
		}

		services = append(services, service)
	}

	return services, nil
}
