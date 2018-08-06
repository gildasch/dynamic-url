package faceswap

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	docker "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/mount"
	"docker.io/go-docker/api/types/network"
	"github.com/pkg/errors"
)

type Wuhuikais struct {
	*docker.Client

	imageID string
}

func (w *Wuhuikais) FaceSwap() error {
	containerID := fmt.Sprintf("faceswap_run_%d", time.Now().Unix())

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = w.ContainerCreate(context.Background(),
		&container.Config{
			Image: "faceswap:latest",
			Cmd: []string{
				"--dst", "/images/Rogelio.png",
				"--src", "/images/gildas2.png",
				"--out", "/images/out.jpg", "--correct_color"},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: dir,
					Target: "/images",
				},
			},
			AutoRemove: true,
		},
		&network.NetworkingConfig{}, containerID)
	if err != nil {
		return err
	}

	err = w.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	statusCh, errCh := w.ContainerWait(context.Background(), containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	return nil
}

type logLine struct {
	Stream string `json:"stream"`
}

func (w *Wuhuikais) build() error {
	buildCtx, err := buildContext()
	if err != nil {
		return err
	}

	resp, err := w.ImageBuild(context.Background(), buildCtx, types.ImageBuildOptions{
		Tags: []string{"faceswap"},
	})
	if err != nil {
		return err
	}

	var line string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line = scanner.Text()
		fmt.Println(line)

		var log logLine
		err := json.Unmarshal([]byte(line), &log)
		if err != nil {
			continue
		}

		if strings.HasPrefix(log.Stream, "Successfully built ") {
			w.imageID = strings.TrimSuffix(
				strings.TrimPrefix(log.Stream, "Successfully built "), "\n")
		}
	}

	if line != `{"stream":"Successfully tagged faceswap:latest\n"}` {
		return errors.Errorf("error building faceswap image, last line was %q", line)
	}

	return nil
}

func buildContext() (*bytes.Buffer, error) {
	dockerfile, err := ioutil.ReadFile("Dockerfile")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	hdr := &tar.Header{
		Name: "Dockerfile",
		Mode: 0600,
		Size: int64(len(dockerfile)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}
	if _, err := tw.Write(dockerfile); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return &buf, nil
}
