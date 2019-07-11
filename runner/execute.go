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

// CreateContainer ...
func (w *Worker) CreateContainer(img string, memoryLimit int64, mounts []mount.Mount) (containerID string, err error) {
	config := &container.Config{
		Image: img,
	}

	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			CpusetCpus: "0",
			PidsLimit:  50,
			Memory:     memoryLimit,
		},
		NetworkMode: "none",
		Mounts:      mounts,
	}
	netConfig := &network.NetworkingConfig{}
	resp, err := w.Cli.ContainerCreate(context.TODO(), config, hostConfig, netConfig, "")

	if err != nil {
		log.Fatalf("%v", err)
	}

	containerID = resp.ID
	log.Printf("Create %v", resp.ID)

	return
}

// Run ...
func (w *Worker) Run() (err error) {
	return
}
