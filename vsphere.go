package main

import (
	"fmt"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
	"golang.org/x/net/context"
)

// Connect to vSphere and return client
func Connect() *govmomi.Client {
	ctx := context.Background()
	u, err := soap.ParseURL(tomlConf.VSphereURL)
	if err != nil {
		panic(err)
	}

	u.User = url.UserPassword(tomlConf.VSphereUser, tomlConf.VSpherePassword)
	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to vSphere")

	return c
}

// Revert VM based on name
func Revert(vmName string) (int, error) {
	vm, err := finder.VirtualMachine(ctx, vmName)
	if err != nil {
		return -1, err
	}

	task, err := vm.RevertToCurrentSnapshot(ctx, true)
	if err != nil {
		return -1, err
	}

	err = task.WaitEx(ctx)
	if err != nil {
		return -1, err
	}

    revertCount := incrementRevertCount(vmName)

	task, err = vm.PowerOn(ctx)
	if err != nil {
		return -1, err
	}

	err = task.WaitEx(ctx)
	if err != nil {
		return -1, err
	}

	return revertCount, nil
}

// Get VM names from a prefix (e.g. pod number)
func GetVMs(team string) ([]string, error) {
	team = team + "*"

	vms, err := finder.VirtualMachineList(ctx, team)
	if err != nil {
		return nil, err
	}

	var vmNames []string
	for _, vm := range vms {
		vmNames = append(vmNames, vm.Name())
	}

	return vmNames, nil
}
