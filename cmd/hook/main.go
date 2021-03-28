package main

import (
	"encoding/json"
	"fmt"
	"github.com/run-ai/runai-container-runtime/pkg/bundle"
	"github.com/run-ai/runai-container-runtime/pkg/logger"
	"github.com/run-ai/runai-container-runtime/pkg/mount"
	"log"
	"os"
)

type containerMetadata struct {
	pid int `json:"pid,omitempty"`
	bundle string `json:"bundle"`
	bundlePath string `json:"bundlePath"`
}

func getContainerMetadata() (*containerMetadata, error) {
	var c containerMetadata
	d := json.NewDecoder(os.Stdin)
	if err := d.Decode(&c); err != nil {
		return nil, err
	}

	if len(c.bundle) == 0 {
		c.bundle = c.bundlePath
	}
	return &c, nil
}

func handleHook(logger *log.Logger) error {
	container, err := getContainerMetadata()
	if err != nil {
		logger.Printf(fmt.Sprintf("Could not decode container state: %v\n", err))
		return err
	}
	logger.Printf("Read container metadata")

	bundler := bundle.New(logger)
	ociSpec, err := bundler.ReadOciSpecFromBundle(container.bundle)
	if err != nil {
		logger.Printf(fmt.Sprintf("Could not get oci spec bundle: %v\n", err))
		return err
	}
	logger.Printf("Read oci spec")

	err = mount.MountDirectoryToContainer(ociSpec.Root.Path)
	if err != nil {
		logger.Printf(fmt.Sprintf("Could not mount directory to container: %v\n", err))
	}
	logger.Printf("Mounted directory into the container")

	return nil
}

func main() {
	logger := logger.New("runai-oci-hook")

	err := handleHook(logger)
	if err != nil {
		logger.Printf("Could not handle hook. Keep running without any changes")
	}
}
