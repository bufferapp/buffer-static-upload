package main

import (
	"flag"
	"fmt"

	"github.com/bufferapp/buffer-static-upload/utils"
)

func main() {
	s3Bucket := "buffer-dan-test"

	filesArg := flag.String("files", "", "the path to the files you'd like to upload")
	flag.Parse()

	files, err := utils.GetFilesFromGlobsList(*filesArg)
	if err != nil {
		fmt.Printf("err %s", err)
	}
	fmt.Printf("Found %d files to upload and version\n", len(files))

	fileVersions, err := utils.VersionAndUploadFiles(s3Bucket, files)
	if err != nil {
		fmt.Printf("Failed to upload files: %s", err)
	}

	for file, version := range fileVersions {
		fmt.Printf("%s = %s\n", file, version)
	}
}
