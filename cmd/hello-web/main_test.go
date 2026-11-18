package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

const browserAccept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"

func TestHandlerHello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}
	if got, want := res.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Fatalf("Content-Type = %q, want %q", got, want)
	}
	if got, want := rec.Body.String(), "hello"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestHandlerBrowserGetsHTML(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", browserAccept)
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}
	if got, want := res.Header.Get("Content-Type"), "text/html; charset=utf-8"; got != want {
		t.Fatalf("Content-Type = %q, want %q", got, want)
	}
	body := rec.Body.String()
	for _, want := range []string{
		"<!doctype html>",
		`<html lang="en">`,
		"<h1>hello</h1>",
		"<title>hello",
		`name="viewport"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("HTML body does not contain %q", want)
		}
	}
	if got, want := res.Header.Get("Content-Length"), strconv.Itoa(len(body)); got != want {
		t.Fatalf("Content-Length = %q, want %q", got, want)
	}
	if got := res.Header.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options = %q, want %q", got, "nosniff")
	}
}

func TestHandlerHTMLNotEmittedForWildcard(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "*/*")
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, req)

	if got, want := rec.Result().Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Fatalf("Content-Type = %q, want %q", got, want)
	}
	if got := rec.Body.String(); got != "hello" {
		t.Fatalf("body = %q, want %q", got, "hello")
	}
}

func TestHandlerHead(t *testing.T) {
	for _, accept := range []string{"", browserAccept} {
		req := httptest.NewRequest(http.MethodHead, "/", nil)
		if accept != "" {
			req.Header.Set("Accept", accept)
		}
		rec := httptest.NewRecorder()
		newHandler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("accept %q: status = %d, want %d", accept, rec.Code, http.StatusOK)
		}
		if got := rec.Body.String(); got != "" {
			t.Fatalf("accept %q: HEAD body = %q, want empty", accept, got)
		}
	}
}

func TestHandlerHTMLUsesRequestHost(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", browserAccept)
	req.Host = "localhost:9000"
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, req)

	body := rec.Body.String()
	if strings.Contains(body, pageAddressPlaceholder) {
		t.Fatalf("HTML body still contains %q", pageAddressPlaceholder)
	}
	if !strings.Contains(body, "localhost:9000") {
		t.Fatalf("HTML body does not mention the request host")
	}
}

func TestPageHTMLEscapesHost(t *testing.T) {
	page := pageHTML(`"><script>alert(1)</script>`)

	if strings.Contains(page, "<script>alert(1)</script>") {
		t.Fatalf("pageHTML did not escape the host")
	}
	if !strings.Contains(page, "&lt;script&gt;") {
		t.Fatalf("pageHTML did not HTML-escape the host")
	}
}

func TestPageHTMLDefaultHost(t *testing.T) {
	if got := pageHTML(""); !strings.Contains(got, defaultHost) {
		t.Fatalf("pageHTML with empty host does not contain %q", defaultHost)
	}
}

func TestHandlerOtherPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestHandlerOtherMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	newHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if got, want := rec.Header().Get("Allow"), "GET, HEAD"; got != want {
		t.Fatalf("Allow = %q, want %q", got, want)
	}
}

func TestWantsHTML(t *testing.T) {
	tests := []struct {
		name   string
		accept string
		want   bool
	}{
		{name: "browser", accept: browserAccept, want: true},
		{name: "html only", accept: "text/html", want: true},
		{name: "with parameters", accept: "text/html;q=0.9", want: true},
		{name: "uppercase", accept: "TEXT/HTML", want: true},
		{name: "curl wildcard", accept: "*/*", want: false},
		{name: "json client", accept: "application/json", want: false},
		{name: "no header", accept: "", want: false},
		{name: "text/* only", accept: "text/*", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			if got := wantsHTML(req); got != tt.want {
				t.Fatalf("wantsHTML(%q) = %v, want %v", tt.accept, got, tt.want)
			}
		})
	}
}
