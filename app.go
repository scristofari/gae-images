package image

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/image"
	"google.golang.org/appengine/log"
)

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
	fmt.Fprint(w, uploadURL.String())
}

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
	fmt.Fprint(w, url)
}

func init() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/generate-url", handleGetUrlForUpload).Methods("GET")
	r.HandleFunc("/upload", handleUpload).Methods("POST")
	http.Handle("/", r)
}
