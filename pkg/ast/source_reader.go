package ast

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// SourceReader implements an Reader for source code.
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

// ReadSource returns an io.Reader for the source contents.
func (s *Sources) ReadSource(location string) io.ReadCloser {
	if fn, ok := s.readers[location]; ok {
		return fn()
	}
	return nil
}

// AddStringSource adds the string contents as a source.
func (s *Sources) AddStringSource(location, source string) {
	s.init()
	s.readers[location] = func() io.ReadCloser {
		return ioutil.NopCloser(strings.NewReader(source))
	}
}

// AddFileSource adds the file contents as a source.
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

// AddPublicURLSource adds the URL contents as a source.
//
// If the Client field of Sources is non-nil, that is used. If not, a
// custom client is used with timeouts filled in.
func (s *Sources) AddPublicURLSource(location, url string) {
	s.init()
	s.readers[location] = s.cacheGet(location, func() ([]byte, error) {
		resp, err := s.client().Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("http.Get failed: %v", resp.StatusCode)
		}

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
		cached, ok := s.cache[location]
		var err error
		if !ok {
			cached, err = fn()
			if err != nil {
				return errReader{err}
			}
			s.cache[location] = cached
		}
		return ioutil.NopCloser(bytes.NewReader(cached))
	}
}

//nolint: gochecknoglobals, gomnd
var httpClient = &http.Client{
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Minute / 2,
			KeepAlive: time.Minute / 2,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 2 * time.Second,
	},
}
