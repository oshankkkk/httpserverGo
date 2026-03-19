# Go TCP HTTP Server

A from-scratch, foundational HTTP/1.1 server implementation built directly on top of raw TCP sockets in Go.

This bypasses Go's standard `net/http` package to manually handle connections, parse byte streams, and construct HTTP responses. 

This codebase serves as a simple sandbox for understanding how HTTP works under the hood over raw TCP connections. Currently, the server scope is limited to receiving, parsing, logging, and sending a static HTTP response.

## Getting Started

To run the server locally:

```bash
go run main.go
```
The server will begin listening on `localhost:8080`. You can test this through curl or postman

To run the test suite:

```bash
go test
```
