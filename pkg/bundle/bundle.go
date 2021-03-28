package bundle

import (
	"encoding/json"
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const iocConfigJsonFileName = "config.json"

type Bundle struct {
	logger *log.Logger
}

func New(logger *log.Logger) *Bundle {
	return &Bundle{
		logger: logger,
	}
}

func (b *Bundle) getOciConfigJsonFileFromBundleDirectory(bundleDirectory string) (string, error) {
	iocConfigJsonFilePath := path.Join(bundleDirectory, iocConfigJsonFileName)
	_, err := os.Stat(iocConfigJsonFilePath)
	if err != nil {
		return "", err
	}
	return iocConfigJsonFilePath, nil
}

func (b *Bundle) getOciFile(bundleDirectory string) (*os.File, error) {
	configJsonFile, err := b.getOciConfigJsonFileFromBundleDirectory(bundleDirectory)
	if err != nil {
		return nil, err
	}
	b.logger.Printf(fmt.Sprintf("Got oci file from bundle [%v]\n", configJsonFile))

	ociFile, err := os.OpenFile(configJsonFile, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return ociFile, nil
}

func (b *Bundle) ReadOciSpecFromBundle(bundleDirectory string) (*specs.Spec, error) {
	ociFile, err := b.getOciFile(bundleDirectory)
	if err != nil {
		return nil, err
	}
	defer ociFile.Close()
	b.logger.Printf(fmt.Sprintf("Opened oci file [%v]\n", ociFile))

	ociJson, err := ioutil.ReadAll(ociFile)
	if err != nil {
		return nil, err
	}
	b.logger.Printf(fmt.Sprintf("Read oci file data [%v]\n", bundleDirectory))

	var ociSpec specs.Spec
	err = json.Unmarshal(ociJson, &ociSpec)
	if err != nil {
		return nil, err
	}
	b.logger.Printf(fmt.Sprintf("Parsed oci json as object [%v]\n", bundleDirectory))

	return &ociSpec, nil
}

func (b *Bundle) WriteOciSpecToBundle(bundleDirectory string, ociSpec *specs.Spec) error {
	ociFile, err := b.getOciFile(bundleDirectory)
	if err != nil {
		return err
	}
	defer ociFile.Close()
	b.logger.Printf(fmt.Sprintf("Opened oci file [%v]\n", ociFile))

	ociJson, err := json.Marshal(ociSpec)
	if err != nil {
		return err
	}
	b.logger.Printf(fmt.Sprintf("Parsed oci object as json [%v]\n", bundleDirectory))

	_, err = ociFile.WriteAt(ociJson, 0)
	if err != nil {
		return err
	}
	b.logger.Printf(fmt.Sprintf("Wrote oci json to file [%v]\n", bundleDirectory))

	return nil
}
