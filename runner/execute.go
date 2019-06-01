//package runner
package main

import (
	"context"
	"log"

    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/network"
	"github.com/moby/moby/client"
)

// Worker ...
type Worker struct {
	cli *client.Client
}

// Run ...
func (w *Worker) Run() {
	config := &container.Config{
		Image: "python",
	}
	hostConfig := &container.HostConfig {}
	netConfig := &network.NetworkingConfig{}
	resp, err := w.cli.ContainerCreate(context.TODO(), config, hostConfig, netConfig, "python1")
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%v", resp.ID)
}

func main() {
	c, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("%v", err)
	}
	w := &Worker{c}
	w.Run()
}
