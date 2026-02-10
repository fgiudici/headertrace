<h1 align="center">
  <br>
    <img width="600" src="./assets/logos/logotype-horizontal.png">
  <br>
</h1>

**headertrace** is a simple HTTP server echoing back HTTP client Headers in the response body.

It allows also to specify custom headers to be put in the HTTP reply.

It's main usage is to help debugging HTTP header transformations.

## Installation

### Download Pre-built Binary

1. Visit the [releases page](https://github.com/fgiudici/headertrace/releases) and download the appropriate binary for your operating system and architecture:
   - **Linux AMD64**: `headertrace-linux-amd64`
   - **Linux ARM64**: `headertrace-linux-arm64`
   - **macOS AMD64**: `headertrace-darwin-amd64`
   - **macOS ARM64** (Apple Silicon): `headertrace-darwin-arm64`
   - **Windows AMD64**: `headertrace-windows-amd64.exe`

2. Rename the downloaded binary to `headertrace` (or `headertrace.exe` on Windows):
   ```bash
   mv headertrace-linux-amd64 headertrace
   ```

3. Make the binary executable (Linux/macOS only):
   ```bash
   chmod +x headertrace
   ```

4. Optionally, move it to a directory in your PATH for easier access:
   ```bash
   sudo mv headertrace /usr/local/bin/
   ```

## Quickstart

Start the HTTP server (default TCP port is 8080):

```bash
headertrace
```

to listen on a particular IP address (or hostname) and TCP port use the `--host` and `--port` parameters:

```bash
headertrace --host 192.168.1.10 --port 1234
```

to add custom Headers in all HTTP server replies use the `-H` flag:

```bash
headertrace -H "X-Custom-HDR1:value,X-Custom-HDR2:value"
```

inline help:

```bash
headertrace -h
```