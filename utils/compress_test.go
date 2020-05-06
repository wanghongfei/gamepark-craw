package utils

import (
	"testing"
)

func TestZipFile(t *testing.T) {
	err := ZipFile("../go.sum", "rd.zip")
	if nil != err {
		t.Fatal(err)
	}
}
