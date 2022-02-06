package flagset

import (
	"flag"
	"fmt"
	"io"
	"strings"

	text "github.com/kr/text"
)

// PrintTitle prints a consistently-formatted title to the given writer.
func PrintTitle(w io.Writer, s string) {
	fmt.Fprintf(w, "%s\n\n", s)
}

// PrintFlag prints a single flag to the given writer.
func PrintFlag(w io.Writer, f *flag.Flag) {
	example, _ := flag.UnquoteUsage(f)
	if example != "" {
		fmt.Fprintf(w, "  -%s=<%s>\n", f.Name, example)
	} else {
		fmt.Fprintf(w, "  -%s\n", f.Name)
	}

	indented := WrapAtLength(f.Usage, 5)
	fmt.Fprintf(w, "%s\n\n", indented)
}

// maxLineLength is the maximum width of any line.
const maxLineLength int = 72

// WrapAtLength wraps the given text at the maxLineLength, taking into account
// any provided left padding.
func WrapAtLength(s string, pad int) string {
	wrapped := text.Wrap(s, maxLineLength-pad)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", pad) + line
	}
	return strings.Join(lines, "\n")
}
