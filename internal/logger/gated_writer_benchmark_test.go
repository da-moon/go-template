package logger_test

import (
	"bytes"
	"strconv"
	"testing"

	logger "github.com/da-moon/go-template/internal/logger"
	test "github.com/da-moon/go-template/internal/testutils"
	serf "github.com/hashicorp/serf/cmd/serf/command/agent"
)

func BenchmarkSmallWriteGoLogger(b *testing.B) {
	cases := test.ConcurrentCases()
	b.ReportAllocs()
	b.ResetTimer()
	const n = 4 << 10
	for _, concurrency := range cases {
		b.SetParallelism(concurrency)
		b.Run("cores "+strconv.Itoa(concurrency), func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				b.ResetTimer()
				b.SetBytes(n)
				backing := make([]byte, 0)
				buf := bytes.NewBuffer(backing)
				logger := logger.NewGatedWriter(buf)
				for pb.Next() {
					for i := 0; i < 1024; i++ {
						bsWrite := []byte("abcd\n")
						if n, err := logger.Write(bsWrite); err != nil {
							errStr := "wrote " + strconv.Itoa(n) + " bytes want " + strconv.Itoa(len(bsWrite)) + " bytes, err:" + err.Error()
							print(errStr)
						}
						logger.Flush()
					}
					buf.Reset()
				}
			})
		})
	}
}
func BenchmarkSmallWriteSerf(b *testing.B) {
	cases := test.ConcurrentCases()
	b.ReportAllocs()
	const n = 4 << 10
	for _, concurrency := range cases {
		b.SetParallelism(concurrency)
		b.Run("cores "+strconv.Itoa(concurrency), func(b *testing.B) {
			b.ResetTimer()
			b.SetBytes(n)
			b.RunParallel(func(pb *testing.PB) {
				backing := make([]byte, 0)
				buf := bytes.NewBuffer(backing)
				logger := serf.GatedWriter{Writer: buf}
				for pb.Next() {
					// b.ResetTimer()
					for i := 0; i < 1024; i++ {
						bsWrite := []byte("abcd\n")
						if n, err := logger.Write(bsWrite); err != nil {
							errStr := "wrote " + strconv.Itoa(n) + " bytes want " + strconv.Itoa(len(bsWrite)) + " bytes, err:" + err.Error()
							print(errStr)
						}
						logger.Flush()
					}
					buf.Reset()
				}
			})
		})
	}
}
