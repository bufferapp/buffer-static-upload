package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestGetFileURL(t *testing.T) {
	tests := []struct {
		bucket         string
		bucketFilename string
		expected       string
	}{
		{"some-bucket", "js/app.1234.js", "https://s3.amazonaws.com/some-bucket/js/app.1234.js"},
		{"static.buffer.com", "dir/js/app.1234.js", "https://static.buffer.com/dir/js/app.1234.js"},
	}

	for _, test := range tests {
		actual := GetFileURL(test.bucket, test.bucketFilename)
		if actual != test.expected {
			t.Errorf("File URL was incorrect, got: %s, expected %s", actual, test.expected)
		}
	}
}

func TestShouldVersionFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"bundle.js", true},
		{"assets.css", true},
		{"another.file", false},
	}

	for _, test := range tests {
		actual := ShouldVersionFile(test.filename)
		if actual != test.expected {
			t.Errorf("ShouldVersionFile result was incorrect, got %t, expected %t for filename %s", actual, test.expected, test.filename)
		}
	}
}

func TestGetUploadFilename(t *testing.T) {
	var AppFs = afero.NewMemMapFs()
	filename := "bundle.js"
	file, _ := AppFs.OpenFile(filename, os.O_CREATE, 0600)
	file.WriteString("some JS content")
	tests := []struct {
		file           afero.File
		filename       string
		skipVersioning bool
		expected       string
	}{
		{file, filename, false, "bundle.d41d8cd98f00b204e9800998ecf8427e.js"},
		{file, filename, true, "bundle.js"},
	}

	for _, test := range tests {
		fmt.Print(file.Name())
		actual, _ := GetUploadFilename(test.file, test.filename, test.skipVersioning)
		if actual != test.expected {
			t.Errorf("GetUploadFilename result was incorrect, got %s, expected %s", actual, test.expected)
		}
	}
}
