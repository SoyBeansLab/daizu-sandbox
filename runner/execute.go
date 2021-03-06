package runner

import (
	"context"
	"io/ioutil"
	"log"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/moby/moby/client"
	"github.com/moby/moby/pkg/stdcopy"
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
	if err != nil {
		log.Fatalf("failed create new client... %v\n", err)
		return
	}

	worker = Worker{
		Cli: cli,
	}
	return
}

// CreateContainer ...
func (w *Worker) CreateContainer(img string, memoryLimit int64, mounts []mount.Mount, cmd []string) (containerID string, err error) {
	// https://godoc.org/github.com/docker/docker/api/types/container#Config
	config := &container.Config{
		Image:        img,
		WorkingDir:   Workspace,
		Cmd:          cmd,
		Tty:          false,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}

	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			CpusetCpus: "0",
			PidsLimit:  0,
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
func (w *Worker) Run(j *Job) (err error) {
	containerID, err := w.CreateContainer(
		j.Image,
		j.MemoryLimit,
		[]mount.Mount{},
		j.Cmd,
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

	hijacked, err := w.Cli.ContainerAttach(context.Background(), containerID, atcOpt)
	if err != nil {
		log.Fatalf("failed hijack... %v\n", err)
	}
	defer hijacked.Close()

	err = w.Cli.ContainerStart(context.Background(), containerID, dtypes.ContainerStartOptions{})
	if err != nil {
		log.Fatalf("failed start... %v\n", err)
	}

	j.Stdout, err = ioutil.TempFile("", containerID+"out")
	if err != nil {
		log.Fatal(err)
	}
	j.Stderr, err = ioutil.TempFile("", containerID+"err")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stdcopy.StdCopy(j.Stdout, j.Stderr, hijacked.Reader)
	if err != nil {
		log.Fatal(err)
	}

	return
}
