# buffer-static-upload

A straightforward static asset uploader which generates a json file for your
application to read the uploaded file locations from.

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
  -files string
    	the path to the files you'd like to upload, ex. "public/**/.*js,public/style.css"
  -o string
    	the json file you'd like your generate (default "staticAssets.json")
```

For example, you can use glob patterns to match multiple sets of files:

```
buffer-static-upload -files "public/js/**/*.js,public/css/*.css"
```

*Note* - The default bucket is used by multiple teams, so if you use that you
must specify a directory to use for your project as not to create unnecessary
collisions.
