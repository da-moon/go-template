package version

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// Build information. Populated at build-time.
var (
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	Toolchain = runtime.Version()
	Arch      = runtime.GOARCH
)

// BuildInformation helps with working with immutable build information.
//go:generate gomodifytags -override -file $GOFILE -struct BuildInformation -add-tags json,mapstructure -w -transform snakecase
type BuildInformation struct {
	Version   string `json:"version" mapstructure:"version"`
	Revision  string `json:"revision" mapstructure:"revision"`
	Branch    string `json:"branch" mapstructure:"branch"`
	BuildDate string `json:"build_date" mapstructure:"build_date"`
	Toolchain string `json:"toolchain" mapstructure:"toolchain"`
	BuildUser string `json:"build_user" mapstructure:"build_user"`
	Arch      string `json:"arch" mapstructure:"arch"`
}

// New returns a new BuildInformation.
func New() *BuildInformation {
	result := &BuildInformation{
		Version:   Version,
		Revision:  Revision,
		Branch:    Branch,
		BuildDate: BuildDate,
		Arch:      Arch,
		Toolchain: Toolchain,
		BuildUser: BuildUser,
	}
	return result
}

// Print returns version information.
func (b *BuildInformation) ToString() string {
	// versionInfoTmpl contains the template used by Info.
	var result bytes.Buffer
	if b.Version != "" {
		fmt.Fprintf(&result, "\nVersion      :  %s", strings.TrimPrefix(b.Version, "v"))
	}
	if b.Revision != "" {
		fmt.Fprintf(&result, "\nRevision     :  %s", b.Revision)
	}
	if b.Branch != "" {
		fmt.Fprintf(&result, "\nBranch       :  %s", b.Branch)
	}
	if b.Arch != "" {
		fmt.Fprintf(&result, "\nArchitecture :  %s", b.Arch)
	}
	if b.BuildUser != "" {
		fmt.Fprintf(&result, "\nBuild User   :  %s", b.BuildUser)
	}
	if b.BuildDate != "" {
		fmt.Fprintf(&result, "\nBuild Date   :  %s", b.BuildDate)
	}
	if b.Toolchain != "" {
		fmt.Fprintf(&result, "\nToolchain    :  %s", b.Toolchain)
	}
	return result.String()
}

// Info returns version, branch and revision information.
func (b *BuildInformation) Info() string {
	return fmt.Sprintf("(version=%s, branch=%s, revision=%s)", b.Version, b.Branch, b.Revision)
}

// BuildContext returns toolchain, buildUser and buildDate information.
func (b *BuildInformation) BuildContext() string {
	return fmt.Sprintf("(go=%s,arch=%s, user=%s, date=%s)", b.Toolchain, b.Arch, b.BuildUser, b.BuildDate)
}
