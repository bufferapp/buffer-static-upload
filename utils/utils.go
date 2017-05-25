package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// GetFileMd5 returns a checksum for a given file
func GetFileMd5(file *os.File) (string, error) {
	var fileHash string
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fileHash, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	fileHash = hex.EncodeToString(hashInBytes)
	return fileHash, nil
}

// GetVersionedFilename returns a new filename with the version before the extension
func GetVersionedFilename(filename string, version string) string {
	ext := filepath.Ext(filename)
	versionedExt := "." + version + ext
	versionedFilename := strings.Replace(filename, ext, versionedExt, 1)
	return versionedFilename
}

// GetFileMimeType returns the mime type of a file using it's extension
func GetFileMimeType(filename string) string {
	ext := filepath.Ext(filename)
	return mime.TypeByExtension(ext)
}

// GetFilesFromGlobsList returns a list of files that match a list of
// comma-deliniated file path globs
func GetFilesFromGlobsList(globList string) ([]string, error) {
	var files []string
	globs := strings.Split(globList, ",")

	for _, glob := range globs {
		fileList, err := filepath.Glob(glob)
		if err != nil {
			return files, err
		}
		files = append(files, fileList...)
	}
	return files, nil
}

// GetS3Uploader returns a configured Uploader
func GetS3Uploader() (*s3manager.Uploader, error) {
	var uploader *s3manager.Uploader

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(endpoints.UsEast1RegionID),
	}))

	_, err := creds.Get()
	if err != nil {
		fmt.Printf("Bad credentials: %s", err)
		return uploader, err
	}

	uploader = s3manager.NewUploader(sess)
	return uploader, nil
}

// VersionAndUploadFiles will verion files and upload them to s3 and return
// a map of filenames and their version hashes
func VersionAndUploadFiles(bucket string, filenames []string) (map[string]string, error) {
	fileVersions := map[string]string{}

	uploader, err := GetS3Uploader()
	if err != nil {
		return fileVersions, err
	}

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return fileVersions, err
		}
		defer file.Close()

		checksum, err := GetFileMd5(file)
		if err != nil {
			return fileVersions, err
		}

		fileVersions[filename] = checksum
		versionedFilename := GetVersionedFilename(filename, checksum)
		mimeType := GetFileMimeType(filename)

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:       aws.String(bucket),
			Key:          aws.String(versionedFilename),
			ContentType:  aws.String(mimeType),
			CacheControl: aws.String("public, max-age=31520626"),
			Expires:      aws.Time(time.Now().AddDate(10, 0, 0)),
			Body:         file,
		})
		if err != nil {
			return fileVersions, err
		}

		fmt.Printf("Uploaded %s\n", result.Location)
	}

	return fileVersions, nil
}
