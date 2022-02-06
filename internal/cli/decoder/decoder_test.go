package decoder_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	decoder "github.com/da-moon/go-template/internal/cli/decoder"
	value "github.com/da-moon/go-template/internal/cli/value"
	mapstructure "github.com/mitchellh/mapstructure"
)

func TestConfigConfigDecodeHook(t *testing.T) {
	t.Parallel()
	type config struct {
		B value.Bool     `mapstructure:"bool"`
		D value.Duration `mapstructure:"duration"`
		S value.String   `mapstructure:"string"`
		U value.Uint     `mapstructure:"uint"`
	}

	cases := []struct {
		in      string
		success string
		failure string
	}{
		{
			`{ }`,
			`"false" "0s" "" "0"`,
			"",
		},
		{
			`{ "bool": true, "duration": "2h", "string": "hello", "uint": 23 }`,
			`"true" "2h0m0s" "hello" "23"`,
			"",
		},
		{
			`{ "bool": "nope" }`,
			"",
			"got 'string'",
		},
		{
			`{ "duration": "nope" }`,
			"",
			`invalid duration "nope"`,
		},
		{
			`{ "string": 123 }`,
			"",
			"got 'float64'",
		},
		{
			`{ "uint": -1 }`,
			"",
			"value cannot be negative",
		},
		{
			`{ "uint": 4294967296 }`,
			"",
			"value is too large",
		},
	}
	for i, c := range cases {
		var raw interface{}
		dec := json.NewDecoder(bytes.NewBufferString(c.in))
		if err := dec.Decode(&raw); err != nil {
			t.Fatalf("(case %d) err: %v", i, err)
		}

		var r config
		msdec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook:  decoder.Hooks,
			Result:      &r,
			ErrorUnused: true,
		})
		if err != nil {
			t.Fatalf("(case %d) err: %v", i, err)
		}

		err = msdec.Decode(raw)
		if c.failure != "" {
			if err == nil || !strings.Contains(err.Error(), c.failure) {
				t.Fatalf("(case %d) err: %v", i, err)
			}
			continue
		}
		if err != nil {
			t.Fatalf("(case %d) err: %v", i, err)
		}

		actual := fmt.Sprintf("%q %q %q %q",
			r.B.String(),
			r.D.String(),
			r.S.String(),
			r.U.String())
		if actual != c.success {
			t.Fatalf("(case %d) bad: %s", i, actual)
		}
	}
}
