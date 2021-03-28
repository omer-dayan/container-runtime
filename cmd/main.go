package main

import (
	"fmt"
	"github.com/run-ai/runai-container-runtime/pkg/patcher"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"
)

const (
	iocConfigJsonFileName = "config.json"
	runcCreateCommand = "create"
	originRuntimeBinary = "nvidia-container-runtime"
	logFilePath = "/var/log/patcher-container-runtime.log"
)

func getCommandAndBundle(args []string) (bool, string, error) {
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

func getIocConfigJsonFileFromBundleDirectory(bundleDirectory string) (string, error) {
	iocConfigJsonFilePath := path.Join(bundleDirectory, iocConfigJsonFileName)
	_, err := os.Stat(iocConfigJsonFilePath)
	if err != nil {
		return "", err
	}
	return iocConfigJsonFilePath, nil
}

func execOrigRuntime() error {
	originRuntimeBinaryFullPath, err := exec.LookPath(originRuntimeBinary)
	if err != nil {
		return err
	}

	err = syscall.Exec(originRuntimeBinaryFullPath, append([]string{originRuntimeBinaryFullPath}, os.Args[1:]...), os.Environ())
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var logger *log.Logger
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		logger = log.New(logFile, "[patcher-container-runtime]", log.LstdFlags)
	} else {
		logger = log.Default()
	}

	isCreateCommand, bundle, err := getCommandAndBundle(os.Args)
	if err != nil {
		logger.Printf(fmt.Sprintf("Could not find bundle or create command from args {%v} due to error: %v\n", os.Args, err))
		logger.Println("Falling back to original runtime")
	}

	if isCreateCommand {
		iocFilePath, err := getIocConfigJsonFileFromBundleDirectory(bundle)
		if err != nil {
			logger.Printf(fmt.Sprintf("Could not find OCI Config.json file due to error: %v\n", err))
		}

		patcher := patcher.New(logger)
		err = patcher.AddPatches(iocFilePath)
		if err != nil {
			logger.Printf(fmt.Sprintf("Could not patch OCI file due to error: %v\n", err))
		}
	}

	err = execOrigRuntime()
	if err != nil {
		logger.Printf(fmt.Sprintf("Could not run original runtime due to error: %v\n", err))
	}
}
