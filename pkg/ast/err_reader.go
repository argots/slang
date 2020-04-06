package ast

// errReader always returns an error for any read call
type errReader struct {
	err error
}

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func (e errReader) Close() error {
	return nil
}
