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

## 🧪 Running Tests

Execute the Mocha test suite to verify the native library bridge and formatting logic:

```bash
npm test
```

---

## 📖 Usage Examples

```javascript
const bytes = require('./adapter');

// Convert string to bytes integer
bytes('1KB');                     // 1024
bytes.parse('1.5MB');             // 1572864

// Convert bytes integer to formatted string
bytes(1024);                      // '1KB'
bytes.format(1000, { thousandsSeparator: ' ' }); // '1 000B'
bytes.format(1024, { unitSeparator: ' ' });      // '1 KB'
bytes.format(1024, { decimalPlaces: 3, fixedDecimals: true }); // '1.000KB'
```

---

## 🏗️ Project Architecture

```text
bytes-go/
├── adapter/
│   └── index.js       # Node.js FFI wrapper powered by Koffi
├── bridge/
│   └── bridge.go      # CGo export functions bridging C types to Go
├── bytesutil/
│   └── bytes.go       # Pure Go implementation of bytes formatting and parsing
├── test/
│   ├── bytes.js       # Core constructor unit tests
│   ├── byte-format.js # Format function test suite
│   └── byte-parse.js  # Parse function test suite
├── go.mod             # Go module definition
├── package.json       # Node.js package configuration
└── README.md          # Project documentation
```

---

## 📄 License

MIT
