package runner

import (
	"context"
	"log"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/moby/moby/client"
)

// Worker ...
type Worker struct {
	Cli *client.Client
}

// constaints ...
const (
	Workspace   = "/tmp/daizu/"
	MemoryLimit = 1024 * 1024 * 512 // nearly equal 512MB  TODO: Atcoderでは1024MB
)

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
func (w *Worker) Run(j Job) (err error) {
	containerID, err := w.CreateContainer(
		j.Image,
		j.MemoryLimit,
		[]mount.Mount{},
	)
	if err != nil {
		log.Fatalf("failed create container... %v\n", err)
	}

	atcOpt := dtypes.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	}

	hijacked, err := w.Cli.ContainerAttach(context.TODO(), containerID, atcOpt)
	if err != nil {
		log.Fatalf("failed hijack... %v\n", err)
	}
	log.Println(hijacked)
	return
}
