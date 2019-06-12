package runner

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/moby/moby/client"
)

// Worker ...
type Worker struct {
	Cli *client.Client
}

// NewWorker ...
func NewWorker() (worker Worker, err error) {
	cli, err := client.NewEnvClient()
	worker = Worker{
		Cli: cli,
	}
	return
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

