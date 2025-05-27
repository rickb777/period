// See https://magefile.org/

//go:build mage

// Build steps for the period API:
package main

import (
	"github.com/magefile/mage/sh"
	"log"
	"os"
)

var Default = Build

func Build() error {
	if err := sh.RunV("go", "test", "-covermode=count", "-coverprofile=period.out", "."); err != nil {
		return err
	}
	if err := sh.RunV("go", "tool", "cover", "-func=period.out"); err != nil {
		return err
	}
	if err := sh.RunV("gofmt", "-l", "-w", "-s", "."); err != nil {
		return err
	}
	if err := sh.RunV("go", "vet", "./..."); err != nil {
		return err
	}
	return nil
}

// tests the module on both amd64 and i386 architectures
func CrossCompile() error {
	for _, arch := range []string{"amd64", "386"} {
		log.Printf("Testing on %s\n", arch)
		env := map[string]string{"GOARCH": arch}
		if _, err := sh.Exec(env, os.Stdout, os.Stderr, "go", "test", "./..."); err != nil {
			return err
		}
		log.Printf("%s is good.\n\n", arch)
	}
	return nil
}
