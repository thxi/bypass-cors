package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

// using custom handler since the default http handler redirects requests
// with double slashes
type handler struct{}

func writeError(w io.Writer, err error) error {
	type errorStruct struct {
		Error string `json:"error"`
	}
	enc := json.NewEncoder(w)
	return enc.Encode(errorStruct{Error: err.Error()})
}

// getRequestURL returns the requested URL to bypass-cors
func getRequestURL(r *http.Request) (*url.URL, error) {
	p := r.URL.String()[1:]
	if !strings.HasPrefix(p, "http") {
		p = "http://" + p
	}

	return url.ParseRequestURI(p)
}

var errRootRequest = errors.New("root request")

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.String(), "/debug/pprof") {
		pprof.Index(w, r)
		return
	}
	if r.URL.String() == "/" {
		// TODO add error
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, errRootRequest)

		log.Warn().Msg("Root request")
		return
	}

	URL, err := getRequestURL(r)
	if err != nil {
		writeError(w, err)
		log.Err(err).Str("url", r.URL.String()).Msg("failed to parse url")
		return
	}

	// create proxy request
	proxyReq, err := http.NewRequest(r.Method, URL.String(), r.Body)
	if err != nil {
		writeError(w, err)
		log.Err(err).Str("method", r.Method).Str("url", URL.String()).Msg("failed to create proxy request")
		return
	}

	// forward headers
	for k, v := range r.Header {
		proxyReq.Header.Add(k, strings.Join(v, " "))
	}

	log.Trace().Str("url", URL.String()).Msg("making a proxy request")
	// TODO: what about following redirects?
	res, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		writeError(w, err)
		log.Err(err).Str("method", r.Method).Str("url", URL.String()).Msg("failed to send proxy request")
		return
	}

	// TODO forward response headers

	// TODO stream body?
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		writeError(w, err)
		log.Err(err).Str("method", r.Method).Str("url", r.URL.String()).Msg("failed to read response body")
		return
	}

	// ignoring error
	err = res.Body.Close()
	if err != nil {
		log.Err(err).Str("method", r.Method).Str("url", r.URL.String()).Msg("failed to close response body")
		return
	}

	_, err = w.Write(body)
	if err != nil {
		log.Err(err).Str("method", r.Method).Str("url", URL.String()).Msg("failed to write response")
	}

	log.Info().Str("method", r.Method).Str("url", r.URL.String()).Msg("succ")
}