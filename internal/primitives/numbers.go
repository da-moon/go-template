package primitives

import (
	"fmt"
	"strconv"
	"time"
)

// RoundUp takes a uint64 greater than 0 and rounds it up to the next
// power of 2.
func RoundUp(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

// IsInt checks if a string is an integer
func IsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// IsFloat checks if a string is a float number
func IsFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// Utoa transforms a uint into a string
func Utoa(u uint) string {
	return fmt.Sprint(u)
}

// consts -
const (
	MaxUint = ^uint(0)
	MaxInt  = int(MaxUint >> 1)
)

// DurationToInt -
func DurationToInt(duration, unit time.Duration) int {
	durationAsNumber := duration / unit

	if int64(durationAsNumber) > int64(MaxInt) {
		return MaxInt
	}
	return int(durationAsNumber)
}
