# vagrant-exec

[![Build Status](https://travis-ci.org/dominodatalab/vagrant-exec.svg?branch=master)](https://travis-ci.org/dominodatalab/vagrant-exec)
[![Go Report Card](https://goreportcard.com/badge/github.com/dominodatalab/vagrant-exec)](https://goreportcard.com/report/github.com/dominodatalab/vagrant-exec)
[![GoDoc](https://godoc.org/github.com/dominodatalab/vagrant-exec?status.svg)](https://godoc.org/github.com/dominodatalab/vagrant-exec)

Golang wrapper around Vagrant CLI.

## Usage

```go
package main

import (
	"fmt"

	ve "github.com/dominodatalab/vagrant-exec"
)

func main() {
	// point to directory where the Vagrantfile is located and enable debug logging
	vagrant := ve.New("/path/to/Vagrantfile/directory", true)

	// check the install version
	version, err := vagrant.Version()
	if err != nil {
		panic(err)
	}
	fmt.Println(version)

	// create and provision VMs
	if err := vagrant.Up(); err != nil {
		panic(err)
	}

	// query the status of all VMs
	statusList, err := vagrant.Status()
	if err != nil {
		panic(err)
	}
	for _, status := range statusList {
		fmt.Printf("%#v", status)
	}

	// stop the VMs
	if err := vagrant.Halt(); err != nil {
		panic(err)
	}

	// destroy the VMs
	if err := vagrant.Destroy(); err != nil {
		panic(err)
	}

	// install a plugin
	plugin := ve.Plugin{
		Name:     "vagrant-disksize",
		Version:  "0.1.3",
		Location: "local",
	}
	if err := vagrant.PluginInstall(plugin); err != nil {
		panic(err)
	}

	// list all plugins
	plugins, err := vagrant.PluginList()
	if err != nil {
		panic(err)
	}
	for _, plugin := range plugins {
		fmt.Printf("%#v", plugin)
	}
}
```

## Contributions

Any suggestions and/or contributions are appreciated. Please submit an issue or PR with your suggested changes.
