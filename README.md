<h1 align="center">
  <br>
    <img width="600" src="./assets/logos/logotype-horizontal.png">
  <br>
</h1>

**headertrace** is a simple HTTP server echoing back HTTP client Headers in the response body.

It allows also to specify custom headers to be put in the HTTP reply.

It's main usage is to help debugging HTTP header transformations.

## Installation

**headertrace** is built for different OS/architectures as a standalone binary (see the #download-pre-built-binaries section).

Linux container images for both AMD64 and ARM64 architectures are also available 
Jump to the #container-images section for more info.

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

Start the HTTP server (default TCP port is 8080):

```bash
headertrace
```

to listen on a particular IP address (or hostname) and TCP port use the `-a` and `-p` parameters:

```bash
headertrace -a 192.168.1.10 -p 1234
```

to add custom Headers in all HTTP server replies use the `-H` flag:

```bash
headertrace -H "X-Custom-HDR1:value,X-Custom-HDR2:value"
```

inline help:

```bash
headertrace -h
```

### Container Images
Linux container images are available at [Quay.io](https://quay.io/fgiudici/headertrace).

You can run the container images directly with **podman** or **docker**:
```bash
docker run quay.io/fgiudici/headertrace:latest
```
and pass the options as usual (see the [Quickstart](#quickstart) section).


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