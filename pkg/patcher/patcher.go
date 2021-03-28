package patcher

import (
	"encoding/json"
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	prestartHookBinary = "patcher-hook"
	srcContainerToolkitDirectory = "/var/lib/runai/shared/"
	dstContainerToolkitDirectory = "/runai/shared/"
)

type Patcher struct {
	logger *log.Logger
}

func New(logger *log.Logger) *Patcher {
	return &Patcher{
		logger: logger,
	}
}

func (p *Patcher) addHookIfNotExists(spec *specs.Spec) error {
	if spec.Hooks == nil {
		spec.Hooks = &specs.Hooks{}
	} else if len(spec.Hooks.Prestart) != 0 {
		for _, preStartHook := range spec.Hooks.Prestart {
			if strings.Contains(preStartHook.Path, prestartHookBinary) {
				p.logger.Printf(fmt.Sprintf("Binary %v exists as prestart hook: %v\n", prestartHookBinary, preStartHook))
				return nil
			}
		}
	}

	prestartHookBinaryFullPath, err := exec.LookPath(prestartHookBinary)
	if err != nil {
		return err
	}

	hook := specs.Hook{
		Path: prestartHookBinaryFullPath,
		Args: []string{prestartHookBinaryFullPath, "prestart"},
	}
	spec.Hooks.Prestart = append(spec.Hooks.Prestart, hook)
	return nil
}

func (p *Patcher) addMountIfNotExists(spec *specs.Spec) {
	if len(spec.Mounts) != 0 {
		for _, mount := range spec.Mounts {
			if mount.Source == srcContainerToolkitDirectory && mount.Destination == dstContainerToolkitDirectory {
				p.logger.Printf(fmt.Sprintf("Mount exists in OCI file: {%v}\n", mount))
				return
			}
		}
	}

	newMount := specs.Mount{
		Source: srcContainerToolkitDirectory,
		Destination: dstContainerToolkitDirectory,
		Options: []string{"bind"},
	}
	spec.Mounts = append(spec.Mounts, newMount)
}

func (p *Patcher) AddPatches(ociFilePath string) error {
	ociFile, err := os.OpenFile(ociFilePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer ociFile.Close()
	p.logger.Printf(fmt.Sprintf("Opened oci file [%v]\n", ociFile))

	ociJson, err := ioutil.ReadAll(ociFile)
	if err != nil {
		return err
	}
	p.logger.Printf(fmt.Sprintf("Read oci file data [%v]\n", ociFilePath))

	var ociSpec specs.Spec
	err = json.Unmarshal(ociJson, &ociSpec)
	if err != nil {
		return err
	}
	p.logger.Printf(fmt.Sprintf("Parsed oci json as object [%v]\n", ociFilePath))

	p.addMountIfNotExists(&ociSpec)
	p.logger.Printf(fmt.Sprintf("Added mount if was not exist [%v]\n", ociFilePath))

	ociJson, err = json.Marshal(ociSpec)
	if err != nil {
		return err
	}
	p.logger.Printf(fmt.Sprintf("Parsed oci object as json [%v]\n", ociFilePath))

	_, err = ociFile.WriteAt(ociJson, 0)
	if err != nil {
		return err
	}
	p.logger.Printf(fmt.Sprintf("Wrote oci json to file [%v]\n", ociFilePath))

	return nil
}
