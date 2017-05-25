package utils

import (
	"os"
	"testing"

	"github.com/bufferapp/buffer-static-upload/utils"
)

func TestGetFileMd5(t *testing.T) {
	file, _ := os.Open("hash-test.txt")
	hash, _ := utils.GetFileMd5(file)
	expected := "bde181f4a40d6e56662537d4e643a487"
	if hash != expected {
		t.Error("Expected", expected, "got", hash)
	}
}
