package test_utils

import "errors"

type DummyReader struct {
	Message string
}

func (r DummyReader) Read(p []byte) (n int, err error) {
	return 0, errors.New(r.Message)
}