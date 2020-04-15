package ast

import "fmt"

// ParseError is returned for all errors
type ParseError struct {
	Reason string
	Source string
	Offset int
}

// Error implements the error interface
func (p *ParseError) Error() string {
	return fmt.Sprintf("%s at %s:%d", p.Reason, p.Source, p.Offset)
}
