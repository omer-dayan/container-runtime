package patcher

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestAddMountIfNotExistsNotExists(t *testing.T) {
	patcher := New(log.Default())
	spec := specs.Spec{}

	patcher.addMountIfNotExists(&spec)

	assert.Equal(t, 1, len(spec.Mounts))
	firstMount := spec.Mounts[0]
	assert.Equal(t, srcContainerToolkitDirectory, firstMount.Source)
	assert.Equal(t, dstContainerToolkitDirectory, firstMount.Destination)
}

func TestAddMountIfNotExistsExists(t *testing.T) {
	patcher := New(log.Default())
	spec := specs.Spec{
		Mounts: []specs.Mount{
			{Source: srcContainerToolkitDirectory,
				Destination: dstContainerToolkitDirectory},
		},
	}
	mountNumberBeforePatches := len(spec.Mounts)

	patcher.addMountIfNotExists(&spec)

	assert.Equal(t, mountNumberBeforePatches, len(spec.Mounts))
}

func TestAddMountIfNotExistsOtherExists(t *testing.T) {
	patcher := New(log.Default())
	spec := specs.Spec{
		Mounts: []specs.Mount{
			{Source: "TEST_SOURCE",
				Destination: "TEST_DESTINATION"},
		},
	}
	mountNumberBeforePatches := len(spec.Mounts)

	patcher.addMountIfNotExists(&spec)

	assert.Equal(t, mountNumberBeforePatches + 1, len(spec.Mounts))
	secondMount := spec.Mounts[1]
	assert.Equal(t, srcContainerToolkitDirectory, secondMount.Source)
	assert.Equal(t, dstContainerToolkitDirectory, secondMount.Destination)
}
