package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestForwardHeadersToRequest(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp})

	headerKey := "Dnt"
	headerVal := "1"
	customHeaderKey := "X-Hello"
	customHeaderVal := "World"

	// setting up a server so that we can make a proxy request to it
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customRecHdr, ok := r.Header[customHeaderKey]
		if !ok || customRecHdr[0] != customHeaderVal {
			t.Errorf("did not forward custom header in proxy request: %v, got: %v", customHeaderKey, customRecHdr)
		}
		recHdr, ok := r.Header[headerKey]
		if !ok || recHdr[0] != headerVal {
			t.Errorf("did not forward header in proxy request: %v, got: %v", headerKey, recHdr)
		}
	}))
	defer ts.Close()

	req, _ := http.NewRequest("GET", "/"+ts.URL, nil)
	req.Header.Add(headerKey, headerVal)
	req.Header.Add(customHeaderKey, customHeaderVal)
	rr := httptest.NewRecorder()

	h := handler{}
	h.ServeHTTP(rr, req)
}

func TestForwardHeadersFromResponse(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp})

	headerKey := "Dnt"
	headerVal := "1"
	customHeaderKey := "X-Hello"
	customHeaderVal := "World"

	// setting up a server so that we can make a proxy request to it
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(headerKey, headerVal)
		w.Header().Add(customHeaderKey, customHeaderVal)
	}))
	defer ts.Close()

	req, _ := http.NewRequest("GET", "/"+ts.URL, nil)
	req.Header.Add(headerKey, headerVal)
	req.Header.Add(customHeaderKey, customHeaderVal)
	rr := httptest.NewRecorder()

	h := handler{}
	h.ServeHTTP(rr, req)

	resp := rr.Result()
	customRecHdr, ok := resp.Header[customHeaderKey]
	if !ok || customRecHdr[0] != customHeaderVal {
		t.Errorf("did not forward custom header to response: %v, got: %v", customHeaderKey, customRecHdr)
	}
	recHdr, ok := resp.Header[headerKey]
	if !ok || recHdr[0] != headerVal {
		t.Errorf("did not forward header to response: %v, got: %v", headerKey, recHdr)
	}
}

func TestForwardQueryParameters(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp})

	queryKey := "hello"
	queryVal := "world"

	// setting up a server so that we can make a proxy request to it
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.URL.Query().Get(queryKey)
		if got != queryVal {
			t.Errorf("did not forward query params to request: %v, got: %v", queryKey, got)
		}
	}))
	defer ts.Close()

	req, _ := http.NewRequest("GET", "/"+ts.URL, nil)
	q := make(url.Values, 1)
	q.Add(queryKey, queryVal)
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	h := handler{}
	h.ServeHTTP(rr, req)
}

func TestForwardBodyToRequest(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp})

	body := "hello world"
	reqBody := strings.NewReader(body)

	// setting up a server so that we can make a proxy request to it
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("got error while reading body: %v", err)
		}
		if string(got) != body {
			t.Errorf("expected body: %v, got: %v", body, string(got))
		}
	}))
	defer ts.Close()

	req, _ := http.NewRequest("GET", "/"+ts.URL, reqBody)

	rr := httptest.NewRecorder()

	h := handler{}
	h.ServeHTTP(rr, req)
}

func TestForwardBodyFromResponse(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp})

	body := "hello world"
	reqBody := strings.NewReader(body)

	// setting up a server so that we can make a proxy request to it
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	}))
	defer ts.Close()

	req, _ := http.NewRequest("GET", "/"+ts.URL, reqBody)

	rr := httptest.NewRecorder()

	h := handler{}
	h.ServeHTTP(rr, req)

	resp := rr.Result()
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("got error while reading body: %v", err)
	}
	if string(got) != body {
		t.Errorf("expected body: %v, got: %v", body, string(got))
	}
}
