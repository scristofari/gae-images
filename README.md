# gae-images
Google app engine - Handling upload / resize images

## Dependency: 

* [blobstore API](https://cloud.google.com/appengine/docs/go/images/) 

## What does the microservice do: 

1. Provide an url for uploading images ( one-shot ).
2. Next you can, upload the image and get the url for it's rendering.
3. Bonus, with the blobstore api, you can [resize and crop dynamically adding (=sXX/=sXX-c)](https://cloud.google.com/appengine/docs/go/images/#Go_Serving_and_re-sizing_images_from_the_Blobstore)
at the end of the url

### The app.yaml.

As simple as possible.

* All the requests are sent to the go app.
* The service will be served in 'https'.

``` yaml
application: "your-app"
version: 1
runtime: go
api_version: go1 

handlers:
- url: /.*
  script: _go_app          # All the request are sent to the go app.
  secure: always           # the service will be in 'https'
```
### The app.go

* Two routes are provided:

 "/url", provide a one-shot upload url. (GET)

  "/upload", the url used by the blobstore to upload the file, and get the image url after that. (POST)

``` go
func init() {
	r := http.NewServeMux()
	r.HandleFunc("/url", handleUploadURL)
	r.HandleFunc("/upload", handleUpload)
	http.Handle("/", r)
}
```

* Get the url for the upload.

``` go
// Generate the url for the upload - url 'one-shot'
func handleGetUrlForUpload(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  uploadURL, err := blobstore.UploadURL(ctx, "/upload", nil)
  if err != nil {
    log.Errorf(ctx, err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(uploadURL.String()))
}
```

* Handle the upload and provide the url for rendering it.

Get the file on the key = 'file'.

Serve the blob with a secure url.

``` go
// Upload the file
func handleUpload(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  blobs, _, err := blobstore.ParseUpload(r)
  if err != nil {
    log.Errorf(ctx, err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  file := blobs["file"]
  if len(file) == 0 {
    log.Errorf(ctx, err.Error())
    http.Error(w, "No file found or url upload not generated", http.StatusBadRequest)
    return
  }

  bkey := file[0].BlobKey
  url, _ := image.ServingURL(ctx, bkey, &image.ServingURLOptions{
    Secure: true,
  })

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(url))
}
```

### Resize - crop the image.

Everything is explained [here](https://cloud.google.com/appengine/docs/go/images/#Go_Serving_and_re-sizing_images_from_the_Blobstorew).
