<h1 align="center">
  <br>
    <img width="600" src="./assets/logos/logotype-horizontal.png">
  <br>
</h1>

**headertrace** is a simple HTTP server echoing back HTTP client Headers in the response body.

It allows also to specify custom headers to be put in the HTTP reply.

It's main usage is to help debugging HTTP header transformations.

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