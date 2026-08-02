# bytes-go

A high-performance, cross-language byte utility library written in **Go** and exposed to **Node.js** via a CGo bridge and [Koffi](https://koffi.dev/) FFI. 

Ported directly from TJ Holowaychuk's original [`bytes`](https://github.com/bytesjs/bytes.js) module with 100% feature parity.

---

## 🛠️ Prerequisites

Ensure you have the following installed on your system:

- **Node.js**: `v18.0.0` or higher
- **Go**: `v1.20` or higher
- **C Compiler / CGo Toolchain**:
  - **Windows**: [MinGW-w64](https://www.winlibs.com/) or GCC via MSYS2 / WinLibs (required for CGo compilation)
  - **macOS**: Xcode Command Line Tools (`xcode-select --install`)
  - **Linux**: `build-essential` (`sudo apt install build-essential`)

---

## 🚀 Quick Setup & Build

### 1. Clone the Repository
```bash
git clone https://github.com/DebadityaDas05/go-bytes.git
cd go-bytes
```

### 2. Install Node.js Dependencies
```bash
npm install
```

### 3. Build Native Shared Library

Compile the Go CGo bridge into a shared library binary for your platform:

- **Windows**:
  ```powershell
  go build -buildmode=c-shared -o libbytes.dll ./bridge
  ```

- **Linux**:
  ```bash
  go build -buildmode=c-shared -o libbytes.so ./bridge
  ```

- **macOS**:
  ```bash
  go build -buildmode=c-shared -o libbytes.dylib ./bridge
  ```

- **Npm Build Script (Cross-Platform Helper)**:
  ```bash
  npm run build
  ```

---

## 🧪 Running Tests & Fuzzing

### 1. Node.js Mocha Test Suite (Original JS Compatibility)
Execute the original Mocha test suite to verify the Koffi FFI dynamic library bridge:
```bash
npm test
```
*(Or `npx mocha test-original`)*

### 2. Go Native Unit Tests
Run the native Go unit tests:
```bash
go test ./test-port
```

### 3. Differential Fuzzing Suite
Run the 60s+ differential fuzzer comparing JS reference output vs Go implementation across millions of randomized inputs:
```bash
npm run fuzz
```

---

## ⚡ Performance & Benchmarks

Run the benchmark suite to evaluate throughput ($\text{ops/sec}$), tail latency ($P_{99}$), startup time, and memory footprint:

```bash
# Run Go standalone benchmark runner (exports bench/results.json)
go run ./bench/main.go

# Run Go standard testing.B microbenchmarks
go test -bench=. ./bench
```

Detailed benchmarking methodologies and empirical results can be found in:
- [`bench/methodology.md`](https://github.com/DebadityaDas05/go-bytes/blob/main/bench/methodology.md)
- [`bench/results.json`](https://github.com/DebadityaDas05/go-bytes/blob/main/bench/results.json)

---

## 📖 Usage Examples

### JavaScript / Node.js
```javascript
const bytes = require('./index.js');

// Convert string to bytes integer
bytes('1KB');                     // 1024
bytes.parse('1.5MB');             // 1572864

// Convert bytes integer to formatted string
bytes(1024);                      // '1KB'
bytes.format(1000, { thousandsSeparator: ' ' }); // '1 000B'
bytes.format(1024, { unitSeparator: ' ' });      // '1 KB'
bytes.format(1024, { decimalPlaces: 3, fixedDecimals: true }); // '1.000KB'
```

### Native Go
```go
package main

import (
    "fmt"
    bytesutil "github.com/DebadityaDas05/go-bytes/src"
)

func main() {
    fmt.Println(bytesutil.Parse("1.5MB")) // Outputs: 1572864
    fmt.Println(bytesutil.Format(1024, nil)) // Outputs: "1KB"
}
```

---

## 🏗️ Project Architecture

```text
bytes-go/
├── bench/
│   ├── benchmark.js        # Node.js performance benchmark runner
│   ├── bytes_bench_test.go # Go standard testing.B microbenchmarks
│   ├── main.go             # Go benchmark runner exporting results.json
│   ├── methodology.md      # Detailed benchmark methodology & formulas
│   └── results.json        # Quantitative p99, RSS, startup & throughput metrics
├── bridge/
│   └── bridge.go           # CGo export functions bridging C types to Go
├── fuzz/
│   ├── harness.js          # Differential fuzzer engine comparing JS vs Go
│   ├── index.js            # Original reference JS library
│   └── log.txt             # 60s+ execution log demonstrating 0 divergences
├── src/
│   └── bytes.go            # Pure Go implementation of bytes formatting and parsing
├── test-original/
│   ├── bytes.js            # Core constructor unit tests (Mocha)
│   ├── byte-format.js      # Format function test suite (Mocha)
│   └── byte-parse.js       # Parse function test suite (Mocha)
├── test-port/
│   └── bytes_test.go       # Go native unit test suite
├── index.js                # Node.js FFI wrapper powered by Koffi
├── go.mod                  # Go module definition
├── package.json            # Node.js package configuration
└── README.md               # Project documentation
```

---

## 📄 License

MIT
