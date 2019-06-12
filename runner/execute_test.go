package runner

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/moby/moby/client"
)

// getContainerList ...
func (c *client.Client) getContainerList() (list []types.Container, err error) {
	filter := map[string][]string{"name": {"<container_name>"}}
	filtBytes, err := json.Marshal(filter)
	if err != nil {
		log.Fatalf("%v", err)
	}
	filt, err := filters.FromParam(string(filtBytes))
	if err != nil {
		log.Fatalf("%v", err)
	}

	opts := types.ContainerListOptions{
		All:     true, // Include stopped containers
		Quiet:   true, // return only containerID
		Filters: filt,
	}

	list, err = c.ContainerList(context.TODO(), opts)
	if err != nil {
		log.Fatal("%v", err)
	}

	return
}

func (w *Worker) TestRun(t *testing.T) (err error) {
	return
}

// TODO: cをWorkerにしないといけない
func TestCreateContainer(t *testing.T) (err error) {
	c, err := NewClient()
	if err != nil {
		log.Fatalf("%v", err)
	}

	return
}
