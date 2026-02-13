<h1 align="center">
  <br>
    <img width="600" src="./assets/logos/logotype-horizontal.png">
  <br>

<p align="center">
  <a href="https://github.com/fgiudici/headertrace/releases/latest"><img alt="Latest Release" src="https://img.shields.io/github/v/release/fgiudici/headertrace"></a>
  <a href="https://quay.io/repository/fgiudici/headertrace"><img alt="Container Image" src="https://img.shields.io/badge/container-quay.io%2Ffgiudici%2Fheadertrace-blue"></a>
  <a href="https://github.com/fgiudici/headertrace/actions/workflows/go.yml"><img alt="Go Build" src="https://github.com/fgiudici/headertrace/actions/workflows/go.yml/badge.svg"></a>
  <a href="https://github.com/fgiudici/headertrace/blob/main/go.mod"><img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/fgiudici/headertrace"></a>
  <a href="https://goreportcard.com/report/github.com/fgiudici/headertrace"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/fgiudici/headertrace"></a>
  <a href="https://github.com/fgiudici/headertrace/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/github/license/fgiudici/headertrace"></a>
</p>

</h1>

**headertrace** is a simple HTTP server echoing back HTTP client Headers in the response body.

It allows also to specify custom headers to be put in the HTTP reply.

It's main usage is to help debugging HTTP header transformations.

## Installation

