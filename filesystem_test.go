package util

import (
	"testing"

	"github.com/spf13/afero"
)

func TestSanitizeExtractPathPass(t *testing.T) {
	err := sanitizeExtractPath("fake.file", "fakePath")
	if err != nil {
		t.Error()
	}
}

func TestSanitizeExtractPathFail(t *testing.T) {
	err := sanitizeExtractPath("../../fake.file", "fakePath")
	if err == nil {
		t.Error()
	}
}

func TestUnZip(t *testing.T) {
	fs := afero.NewOsFs()
	err := UnZip("__testdata__/test.zip", "__testdata__/test-zip")
	if err != nil {
		t.Error()
	}
	_ = fs.RemoveAll("__testdata__/test-zip")
}

func TestFindFiles(t *testing.T) {
	fs := afero.NewOsFs()
	files, err := FindFiles(fs, "__testdata__", "(.*)\\.txt")
	if len(files) != 1 {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}
