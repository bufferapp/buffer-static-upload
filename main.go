package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var defaultS3Bucket = "static.buffer.com"
var uploader *s3manager.Uploader
var svc *s3.S3

func fatal(format string, a ...interface{}) {
	s := "Error: " + format + "\n"
	if a != nil {
		fmt.Printf(s, a)
	} else {
		fmt.Print(s)
	}
	os.Exit(1)
}

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

// SetupS3Uploader configures and assigns the global "uploader" and "svc" variables
func SetupS3Uploader() {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(endpoints.UsEast1RegionID),
	}))
	_, err := creds.Get()
	if err != nil {
		fatal("failed to load AWS credentials %s", err)
	}

	uploader = s3manager.NewUploader(sess)
	svc = s3.New(sess)
}

// HasPreviousUpload performs a HEAD request to check if a file has been uploaded already
func HasPreviousUpload(svc *s3.S3, bucket string, filename string) bool {
	req, _ := svc.HeadObjectRequest(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	err := req.Send()
	if err == nil {
		return true
	}
	return false
}

// GetFileURL returns the final url of the file
func GetFileURL(filename string, bucket string) string {
	// the static.buffer.com bucket has a domain alias
	if bucket == defaultS3Bucket {
		return "https://" + path.Join(bucket, filename)
	}
	return "https://s3.amazonaws.com" + path.Join("/", filename)
}

// UploadFile ensures a given file is uploaded to the s3 bucket and returns
// the filename
func UploadFile(file *os.File, filename string, bucket string) (fileURL string, err error) {
	mimeType := GetFileMimeType(filename)

	var action string
	if !HasPreviousUpload(svc, bucket, filename) {
		_, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:       aws.String(bucket),
			Key:          aws.String(filename),
			ContentType:  aws.String(mimeType),
			CacheControl: aws.String("public, max-age=31520626"),
			Expires:      aws.Time(time.Now().AddDate(10, 0, 0)),
			Body:         file,
		})
		if err != nil {
			return fileURL, err
		}
		action = "Uploaded"
	} else {
		action = "Skipped"
	}

	fileURL = GetFileURL(filename, bucket)
	fmt.Printf("%-10s %s\n", action, fileURL)
	return fileURL, nil
}

// VersionAndUploadFiles will verion files and upload them to s3 and return
// a map of filenames and their version hashes
func VersionAndUploadFiles(
	bucket string,
	directory string,
	filenames []string,
) (map[string]string, error) {
	fileVersions := map[string]string{}

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return fileVersions, err
		}
		defer file.Close()

		ext := filepath.Ext(filename)
		uploadFilename := filename
		if ext == ".js" || ext == ".css" {
			checksum, errMd5 := GetFileMd5(file)
			if errMd5 != nil {
				return fileVersions, errMd5
			}
			uploadFilename = GetVersionedFilename(filename, checksum)
		}
		bucketFilename := path.Join(directory, uploadFilename)

		fileURL, err := UploadFile(file, bucketFilename, bucket)
		if err != nil {
			return fileVersions, err
		}
		fileVersions[filename] = fileURL
	}

	return fileVersions, nil
}

func main() {
	s3Bucket := flag.String("bucket", defaultS3Bucket, "the s3 bucket to upload to")
	directory := flag.String("dir", "", "required, the directory to upload files to in the bucket")
	filesArg := flag.String("files", "", "the path to the files you'd like to upload, ex. \"public/**/.*js,public/style.css\"")
	outputFilename := flag.String("o", "staticAssets.json", "the json file you'd like your generate")
	flag.Parse()

	if *directory == "" && *s3Bucket == defaultS3Bucket {
		fatal("To use the default bucket you need to specify an upload directory (-dir)")
	}

	start := time.Now()
	files, err := GetFilesFromGlobsList(*filesArg)
	if err != nil {
		fatal("failed to get files %s", err)
	}
	fmt.Printf("Found %d files to upload and version:\n", len(files))

	SetupS3Uploader()
	fileVersions, err := VersionAndUploadFiles(*s3Bucket, *directory, files)
	if err != nil {
		fatal("failed to upload files %s", err)
	}

	output, err := json.MarshalIndent(fileVersions, "", "  ")
	if err != nil {
		fatal("failed to generate versions json file %s", err)
	}

	err = ioutil.WriteFile(*outputFilename, output, 0644)
	if err != nil {
		fatal("failed to write versions json file %s", err)
	}

	elapsed := time.Since(start)
	fmt.Printf(
		"\nSuccessfully uploaded static assets and generated %s in %s\n",
		*outputFilename,
		elapsed,
	)
}
