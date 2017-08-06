package image

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/image"
	"google.golang.org/appengine/log"
)

func handleUploadURL(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "No file found", http.StatusBadRequest)
		return
	}

	bkey := file[0].BlobKey
	url, _ := image.ServingURL(ctx, bkey, &image.ServingURLOptions{
		Secure: true,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(url.String()))
}

func init() {
	r := http.NewServeMux()
	r.HandleFunc("/url", handleUploadURL)
	r.HandleFunc("/upload", handleUpload)
	http.Handle("/", r)
}
