package ast_test

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/argots/slang/pkg/ast"
)

func TestSourcesMissing(t *testing.T) {
	s := ast.Sources{}
	if v := s.ReadSource("something"); v != nil {
		t.Fatal("Unexpected source", v)
	}
}

func TestSourcesString(t *testing.T) {
	s := ast.Sources{}
	tests := map[string]string{
		"hello": "world",
		"world": "hello",
	}
	for k, v := range tests {
		s.AddStringSource(k, v)
	}
	for k, v := range tests {
		w, err := readAll(s.ReadSource(k))
		if string(w) != v || err != nil {
			t.Fatal("Failed to read", k, string(w), err)
		}
	}
}

func TestSourcesFile(t *testing.T) {
	s := ast.Sources{}
	tests := map[string]string{
		"hello": "world\n",
		"world": "hello\n",
	}
	for k := range tests {
		s.AddFileSource(k, "testdata/"+k+".txt")
	}
	for k, v := range tests {
		w, err := readAll(s.ReadSource(k))
		if string(w) != v || err != nil {
			t.Fatal("Failed to read", k, string(w), err)
		}
	}
}

func TestSourcesMissingFile(t *testing.T) {
	s := ast.Sources{}
	s.AddFileSource("missing", "testdata/missing.txt")
	if v, err := readAll(s.ReadSource("missing")); len(v) != 0 || err == nil {
		t.Fatal("Unexpected source", v, err)
	}
}

func TestSourcesUrl(t *testing.T) {
	s := ast.Sources{}
	s.AddPublicUrlSource("robots.txt", "https://www.google.com/robots.txt")
	data, err := readAll(s.ReadSource("robots.txt"))
	if err != nil || !strings.Contains(string(data), "Disallow") {
		t.Fatal("Unexpected source", string(data), err)
	}
	data2, err := readAll(s.ReadSource("robots.txt"))
	if err != nil || string(data) != string(data2) {
		t.Fatal("Unexpected source", string(data2), err)
	}
}

func TestSourcesInvalidUrl(t *testing.T) {
	s := ast.Sources{}
	s.AddPublicUrlSource("invalid_url1", "https://wwwxx.google.com/invalid_url")
	data, err := readAll(s.ReadSource("invalid_url1"))
	if err == nil || len(data) != 0 {
		t.Fatal("Unexpected source", string(data), err)
	}

	s.AddPublicUrlSource("invalid_url2", "https://www.google.com/invalid_url")
	data, err = readAll(s.ReadSource("invalid_url2"))
	if err == nil || len(data) != 0 {
		t.Fatal("Unexpected source", string(data), err)
	}
}

func TestSourcesCustomHTTPClient(t *testing.T) {
	err := errors.New("some error")
	s := ast.Sources{}
	s.Client = &http.Client{}
	s.Client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, err
	})
	s.AddPublicUrlSource("robots.txt", "https://www.google.com/robots.txt")
	data, err2 := readAll(s.ReadSource("robots.txt"))
	if err2 == nil || !strings.Contains(err2.Error(), err.Error()) || len(data) != 0 {
		t.Fatal("Unexpected source", string(data), err2)
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func readAll(r io.ReadCloser) ([]byte, error) {
	defer r.Close()
	return ioutil.ReadAll(r)
}
