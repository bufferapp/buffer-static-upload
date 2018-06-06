package main

import (
	"testing"
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
