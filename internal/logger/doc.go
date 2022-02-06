// Package logger this package combines various logger
// libraries from `hashicorp` and adds some extra features, safety checks and optimizations.
// ─── ACKNOWLEDGEMENTS ───────────────────────────────────────────────────────────
//  - github.com/hashicorp/serf/
//  - github.com/hashicorp/logutils/
// ─── BENCHMARKS ─────────────────────────────────────────────────────────────────
// `go-logger` implementation of gated writer vs [`Serf`](github.com/hashicorp/serf) implementation (go version `go1.15 linux/amd64`)
// ────────────────────────────────────────────────────────────────────────────────
// Running tool: /home/gitpod/go/bin/go test -benchmem -run=^$ github.com/da-moon/go-template/internal/logger -bench ^(BenchmarkSmallWriteGoLogger|BenchmarkSmallWriteSerf)$ -v
// goos: linux
// goarch: amd64
// pkg: github.com/da-moon/go-template/internal/logger
// BenchmarkSmallWriteGoLogger
// BenchmarkSmallWriteGoLogger/cores_1
// BenchmarkSmallWriteGoLogger/cores_1-16         	   52780	     27048 ns/op	 151.43 MB/s	    5427 B/op	    1017 allocs/op
// BenchmarkSmallWriteGoLogger/cores_2
// BenchmarkSmallWriteGoLogger/cores_2-16         	   46780	     26178 ns/op	 156.47 MB/s	    5254 B/op	     984 allocs/op
// BenchmarkSmallWriteGoLogger/cores_4
// BenchmarkSmallWriteGoLogger/cores_4-16         	   45724	     26675 ns/op	 153.55 MB/s	    5266 B/op	     987 allocs/op
// BenchmarkSmallWriteGoLogger/cores_8
// BenchmarkSmallWriteGoLogger/cores_8-16         	   44344	     28360 ns/op	 144.43 MB/s	    5426 B/op	    1017 allocs/op
// BenchmarkSmallWriteGoLogger/cores_16
// BenchmarkSmallWriteGoLogger/cores_16-16        	   36810	     29573 ns/op	 138.51 MB/s	    5422 B/op	    1015 allocs/op
// BenchmarkSmallWriteSerf
// BenchmarkSmallWriteSerf/cores_1
// BenchmarkSmallWriteSerf/cores_1-16             	   38850	     32975 ns/op	 124.21 MB/s	    5470 B/op	    1024 allocs/op
// BenchmarkSmallWriteSerf/cores_2
// BenchmarkSmallWriteSerf/cores_2-16             	   41787	     28999 ns/op	 141.25 MB/s	    5469 B/op	    1024 allocs/op
// BenchmarkSmallWriteSerf/cores_4
// BenchmarkSmallWriteSerf/cores_4-16             	   43066	     30472 ns/op	 134.42 MB/s	    5469 B/op	    1024 allocs/op
// BenchmarkSmallWriteSerf/cores_8
// BenchmarkSmallWriteSerf/cores_8-16             	   35313	     28614 ns/op	 143.15 MB/s	    5471 B/op	    1024 allocs/op
// BenchmarkSmallWriteSerf/cores_16
// BenchmarkSmallWriteSerf/cores_16-16            	   39213	     29698 ns/op	 137.92 MB/s	    5470 B/op	    1024 allocs/op
// PASS
// ok  	github.com/da-moon/go-template/internal/logger	16.458s
// ────────────────────────────────────────────────────────────────────────────────
package logger