**headertrace** is built for different OS/architectures as a standalone binary (see the [Download Binaries](#download-binaries) section).

Linux container images for both AMD64 and ARM64 architectures are also available 
Jump to the [Container Images](#container-images) section for more info.

### Download Binaries

1. Visit the [releases page](https://github.com/fgiudici/headertrace/releases) and download the appropriate binary for your operating system and architecture:
   - **Linux AMD64**:
   [`headertrace-linux-amd64`](https://github.com/fgiudici/headertrace/releases/latest/download/headertrace-linux-amd64)
   - **Linux ARM64**:
   [`headertrace-linux-arm64`](https://github.com/fgiudici/headertrace/releases/latest/download/headertrace-linux-arm64)
   - **macOS AMD64**:
   [`headertrace-darwin-amd64`](https://github.com/fgiudici/headertrace/releases/latest/download/headertrace-darwin-amd64)
   - **macOS ARM64** (Apple Silicon):
   [`headertrace-darwin-arm64`](https://github.com/fgiudici/headertrace/releases/latest/download/headertrace-darwin-arm64)
   - **Windows AMD64**:
   [`headertrace-windows-amd64.exe`](https://github.com/fgiudici/headertrace/releases/latest/download/headertrace-windows-amd64.exe)

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

Running the binary with no arguments starts the HTTP server on port 8080):

```bash
headertrace
```

connect to the server with a browser or an HTTP client to get headers echoed back:

```bash
$ curl 127.0.0.1:8080 
{
  "headers": {
    "Accept": "*/*",
    "User-Agent": "curl/8.5.0"
  },
  "host": "localhost:8080",
  "method": "GET",
  "path": "/",
  "protocol": "HTTP/1.1"
}

```
For a list of available options, see the [Usage section](#usage) or view the inline help:

```bash
headertrace -h
```

### Container Images
Linux container images are available at [Quay.io](https://quay.io/fgiudici/headertrace).

You can run the container images directly with **podman** or **docker**:
```bash
docker run quay.io/fgiudici/headertrace:latest
```
and pass the arguments as usual (see the [Usage section](#usage) or view the inline help `-h`).


A sample Kubernetes Deployment and an associated Service would look like:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: headertrace
  name: headertrace
  namespace: headertrace
spec:
  replicas: 2
  selector:
    matchLabels:
      app: headertrace
  strategy: {}
  template:
    metadata:
      labels:
        app: headertrace
    spec:
      containers:
      - image: quay.io/fgiudici/headertrace:latest
        args: ["-H X-Served-By:Headertrace"]
        name: headertrace
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: headertrace
  name: headertrace
  namespace: headertrace
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: headertrace
  type: NodePort
```

## Usage

### Command Line Reference

```
headertrace [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--address` | `-a` | `0.0.0.0` | IP address (or domain) to bind the server to. |
| `--port` | `-p` | `8080` | TCP port to bind the server to. |
| `--header` | `-H` | _(none)_ | Custom HTTP headers to add to every response (format: `key1:value1,key2:value2`). |
| `--drop-header` | `-D` | _(none)_ | HTTP headers to redact from request headers echoed in the response body (format: `key1,key2`). |
| `--privacy` | `-P` | `false` | Drop `X-Forwarded-*`, `X-Real-IP`, and `Cf-*` (Cloudflare) headers from echoed request headers. |
| `--sent` | `-s` | `false` | Include the HTTP headers added in the server response inside the response body. |
| `--log-level` | `-l` | _(none)_ | Set the logging verbosity. Accepted values: `TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`. Overrides the `LOG_LEVEL` environment variable. |
| `--version` | `-v` | | Print version and exit. |
| `--help` | `-h` | | Print help and exit. |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `LOG_LEVEL` | Sets the logging verbosity (`TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`). Can be overridden by the `--log-level` flag. Defaults to `INFO` if unset. |

### Response Format

**headertrace** responds with a JSON body containing information about the received HTTP request:

```json
{
  "headers": {
    "Accept": "*/*",
    "User-Agent": "curl/8.5.0"
  },
  "host": "localhost:8080",
  "method": "GET",
  "path": "/",
  "protocol": "HTTP/1.1"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `headers` | object | HTTP headers received in the client request. |
| `host` | string | Host (and port) the request was sent to. |
| `method` | string | HTTP method of the request (e.g. `GET`). |
| `path` | string | Request URI path. |
| `protocol` | string | HTTP protocol version (e.g. `HTTP/1.1`). |
| `sent` | object | _(Optional)_ HTTP headers added in the server response. Only present when `-s` / `--sent` is enabled. |

### Examples

#### Basic server

Start **headertrace** with default settings (listening on `0.0.0.0:8080`):

```bash
headertrace
```

Test it with curl:

```bash
$ curl -s http://localhost:8080 | jq .
{
  "headers": {
    "Accept": "*/*",
    "User-Agent": "curl/8.5.0"
  },
  "host": "localhost:8080",
  "method": "GET",
  "path": "/",
  "protocol": "HTTP/1.1"
}
```

#### Bind to a specific address and port

```bash
headertrace -a 192.168.1.10 -p 3000
```

#### Add custom response headers

Inject custom headers into every HTTP response. Useful for simulating upstream services that set specific headers:

```bash
headertrace -H "X-Served-By:headertrace,X-Request-Region:eu-west-1"
```

```bash
$ curl -s -D - http://localhost:8080
HTTP/1.1 200 OK
Content-Type: application/json
X-Request-Region: eu-west-1
X-Served-By: headertrace

{
  "headers": { ... },
  ...
}
```

#### Inspect response headers in the body (`--sent`)

Include the headers sent by the server in the JSON response body — useful for verifying what headers the server is actually returning:

```bash
headertrace -s -H "X-Custom:hello"
```

```bash
$ curl -s http://localhost:8080 | jq .
{
  "headers": {
    "Accept": "*/*",
    "User-Agent": "curl/8.5.0"
  },
  "host": "localhost:8080",
  "method": "GET",
  "path": "/",
  "protocol": "HTTP/1.1",
  "sent": {
    "Content-Type": "application/json",
    "X-Custom": "hello"
  }
}
```

#### Redact specific request headers (`--drop-header`)

Hide specific request headers from the response body. This is useful when you want to filter out noisy or irrelevant headers:

```bash
headertrace -D "Authorization,Cookie"
```

```bash
$ curl -s -H "Authorization: Bearer secret" -H "Cookie: session=abc" http://localhost:8080 | jq .
{
  "headers": {
    "Accept": "*/*",
    "User-Agent": "curl/8.5.0"
  },
  "host": "localhost:8080",
  "method": "GET",
  "path": "/",
  "protocol": "HTTP/1.1"
}
```

The `Authorization` and `Cookie` headers are received by the server but omitted from the response body.

#### Privacy mode

Enable privacy mode to automatically redact proxy-related headers (`X-Forwarded-*`, `X-Real-IP`) and Cloudflare headers (`Cf-*`):

```bash
headertrace -P
```

This is particularly useful when the server sits behind a reverse proxy or CDN and you want to avoid echoing back internal network information.

#### Verbose logging

Increase log verbosity for troubleshooting. At `DEBUG` level, redacted headers are logged; at `TRACE` level, all header values are logged:

```bash
headertrace -l DEBUG
```

Or using the environment variable:

```bash
LOG_LEVEL=TRACE headertrace
```

#### Combining multiple options

Run on a custom port, add response headers, redact sensitive request headers, enable privacy mode, and dump sent headers — all at once:

```bash
headertrace -p 9090 -H "X-Served-By:headertrace" -D "Authorization" -P -s -l DEBUG
```

#### Running with containers

```bash
docker run -p 8080:8080 quay.io/fgiudici/headertrace:latest -H "X-Served-By:headertrace" -P -s
```
