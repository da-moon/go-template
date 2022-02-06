//go:build tools
// +build tools

package tools

// Manage tool dependencies via go.mod.
//
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/golang/go/issues/25922
//
// nolint
import (
	_ "github.com/davidrjenni/reftools/cmd/fillstruct"
	_ "github.com/fatih/gomodifytags"
	_ "github.com/fatih/motion"
	_ "github.com/go-delve/delve/cmd/dlv"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/josharian/impl"
	_ "github.com/jstemmer/gotags"
	_ "github.com/kisielk/errcheck"
	_ "github.com/klauspost/asmfmt/cmd/asmfmt"
	_ "github.com/koron/iferr"
	_ "github.com/magefile/mage"
	_ "github.com/rogpeppe/godef"
	_ "github.com/stretchr/gorc"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/tools/cmd/gorename"
	_ "golang.org/x/tools/cmd/guru"
	_ "golang.org/x/tools/gopls"
	_ "honnef.co/go/tools/cmd/keyify"
	_ "honnef.co/go/tools/cmd/staticcheck"
)

//go:generate go install -v "github.com/davidrjenni/reftools/cmd/fillstruct"
//go:generate go install -v "github.com/fatih/gomodifytags"
//go:generate go install -v "github.com/fatih/motion"
//go:generate go install -v "github.com/go-delve/delve/cmd/dlv"
//go:generate go install -v "github.com/josharian/impl"
//go:generate go install -v "github.com/jstemmer/gotags"
//go:generate go install -v "github.com/kisielk/errcheck"
//go:generate go install -v "github.com/klauspost/asmfmt/cmd/asmfmt"
//go:generate go install -v "github.com/koron/iferr"
//go:generate go install -v "github.com/magefile/mage"
//go:generate go install -v "github.com/rogpeppe/godef"
//go:generate go install -v "github.com/stretchr/gorc"
//go:generate go install -v "golang.org/x/lint/golint"
//go:generate go install -v "golang.org/x/tools/cmd/goimports"
//go:generate go install -v "golang.org/x/tools/cmd/gorename"
//go:generate go install -v "golang.org/x/tools/cmd/guru"
//go:generate go install -v "golang.org/x/tools/gopls"
//go:generate go install -v "honnef.co/go/tools/cmd/keyify"
//go:generate go install -v "honnef.co/go/tools/cmd/staticcheck"
//go:generate go install -v "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
