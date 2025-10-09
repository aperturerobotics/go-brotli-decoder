// Copyright 2016 Google Inc. All Rights Reserved.
//
// Distributed under MIT license.
// See file LICENSE for detail or copy at https://opensource.org/licenses/MIT

package brotli

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func Decode(encodedData []byte) ([]byte, error) {
	r := NewReader(bytes.NewReader(encodedData))
	return ioutil.ReadAll(r)
}

func TestDecodeInvalidInput(t *testing.T) {
	content := []byte("---\nthis-is-not-brotli: \"it is actually yaml\"")
	input := bytes.NewBuffer(content)

	r := NewReader(input)

	buf, err := io.ReadAll(r)
	if err == nil {
		t.Fatalf("expected error, got none and read:\n%x\n%s\n%v", buf, buf, buf)
	}
}
