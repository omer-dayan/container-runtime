package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetCommandAndBundleCreateWithBundle(t *testing.T) {
	args := []string{"BIN", "create", "--bundle", "/test"}

	isCommand, bundle, err := getCommandAndBundle(args)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, isCommand)
	assert.Equal(t, "/test", bundle)
}

func TestGetCommandAndBundleCreateWithBundleShortcut(t *testing.T) {
	args := []string{"BIN", "create", "-b", "/test"}

	isCommand, bundle, err := getCommandAndBundle(args)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, isCommand)
	assert.Equal(t, "/test", bundle)
}

func TestGetCommandAndBundleNoCreateNoBundle(t *testing.T) {
	args := []string{"BIN", "delete", "test"}

	isCommand, bundle, err := getCommandAndBundle(args)

	assert.Equal(t, nil, err)
	assert.Equal(t, false, isCommand)
	assert.Equal(t, "", bundle)
}

func TestGetCommandAndBundleCreateNoBundle(t *testing.T) {
	args := []string{"BIN", "create", "test"}

	isCommand, bundle, err := getCommandAndBundle(args)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, isCommand)
	currentWorkingDirectory, _ := os.Getwd()
	assert.Equal(t, currentWorkingDirectory, bundle)
}
