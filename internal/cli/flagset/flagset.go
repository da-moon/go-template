package flagset

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
)

type FlagSet struct {
	value map[string]*flag.FlagSet
	usage string
	name  string
}

func New(name string, usage string) *FlagSet {
	result := &FlagSet{
		usage: usage,
		name:  name,
		value: map[string]*flag.FlagSet{
			name: flag.NewFlagSet(name, flag.ContinueOnError),
		},
	}
	return result
}
func (fs *FlagSet) Value() *flag.FlagSet {
	result := flag.NewFlagSet(fs.name, flag.ContinueOnError)
	for _, val := range fs.value {
		val.VisitAll(func(f *flag.Flag) {
			result.Var(f.Value, f.Name, f.Usage)
		})
	}
	return result
}
func (fs *FlagSet) Var(value flag.Value, name string, usage string) {
	_, ok := fs.value[fs.name]
	if !ok {
		fs.value[fs.name] = flag.NewFlagSet(fs.name, flag.ContinueOnError)
	}
	fs.value[fs.name].Var(value, name, usage)
}

func (fs *FlagSet) Merge(input *flag.FlagSet) {
	if input == nil {
		return
	}
	fs.value[input.Name()] = input
}

func (fs *FlagSet) Usage() string {
	out := new(bytes.Buffer)
	out.WriteString(strings.TrimSpace(fs.usage))
	out.WriteString("\n")
	out.WriteString("\n")
	for _, fs := range fs.value {
		PrintTitle(out, fmt.Sprintf("%s Options:", fs.Name()))
		fs.VisitAll(func(f *flag.Flag) {
			PrintFlag(out, f)
		})
	}
	return strings.TrimRight(out.String(), "\n")
}

// Contains returns true if the given flag is contained in the given flag
// set or false otherwise.
func Contains(fs *flag.FlagSet, f *flag.Flag) bool {
	if fs == nil {
		return false
	}
	var in bool
	fs.VisitAll(func(hf *flag.Flag) {
		in = in || f.Name == hf.Name
	})
	return in
}

func Merge(dst, src *flag.FlagSet) {
	if dst == nil {
		panic("dst cannot be nil")
	}
	if src == nil {
		return
	}
	src.VisitAll(func(f *flag.Flag) {
		dst.Var(f.Value, f.Name, f.Usage)
	})
}
