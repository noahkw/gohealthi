# gohealthi

A small gRPC server in Go to fetch a machine's health stats.

## Building for other architectures

Building for linux-arm64 on Windows:

```PS
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o gohealthi-arm64 ./cmd/gohealthi  
```

This creates a `gohealthi-arm64` binary in the current directory.

## Hacking

Regenerate proto sources:

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/health.proto
```
