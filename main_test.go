package main

import (
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestFormatBytes_ParsesAllLines(t *testing.T) {
	// NOTE: This mimics what readConnection returns: bytes *before* \r\n\r\n.
	// That means it typically does NOT end with '\n' after the last header line.
	raw := []byte(
		"GET / HTTP/1.1\r\n" +
			"Host: localhost:8080\r\n" +
			"User-Agent: curl/8.0.0",
	)

	lines := formatBytes(raw)

	// What we WANT:
	want := []string{
		"GET / HTTP/1.1",
		"Host: localhost:8080",
		"User-Agent: curl/8.0.0",
	}

	if len(lines) != len(want) {
		t.Fatalf("expected %d lines, got %d: %#v", len(want), len(lines), lines)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line[%d]: expected %q, got %q", i, want[i], lines[i])
		}
	}
}
// idk what t *testing.T mean, will learn it on the 2nd
func TestHeaderFieldParser(t *testing.T){

	in:=[]string{
"Host: localhost:42069",
"User-Agent: curl/7.81.0",
"Accept: *//*",
"Content-Length: 21"}

headermap,bodyflag:=headerfieldParser(in)
want:=map[string]string{
	"Host":"localhost:42069",
"User-Agent": "curl/7.81.0",
"Accept": "*//*",
"Content-Length": "21"}
	if len(headermap) != len(want) {
		t.Fatalf("expected %d header fields, got %d: %#v", len(want), len(headermap), headermap)
	}

	for k, v := range want {
		if got, ok := headermap[k]; !ok {
			t.Fatalf("expected key %q missing in header map", k)
		} else if got != v {
			t.Fatalf("for key %q, expected value %q, got %q", k, v, got)
		}
	}

	if !bodyflag {
		t.Fatalf("expected bodyflag = true, got false")
	}
}



func TestHeaderParser_OK(t *testing.T) {
	in := []string{
		"GET /hello HTTP/1.1",
		"Host: localhost:8080",
	}

	sl, err := headerParser(in)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if sl.method != "GET" {
		t.Fatalf("method: expected %q, got %q", "GET", sl.method)
	}
	if sl.path != "/hello" {
		t.Fatalf("path: expected %q, got %q", "/hello", sl.path)
	}
	if sl.version != "HTTP/1.1" {
		t.Fatalf("version: expected %q, got %q", "HTTP/1.1", sl.version)
	}
}

func TestHeaderParser_InvalidMethodLowercase(t *testing.T) {
	in := []string{
		"Get / HTTP/1.1", // has lowercase letters
		"Host: localhost:8080",
	}
	_, err := headerParser(in)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "incorrect http method") {
		t.Fatalf("expected method error, got %v", err)
	}
}

func TestHeaderParser_InvalidVersion(t *testing.T) {
	in := []string{
		"GET / HTTP/2.0",
		"Host: localhost:8080",
	}
	_, err := headerParser(in)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "wrong http version") {
		t.Fatalf("expected version error, got %v", err)
	}
}

func TestReadConnection_ReturnsHeaderOnly(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Write a request that includes headers + delimiter + body.
	req := "POST /coffee HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"Content-Length: 5\r\n" +
		"\r\n" +
		"HELLO"

	// Run readConnection on the "server" end.
	done := make(chan []byte, 1)
	go func() {
		h := readConnection(server)
		done <- h
	}()

	// Write from the client end in one shot (or chunks, both should work).
	_, err := client.Write([]byte(req))
	if err != nil {
		t.Fatalf("client write failed: %v", err)
	}

	select {
	case header := <-done:
		got := string(header)
		if strings.Contains(got, "\r\n\r\n") {
			t.Fatalf("header should not contain delimiter, got: %q", got)
		}
		if strings.Contains(got, "HELLO") {
			t.Fatalf("header should not include body, got: %q", got)
		}
		if !strings.Contains(got, "POST /coffee HTTP/1.1") {
			t.Fatalf("missing request line, got: %q", got)
		}
		if !strings.Contains(got, "Content-Length: 5") {
			t.Fatalf("missing header field, got: %q", got)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout waiting for readConnection to return")
	}
}

func TestSendResponse_WritesExpected(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Call sendResponse on server side
	go sendResponse(server)

	// Read all data from client side (it will return once the other side closes,
	// but sendResponse doesn't close; so we read a fixed amount with ReadFull-ish logic.)
	buf := make([]byte, 1024)
	n, err := client.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("client read failed: %v", err)
	}
	got := string(buf[:n])

	wantPrefix := "HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\nHello\r\n"
	if got != wantPrefix {
		t.Fatalf("unexpected response:\n--- got ---\n%q\n--- want ---\n%q", got, wantPrefix)
	}
}
