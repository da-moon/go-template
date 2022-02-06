package compile

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	figure "github.com/common-nighthawk/go-figure"
	git "github.com/da-moon/go-template/build/mage/git"
	primitives "github.com/da-moon/go-template/internal/primitives"
	color "github.com/fatih/color"
	hashimultierr "github.com/hashicorp/go-multierror"
	mg "github.com/magefile/mage/mg"
	sh "github.com/magefile/mage/sh"
	stacktrace "github.com/palantir/stacktrace"
)

const (
	packageName = "github.com/da-moon/go-template"
)

func mkdir() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := "bin"
	path = primitives.PathJoin(wd, path)
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		err = stacktrace.Propagate(err, "could not open a file descriptor at '%s'", path)
		return err
	}
	err = f.Sync()
	if err != nil {
		_ = f.Close()
		err = stacktrace.Propagate(err, "could not flush path '%s' to disk", path)
		return err
	}
	err = f.Close()
	if err != nil {
		err = stacktrace.Propagate(err, "could not close file descriptor at '%s' to disk", path)
		return err
	}
	return nil
}
func listDirs(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

// Target compiles all binaries stored in cmd.
// nolint:gocognit,gocyclo // this is a build function so we can ignore
// these two warnings
// [ TODO ] make it concurrent
func Target() error {
	mg.Deps(mkdir)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := "cmd"
	path = primitives.PathJoin(wd, path)
	binaries, err := listDirs(path)
	if err != nil {
		return err
	}
	env := map[string]string{
		"GO111MODULE": "on",
		"CGO_ENABLED": "0",
		"CGO_LDFLAGS": "-s -w -extldflags '-static'",
	}
	var result error
	for _, bin := range binaries {
		banner := figure.NewFigure(bin, "", true)
		banner.Print()
		fmt.Println("")
		versionpkg := fmt.Sprintf("%s/build/version", packageName)
		ldflags := make([]string, 0)
		buildDate := time.Now().Format("01/02/06")
		ldflags = append(ldflags, []string{
			"-X",
			fmt.Sprintf("%s.BuildDate=%v", versionpkg, buildDate),
		}...)
		buildUser, err := user.Current()
		if err == nil && buildUser != nil {
			ldflags = append(ldflags, []string{
				"-X",
				fmt.Sprintf("%s.BuildUser=%v", versionpkg, buildUser.Username),
			}...)
		}
		repo, err := git.Open()
		if err == nil && repo != nil {
			branch, err := repo.Branch()
			if err == nil && branch != "" {
				ldflags = append(ldflags, []string{
					"-X",
					fmt.Sprintf("%s.Branch=%v", versionpkg, branch),
				}...)
			}
			revision, err := repo.Commit()
			if err == nil && revision != "" {
				ldflags = append(ldflags, []string{
					"-X",
					fmt.Sprintf("%s.Revision=%v", versionpkg, revision),
				}...)
			}
			version, err := repo.Tag()
			if err == nil && version != "" {
				ldflags = append(ldflags, []string{
					"-X",
					fmt.Sprintf("%s.Version=%v", versionpkg, version),
				}...)
			}
		}
		args := []string{
			"build",
			"-ldflags",
			strings.Join(ldflags, " "),
			"-o",
			primitives.PathJoin(wd, "bin", bin),
			primitives.PathJoin(wd, "cmd", bin),
		}
		color.Green("# Build Command ------------------------------------------")
		fmt.Println("")
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "go build")
		fmt.Fprintf(&buf, " \\\n")
		fmt.Fprintf(&buf, "  -ldflags \\\n")
		fmt.Fprintf(&buf, "  '\n")
		for idx, val := range ldflags {
			if idx%2 == 0 {
				fmt.Fprintf(&buf, "  ")
			}
			fmt.Fprintf(&buf, `%s`, val)
			if idx%2 == 0 {
				fmt.Fprintf(&buf, ` "`)
			}
			if idx%2 == 1 {
				// if idx != len(ldflags)-1 {
				fmt.Fprintf(&buf, "\" \n")
				// }
			}
		}
		fmt.Fprintf(&buf, "  '")
		fmt.Fprintf(&buf, " -o %s", primitives.PathJoin(wd, "bin", bin))
		fmt.Fprintf(&buf, "  %s", primitives.PathJoin(wd, "cmd", bin))
		fmt.Println(buf.String())
		err = sh.RunWithV(env, "go", args...)
		if err != nil {
			hashimultierr.Append(result, err)
			continue
		}
		binaryPath := primitives.PathJoin(wd, "bin", bin)
		strip, err := exec.LookPath("strip")
		if strip != "" && err == nil {
			err = sh.Run("strip", binaryPath)
			if err != nil {
				err = stacktrace.Propagate(err, "failed to run strip on %q", binaryPath)
				hashimultierr.Append(result, err)
				continue
			}
		}
		upx, err := exec.LookPath("upx")
		if upx != "" && err == nil {
			err = sh.Run("upx", binaryPath)
			if err != nil {
				err = stacktrace.Propagate(err, "failed to run upx on %q", binaryPath)
				hashimultierr.Append(result, err)
			}
		}
	}
	return result
}
