package decoder

import (
	"reflect"

	value "github.com/da-moon/go-template/internal/cli/value"
	mapstructure "github.com/mitchellh/mapstructure"
	stacktrace "github.com/palantir/stacktrace"
)

var Hooks = mapstructure.ComposeDecodeHookFunc(
	BoolToBoolFunc(),
	StringToDurationFunc(),
	StringToStringFunc(),
	Float64ToUintFunc(),
)

// BoolToBoolFunc is a mapstructure hook that looks for an incoming bool
// mapped to a Bool and does the translation.
func BoolToBoolFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Bool {
			return data, nil
		}

		val := value.Bool{}
		if t != reflect.TypeOf(val) {
			return data, nil
		}
		val.RawSet(data.(bool))
		return val, nil
	}
}

// StringToDurationFunc is a mapstructure hook that looks for an incoming
// string mapped to a Duration and does the translation.
func StringToDurationFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		val := value.Duration{}
		if t != reflect.TypeOf(val) {
			return data, nil
		}
		if err := val.Set(data.(string)); err != nil {
			return nil, err
		}
		return val, nil
	}
}

// StringToStringFunc is a mapstructure hook that looks for an incoming
// string mapped to a String and does the translation.
func StringToStringFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		val := value.String{}
		if t != reflect.TypeOf(val) {
			return data, nil
		}
		val.RawSet(data.(string))
		return val, nil
	}
}

// Float64ToUintFunc is a mapstructure hook that looks for an incoming
// float64 mapped to a Uint and does the translation.
func Float64ToUintFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Float64 {
			return data, nil
		}

		val := value.Uint{}
		if t != reflect.TypeOf(val) {
			return data, nil
		}

		fv := data.(float64)
		if fv < 0 {
			return nil, stacktrace.NewError("value cannot be negative")
		}

		// The standard guarantees at least this, and this is fine for
		// values we expect to use in configs vs. being fancy with the
		// machine's size for uint.
		if fv > (1<<32 - 1) {
			return nil, stacktrace.NewError("value is too large")
		}
		val.RawSet((uint)(fv))
		return val, nil
	}
}
