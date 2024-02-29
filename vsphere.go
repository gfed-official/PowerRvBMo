package main

import (
	"net/url"
    "fmt"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
	"golang.org/x/net/context"
)

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

func Revert(vmName string) error {
	vm, err := finder.VirtualMachine(ctx, vmName)
	if err != nil {
		return err
	}

	task, err := vm.RevertToCurrentSnapshot(ctx, true)
	if err != nil {
		return err
	}

	err = task.WaitEx(ctx)
	if err != nil {
		return err
	}

    task, err = vm.PowerOn(ctx)
    if err != nil {
        return err
    }

    err = task.WaitEx(ctx)
    if err != nil {
        return err
    }

	return nil
}

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
