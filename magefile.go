//go:build mage
// +build mage

package main

import (
	"os"

	build "github.com/da-moon/go-template/internal/mage/build"
	test "github.com/da-moon/go-template/internal/mage/test"
	primitives "github.com/da-moon/go-template/internal/primitives"
	mg "github.com/magefile/mage/mg"
	sh "github.com/magefile/mage/sh"
	stacktrace "github.com/palantir/stacktrace"
)

var (
	// Default target to run when none is specified.
	Default = Build
	// Aliases can be used interchangeably with their targets.
	Aliases = map[string]interface{}{
		"d": Deps,
		"t": Test,
		"b": Build,
		"c": Clean,
	}
)

// Deps tidies go mod and downloads dependencies
func Deps() error {
	env := map[string]string{
		"GO111MODULE": "on",
		"CGO_ENABLED": "0",
		"CGO_LDFLAGS": "-s -w -extldflags '-static'",
	}
	args := []string{
		"mod",
		"tidy",
	}
	err := sh.RunWithV(env, "go", args...)
	if err != nil {
		err = stacktrace.Propagate(err, "could not tidy go modules")
		return err
	}
	args = []string{
		"get",
	}
	if mg.Verbose() {
		args = append(args, "-v")
	}
	args = append(args, "."+string(os.PathSeparator)+"...")
	err = sh.RunWithV(env, "go", args...)
	if err != nil {
		err = stacktrace.Propagate(err, "could not download dependant packages")
		return err
	}
	return nil

}

func Clean() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := "bin"
	path = primitives.PathJoin(wd, path)
	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		os.RemoveAll(path)
	}
	return nil
}
func Build() error {
	mg.Deps(Clean, Deps)
	return build.Target()
}

// Test will run all tests across all sub-directories once.
func Test() error {
	mg.Deps(Deps)
	return test.Target()
}
