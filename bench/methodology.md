# Performance Benchmark Methodology

This document outlines the engineering methodology, mathematical formulas, measurement procedures, and environment configurations used to evaluate the performance of the **native Go implementation** ([`github.com/DebadityaDas05/go-bytes`](https://github.com/DebadityaDas05/go-bytes)) compared to the original V8 JavaScript reference library and C-shared Koffi FFI runtime bindings.

---

## 1. System Overview & Architecture

The benchmark suite evaluates performance across three distinct execution layers:

1. **Native Go Library ([`src/bytes.go`](https://github.com/DebadityaDas05/go-bytes/blob/main/src/bytes.go))**: Compiled, zero-overhead native Go binary routines.
2. **Go C-Shared FFI Bridge ([`bridge/bridge.go`](https://github.com/DebadityaDas05/go-bytes/blob/main/bridge/bridge.go))**: Go exports compiled into a C-shared dynamic library (`libbytes.dll` / `libbytes.so`) invoked via Node.js Koffi FFI.
3. **JS Reference Implementation**: Original pure V8 JavaScript library executing on Node.js v20.x.

---

## 2. Benchmark Suite Components

All benchmarking code resides in the [`bench/`](https://github.com/DebadityaDas05/go-bytes/tree/main/bench) directory:

- **[`bench/bytes_bench_test.go`](https://github.com/DebadityaDas05/go-bytes/blob/main/bench/bytes_bench_test.go)**: Standard Go microbenchmark suite leveraging Go's `testing.B` harness with memory allocation profiling (`b.ReportAllocs()`).
- **[`bench/main.go`](https://github.com/DebadityaDas05/go-bytes/blob/main/bench/main.go)**: Standalone Go application that executes 500,000 iterations per target function, measures nanosecond-precision tail latencies ($P_{99}$), collects Go runtime memory allocation statistics (`runtime.ReadMemStats`), and writes results to [`bench/results.json`](https://github.com/DebadityaDas05/go-bytes/blob/main/bench/results.json).

---

## 3. Measured Metrics & Mathematical Definitions

### A. Startup Time (`startup_ms`)
- **Definition**: Time in milliseconds ($\text{ms}$) elapsed between process/module invocation and ready state for processing requests.
- **Go Native Measurement**:
  $$\text{Startup} = \frac{t_{\text{first\_call}} - t_{\text{init}}}{1,000,000} \quad (\text{ms})$$
  Captured using `time.Now()` nanosecond monotonic clock resolution.
- **JS / FFI Measurement**: Captures V8 module parsing, symbols binding, and dynamic C-shared library (`libbytes.dll`) loading time.

### B. Memory Allocation & Footprint (`memory_rss_mb`)
- **Heap Allocation (`go_native_alloc_mb`)**: Currently active allocated heap objects tracked by Go's runtime allocator (`m.Alloc`).
- **System Memory (`go_native_sys_mb`)**: Total virtual memory reserved from the OS kernel for heap, stack, and runtime structures (`m.Sys`).
- **Resident Set Size (`js_reference_rss_mb`)**: Peak physical RAM pages allocated to the process (`process.memoryUsage().rss`).

### C. Throughput (`throughput_ops_per_sec`)
- Evaluates maximum sustained operations per second ($\text{ops/sec}$) for each function (`Parse`, `Format`, `Bytes`) over $N = 500,000$ iterations:
  $$\text{Throughput} = \frac{N}{\Delta t_{\text{seconds}}}$$
- **Combined Total Throughput**:
  $$\text{Throughput}_{\text{combined}} = \frac{3 \times N}{\Delta t_{\text{parse}} + \Delta t_{\text{format}} + \Delta t_{\text{bytes}}}$$

### D. 99th Percentile Tail Latency (`p99_latency_us`)
- Individual operation duration $L_i$ logged in microsecond ($\mu\text{s}$) resolution using nanosecond monotonic timestamps (`time.Since(t0)`).
- Latency array $L = [L_1, L_2, \dots, L_N]$ is sorted in ascending order.
- $P_{99}$ tail latency is computed at the 99th percentile index:
  $$P_{99} = L_{\lfloor 0.99 \times N \rfloor}$$

---

## 4. Statistical Rigor & Warmup Methodology

To ensure reproducible, low-variance measurements:
1. **Pre-test Warmup Phase**: Executed 10,000 preliminary operations per function prior to data collection to stabilize:
   - Go runtime thread scheduling and stack growth.
   - CPU frequency scaling / power state transitions.
   - Garbage collector heap pacing.
2. **Sample Size**: $500,000$ iterations per function ($1.5 \times 10^6$ total operations per test run).
3. **Allocation Profiling**: Go benchmarks run with zero unexpected heap escapes enabled by value semantic data structures.

---

## 5. Execution Instructions

### Run Standard Go Microbenchmarks
```bash
go test -bench=. ./bench
```

### Run Full Benchmarking Suite & Export `results.json`
```bash
go run ./bench/main.go
```
*(Updates [`bench/results.json`](https://github.com/DebadityaDas05/go-bytes/blob/main/bench/results.json) with quantitative results).*
