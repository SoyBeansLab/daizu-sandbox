package runner

import (
	"context"
	"log"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
)

// getContainerList ...
func (w *Worker) getContainerList() (list []types.Container, err error) {
	opts := types.ContainerListOptions{
		All: true, // Include stopped containers
	}

	list, err = w.Cli.ContainerList(context.TODO(), opts)
	if err != nil {
		log.Println(err)
	}

	return
}

func TestRun(t *testing.T) {
	job := Job{
		Language: "python",
		Image: "python",
		Source: "./",
		TimeLimit: 2,
		MemoryLimit: 1024*1024*1024*1024,
	}

	w, err := NewWorker()
	if err != nil {
		log.Println(err)
	}

	err = w.Run(job)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestCreateContainer(t *testing.T) {
	w, err := NewWorker()
	if err != nil {
		log.Println(err)
	}

	id, err := w.CreateContainer("python", 1024*1024*5, []mount.Mount{}, []string{})
	if err != nil {
		log.Println(err)
	}

	list, err := w.getContainerList()
	if err != nil {
		log.Println(err)
	}

	flag := false
	for _, i := range list {
		if i.ID == id {
			flag = true
			break
		}
	}

	if !flag {
		t.Errorf("cannot find %v in container list", id)
	} else {
		// remove container
		ctx := context.Background()
		err = w.Cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
		if err != nil {
			log.Println(err)
		}
	}
}
