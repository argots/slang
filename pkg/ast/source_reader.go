package ast

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// SourceReader implements the abiltiy to read source code.
type SourceReader interface {
	ReadSource(location string) io.ReadCloser
}

// Sources implements a source reaader for reading for text, file or url.
type Sources struct {
	// Client allows overriding the default client if needed.
	*http.Client
	readers map[string]func() io.ReadCloser
	cache   map[string][]byte
}

func (s *Sources) AddStringSource(location, source string) {
	s.init()
	s.readers[location] = func() io.ReadCloser {
		return ioutil.NopCloser(strings.NewReader(source))
	}
}

func (s *Sources) AddFileSource(location, filePath string) {
	s.init()
	// should this be cached?
	s.readers[location] = func() io.ReadCloser {
		f, err := os.Open(filePath)
		if err != nil {
			return errReader{err}
		}
		return f
	}
}

func (s *Sources) AddPublicUrlSource(location, url string) {
	s.init()
	s.readers[location] = s.cacheGet(location, func() ([]byte, error) {
		resp, err := s.client().Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	})
}

func (s *Sources) init() {
	if s.readers == nil {
		s.readers = map[string]func() io.ReadCloser{}
		s.cache = map[string][]byte{}
	}
}

func (s *Sources) client() *http.Client {
	if s.Client != nil {
		return s.Client
	}
	return httpClient
}

func (s *Sources) cacheGet(location string, fn func() ([]byte, error)) func() io.ReadCloser {
	return func() io.ReadCloser {
		// TODO: implement cache limtis, LRU eviction etc
		cached, ok := s.cache[location]
		if !ok {
			cached, err := fn()
			if err != nil {
				return errReader{err}
			}
			s.cache[location] = cached
		}
		return ioutil.NopCloser(bytes.NewReader(cached))
	}
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 2 * time.Second,
	},
}
