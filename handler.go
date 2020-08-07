package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/pprof"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

// using custom handler since the default http handler redirects requests
// with double slashes
type handler struct{}

func writeError(w http.ResponseWriter, err error, status int) error {
	type errorStruct struct {
		Error string `json:"error"`
	}
	enc := json.NewEncoder(w)
	w.WriteHeader(status)
	return enc.Encode(errorStruct{Error: err.Error()})
}

// getRequestURL returns the proxy request URL
func getRequestURL(r *http.Request) (*url.URL, error) {
	p := r.URL.String()[1:]
	if !strings.HasPrefix(p, "http") {
		p = "http://" + p
	}

	return url.ParseRequestURI(p)
}

var errRootRequest = errors.New("root request")

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// can ignore the error here
	defer func() {
		if r.Body != nil {
			r.Body.Close()
		}
	}()

	if strings.HasPrefix(r.URL.String(), "/debug/pprof") {
		pprof.Index(w, r)
		return
	}
	if r.URL.String() == "/" {
		writeError(w, errRootRequest, http.StatusBadRequest)

		log.Warn().Msg("Root request")
		return
	}

	URL, err := getRequestURL(r)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		log.Err(err).Str("url", r.URL.String()).Msg("failed to parse url")
		return
	}

	// create proxy request
	proxyReq, err := http.NewRequest(r.Method, URL.String(), r.Body)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		log.Err(err).Str("method", r.Method).Str("url", URL.String()).Msg("failed to create proxy request")
		return
	}

	// forward headers to request
	for k, v := range r.Header {
		proxyReq.Header.Add(k, strings.Join(v, " "))
	}

	log.Trace().Str("url", URL.String()).Msg("making a proxy request")
	proxyResp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		log.Err(err).Str("method", r.Method).Str("url", URL.String()).Msg("failed to send proxy request")
		return
	}

	// forward headers to response
	for k, v := range proxyResp.Header {
		w.Header().Add(k, strings.Join(v, " "))
	}

	_, err = io.Copy(w, proxyResp.Body)
	// note: calling Close twice is okay in this case
	defer proxyResp.Body.Close()
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		log.Err(err).Str("method", r.Method).Str("url", r.URL.String()).Msg("failed to copy response body")
		return
	}

	err = proxyResp.Body.Close()
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		log.Err(err).Str("method", r.Method).Str("url", r.URL.String()).Msg("failed to close response body")
		return
	}

	log.Info().Str("method", r.Method).Str("url", URL.String()).Msg("succ")
}
