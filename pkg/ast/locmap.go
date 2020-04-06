package ast

import (
	"io"
	"io/ioutil"
)

// Loc represents a location.
//
// This can be resolved to an offset using a location map.
type Loc uint32

// Offset returns the location for a token location
func (l Loc) Offset(lm LocMap) (source string, start, end uint32) {
	return lm.Get(l)
}

// Token returns the actual token given a source reader
func (l Loc) Token(lm LocMap, sources SourceReader) (string, error) {
	location, start, end := lm.Get(l)
	src := sources.ReadSource(location)
	defer src.Close()

	if _, err := io.CopyN(ioutil.Discard, src, int64(start)); err != nil {
		return "", err
	}
	result := make([]byte, end-start)
	if _, err := io.ReadAtLeast(src, result, int(end-start)); err != nil {
		return "", err
	}
	return string(result), nil
}

// LocMap implements a map of token offsets to a Loc handle
type LocMap interface {
	Get(handle Loc) (location string, start, end uint32)
	Add(location string, start, end uint32) (handle Loc)
}

func NewLocMap() LocMap {
	return &locMap{entriesMap: map[locEntry]Loc{}}
}

type locEntry struct {
	location   string
	start, end uint32
}

type locMap struct {
	entries    []locEntry
	entriesMap map[locEntry]Loc
}

func (l *locMap) Get(handle Loc) (location string, start, end uint32) {
	entry := l.entries[int(uint32(handle))]
	return entry.location, entry.start, entry.end
}

func (l *locMap) Add(location string, start, end uint32) (handle Loc) {
	entry := locEntry{location, start, end}
	handle, ok := l.entriesMap[entry]
	if !ok {
		handle = Loc(uint32(len(l.entries)))
		l.entries = append(l.entries, entry)
		l.entriesMap[entry] = handle
	}
	return handle
}
