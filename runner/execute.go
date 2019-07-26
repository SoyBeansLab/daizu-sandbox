package runner

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"

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
		AutoRemove:  true,
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
	defer hijacked.Close()

	return
}

// Exec ...
func Exec() (err error) {
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWUTS,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	}
	return
}

// InitContainer ...
// reference, https://github.com/rrreeeyyy/container-internship/tree/master/02
func InitContainer() error {
	if err := syscall.Sethostname([]byte("container")); err != nil {
		return fmt.Errorf("Setting hostname failed: %v", err)
	}

	if err := os.MkdirAll("/sys/fs/cgroup/cpu/daizu", 0700); err != nil {
		return fmt.Errorf("Cgroups namespace daizu create failed: %v", err)
	}
	if err := ioutil.WriteFile("/sys/fs/cgroup/cpu/daizu/tasks", []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644); err != nil {
		return fmt.Errorf("Cgroups register tasks to daizu namespace failed: %v", err)
	}
	if err := ioutil.WriteFile("/sys/fs/cgroup/cpu/daizu/cpu.cfs_quota_us", []byte("1000\n"), 0644); err != nil {
		return fmt.Errorf("Cgroups add limit cpu.cfs_quota_us to 1000 failed: %v", err)
	}

	if err := syscall.Mount("proc", "/root/rootfs/proc", "proc", syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, ""); err != nil {
		return fmt.Errorf("Proc mount failed: %v", err)
	}
	if err := os.Chdir("/root"); err != nil {
		return fmt.Errorf("Chdir /root failed: %v", err)
	}
	if err := syscall.Mount("rootfs", "/root/rootfs", "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Rootfs bind mount failed: %v", err)
	}
	if err := os.MkdirAll("/root/rootfs/oldrootfs", 0700); err != nil {
		return fmt.Errorf("Oldrootfs create failed: %v", err)
	}
	if err := syscall.PivotRoot("rootfs", "/root/rootfs/oldrootfs"); err != nil {
		return fmt.Errorf("PivotRoot failed: %v", err)
	}
	if err := syscall.Unmount("/oldrootfs", syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("Oldrootfs umount failed: %v", err)
	}
	if err := os.RemoveAll("/oldrootfs"); err != nil {
		return fmt.Errorf("Remove oldrootfs failed: %v", err)
	}
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("Chdir failed: %v", err)
	}
	if err := syscall.Exec("/bin/sh", []string{"/bin/sh"}, os.Environ()); err != nil {
		return fmt.Errorf("Exec failed: %v", err)
	}
	return nil
}
