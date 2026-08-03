package proxy

import (
	"net/http"
	"testing"
)

func TestPreserveEventStream(t *testing.T) {
	response := &http.Response{
		Header: http.Header{
			"Cache-Control":  {"no-cache"},
			"Content-Length": {"128"},
			"Content-Type":   {"text/event-stream; charset=utf-8"},
		},
		ContentLength: 128,
	}

	if err := preserveEventStream(response); err != nil {
		t.Fatalf("preserve event stream: %v", err)
	}
	if got := response.Header.Get("Cache-Control"); got != "no-cache, no-transform" {
		t.Fatalf("expected no-transform cache control, got %q", got)
	}
	if got := response.Header.Get("Content-Length"); got != "" {
		t.Fatalf("expected no content length, got %q", got)
	}
	if response.ContentLength != -1 {
		t.Fatalf("expected unknown content length, got %d", response.ContentLength)
	}
}

func TestPreserveEventStreamLeavesOtherResponsesUnchanged(t *testing.T) {
	response := &http.Response{
		Header: http.Header{
			"Cache-Control":  {"public"},
			"Content-Length": {"128"},
			"Content-Type":   {"application/json"},
		},
		ContentLength: 128,
	}

	if err := preserveEventStream(response); err != nil {
		t.Fatalf("preserve event stream: %v", err)
	}
	if got := response.Header.Get("Cache-Control"); got != "public" {
		t.Fatalf("expected cache control to remain unchanged, got %q", got)
	}
	if got := response.Header.Get("Content-Length"); got != "128" {
		t.Fatalf("expected content length to remain unchanged, got %q", got)
	}
	if response.ContentLength != 128 {
		t.Fatalf("expected content length 128, got %d", response.ContentLength)
	}
}
