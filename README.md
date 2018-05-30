# buffer-static-upload

A straightforward static asset uploader which versions files by their contents
and generates a json file for your application to read the uploaded file
locations from.

## Versioned files and Images

`.js` and `.css` files are versioned during upload using a hash of the file's
contents as to allow for cache-busting.

Images and other files are *not* versioned to allow for maximum caching and due
to their contents not changing very often like `.css` and `.js` files do.

## Install

A pre-compiled binary is available for download for both Linux and macOS.
Replace the version (ex. `0.2.0`) in the URL below for the version you require:

```
curl -L https://github.com/bufferapp/buffer-static-upload/releases/download/0.2.0/buffer-static-upload-`uname -s` > /usr/local/bin/buffer-static-upload
chmod +x /usr/local/bin/buffer-static-upload
```

## Usage

Ensure your AWS credentials environment variables are set (`AWS_ACCESS_KEY_ID`,
`AWS_SECRET_ACCESS_KEY`). The cli has the following argument options:

```
$ buffer-static-upload -h
Usage of buffer-static-upload:
  -bucket string
      the s3 bucket to upload to (default "static.buffer.com")
  -dir string
      required, the directory to upload files to in the bucket
  -dry-run
      print the output only, skip file uploads and manifest creation
  -files string
      the path to the files you'd like to upload, ex. "public/**/.*js,public/style.css"
  -format string
      format of the output [json,csv] (default "json")
  -o string
      the filename for the versions manifest (default "staticAssets.json")
  -v	print the current buffer-static-upload version
```

For example, you can use glob patterns to match multiple sets of files:

```
buffer-static-upload -files "public/js/**/*.js,public/css/*.css,public/img/*.*" -bucket my-bucket
```

This will generate a `staticAssets.json` file in this directory like this:

```json
{
  "public/css/style.css": "https://my-bucket.s3.amazonaws.com/public/css/style.11985b07e3121564a73d4d6821bfcfe7.css",
  "public/js/x/another.js": "https://my-bucket.s3.amazonaws.com/public/js/x/another.bfa2d0f60841707efe7be0a94c4caacf.js",
  "public/js/script.js": "https://my-bucket.s3.amazonaws.com/public/js/script.d55002b60fcfff0b3d355184d23af6f7.js",
  "public/img/home.jpg": "https://my-bucket.s3.amazonaws.com/public/img/home.jpg",
}
```

*Note* - The default bucket is used by multiple teams, so if you use that you
must specify a directory to use for your project as not to create unnecessary
collisions.

### Development

To work on this project you'll need [Golang](https://golang.org/dl/) and
the [Glide package manager](https://glide.sh/) installed. You should have
your `GOPATH` environment variable set and this repo should be cloned within
your `$GOPATH/src/github.com/bufferapp` directory.

To install the dependencies, run:

```
glide install
```

To test the script run:

```
go run main.go <your cli arguments here>
```

When distributing a new release version, run this script to generate the
binaries for Linux and Mac:

```
./build.sh
```

### License

MIT
