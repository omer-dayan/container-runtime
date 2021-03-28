package main

import (
	"fmt"
	"github.com/run-ai/runai-container-runtime/pkg/bundle"
	"github.com/run-ai/runai-container-runtime/pkg/logger"
	"github.com/run-ai/runai-container-runtime/pkg/patcher"
	"os"
	"os/exec"
	"syscall"
)

const (
	runcCreateCommand = "create"
	originRuntimeBinary = "nvidia-container-runtime"
)

func getCommandAndBundlePath(args []string) (bool, string, error) {
	isCommand := false
	var bundle string
	for i, arg := range args {
		if arg == runcCreateCommand {
			isCommand = true
		} else if arg == "--bundle" || arg == "-b" {
			if len(args) == i + 1 {
				return false, "", fmt.Errorf("could not find bundle path after bundle flag in args {%v}", args)
			}
			bundle = args[i + 1]
		}
	}

	if isCommand && bundle == "" {
		currentWorkingDirectory, err := os.Getwd()
		if err != nil {
			return false, "", err
		}
		bundle = currentWorkingDirectory
	}
	return isCommand, bundle, nil
}

func execOrigRuntime() error {
	originRuntimeBinaryFullPath, err := exec.LookPath(originRuntimeBinary)
	if err != nil {
		return err
	}

	return syscall.Exec(originRuntimeBinaryFullPath, append([]string{originRuntimeBinaryFullPath}, os.Args[1:]...), os.Environ())
}

func main() {
	logger := logger.New("runai-container-runtime")

	isCreateCommand, bundlePath, err := getCommandAndBundlePath(os.Args)
	if err != nil {
		logger.Printf(fmt.Sprintf("Could not find bundle or create command from args {%v} due to error: %v\n", os.Args, err))
		logger.Println("Falling back to original runtime")
	}

	if isCreateCommand {
		patcher := patcher.New(logger, bundle.New(logger))
		err = patcher.AddPatches(bundlePath)
		if err != nil {
			logger.Printf(fmt.Sprintf("Could not patch OCI file due to error: %v\n", err))
		}
	}

	err = execOrigRuntime()
	if err != nil {
		logger.Printf(fmt.Sprintf("Could not run original runtime due to error: %v\n", err))
	}
}
