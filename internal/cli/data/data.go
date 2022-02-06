package data

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/palantir/stacktrace"
)

func Load(input string, testStdin io.Reader) (string, error) {
	// [ NOTE ] for handling empty quoted shell parameters
	if len(input) == 0 {
		return "", nil
	}
	switch input[0] {
	case '@':
		return fromFile(input[1:])
	case '-':
		if len(input) > 1 {
			return input, nil
		}
		return fromStandardInput(testStdin)
	default:
		return input, nil
	}
}
func fromFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		err = stacktrace.Propagate(err, "failed to read file")
		return "", err
	}
	return string(data), nil
}

func fromStandardInput(testStdin io.Reader) (string, error) {
	var stdin io.Reader = os.Stdin
	if testStdin != nil {
		stdin = testStdin
	}

	var b bytes.Buffer
	_, err := io.Copy(&b, stdin)
	if err != nil {
		err = stacktrace.Propagate(err, "Failed to read stdin")
		return "", err
	}
	return b.String(), nil
}
