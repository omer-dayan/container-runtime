package patcher

import (
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/run-ai/runai-container-runtime/pkg/bundle"
	"log"
)

const (
	srcContainerToolkitDirectory = "/var/lib/runai/shared/"
	dstContainerToolkitDirectory = "/runai/shared/"
)

type Patcher struct {
	logger *log.Logger
	bundle *bundle.Bundle
}

func New(logger *log.Logger, bundler *bundle.Bundle) *Patcher {
	return &Patcher{
		logger: logger,
		bundle: bundler,
	}
}

func (p *Patcher) addMountIfNotExists(spec *specs.Spec) {
	if len(spec.Mounts) != 0 {
		for _, mount := range spec.Mounts {
			if mount.Destination == dstContainerToolkitDirectory {
				p.logger.Printf(fmt.Sprintf("mount to %v exists in OCI file: {%v}\n", dstContainerToolkitDirectory, mount))
				return
			}
		}
	}

	newMount := specs.Mount{
		Source:      srcContainerToolkitDirectory,
		Destination: dstContainerToolkitDirectory,
		Options:     []string{"bind"},
	}
	spec.Mounts = append(spec.Mounts, newMount)
}

func (p *Patcher) addPocEnvVar(spec *specs.Spec) {
	spec.Process.Env = append(spec.Process.Env, "POC=success")
}

func (p *Patcher) AddPatches(bundleDirectory string) error {
	ociSpec, err := p.bundle.ReadOciSpecFromBundle(bundleDirectory)
	if err != nil {
		return err
	}

	p.addPocEnvVar(ociSpec)
	p.logger.Printf(fmt.Sprintf("Added env var as a POC [%v]\n", bundleDirectory))

	return p.bundle.WriteOciSpecToBundle(bundleDirectory, ociSpec)
}
