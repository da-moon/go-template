package value_test

import (
	"flag"
	"reflect"
	"testing"

	value "github.com/da-moon/go-template/internal/cli/value"
)

func TestAppendSlice_implements(t *testing.T) {
	t.Parallel()
	var raw interface{}
	raw = new(value.AppendSlice)
	if _, ok := raw.(flag.Value); !ok {
		t.Fatalf("AppendSlice should be a Value")
	}
}

func TestAppendSliceSet(t *testing.T) {
	t.Parallel()
	sv := new(value.AppendSlice)
	err := sv.Set("foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = sv.Set("bar")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{"foo", "bar"}
	if !reflect.DeepEqual([]string(*sv), expected) {
		t.Fatalf("Bad: %#v", sv)
	}
}
