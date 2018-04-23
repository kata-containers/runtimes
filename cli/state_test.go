// Copyright (c) 2017 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"flag"
	"testing"

	vc "github.com/kata-containers/runtime/virtcontainers"
	vcAnnotations "github.com/kata-containers/runtime/virtcontainers/pkg/annotations"
	"github.com/kata-containers/runtime/virtcontainers/pkg/vcmock"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestStateCliAction(t *testing.T) {
	assert := assert.New(t)

	actionFunc, ok := stateCLICommand.Action.(func(ctx *cli.Context) error)
	assert.True(ok)

	flagSet := flag.NewFlagSet("flag", flag.ContinueOnError)

	// without container id
	flagSet.Parse([]string{"runtime"})
	ctx := cli.NewContext(&cli.App{}, flagSet, nil)
	err := actionFunc(ctx)
	assert.Error(err)

	// with container id
	flagSet.Parse([]string{"runtime", testContainerID})
	ctx = cli.NewContext(&cli.App{}, flagSet, nil)
	err = actionFunc(ctx)
	assert.Error(err)
}

func TestStateSuccessful(t *testing.T) {
	assert := assert.New(t)

	sandbox := &vcmock.Sandbox{
		MockID: testContainerID,
	}

	sandbox.MockContainers = []*vcmock.Container{
		{
			MockID:      sandbox.ID(),
			MockSandbox: sandbox,
		},
	}

	testingImpl.ContainerSandboxListFunc = func(containerID string) ([]string, bool, error) {
		return []string{}, false, nil
	}

	defer func() {
		testingImpl.ContainerSandboxListFunc = nil
	}()

	// trying with an inexistent id
	err := state("123456789")
	assert.Error(err)

	testingImpl.ContainerSandboxListFunc = func(containerID string) ([]string, bool, error) {
		return []string{sandbox.ID()}, true, nil
	}

	testingImpl.StatusContainerFunc = func(sandboxID, containerID string) (vc.ContainerStatus, error) {
		return vc.ContainerStatus{
			ID: sandbox.ID(),
			Annotations: map[string]string{
				vcAnnotations.ContainerTypeKey: string(vc.PodContainer),
			},
		}, nil
	}

	defer func() {
		testingImpl.StatusContainerFunc = nil
	}()

	err = state(sandbox.ID())
	assert.NoError(err)
}
