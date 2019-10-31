package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/warkio/logger-go"

	"github.com/warkio/pomf-go"
	storagefs "github.com/warkio/pomf-go/storage/filesystem"
)

func main() {
	log := logger.StandardOutput

	addr := os.Getenv("POMF_LISTEN_ADDRESS")
	if addr == "" {
		addr = ":3000"
	}

	log.Printf("Address: %#v", addr)

	urlPrefix := os.Getenv("POMF_URL_PREFIX")
	if urlPrefix == "" {
		urlPrefix = "/"
	}

	log.Printf("URL prefix: %#v", urlPrefix)

	uploadDirectory := os.Getenv("POMF_UPLOAD_DIRECTORY")
	if uploadDirectory == "" {
		uploadDirectory = "files"
	}

	var err error

	uploadDirectory, err = filepath.Abs(uploadDirectory)
	if err != nil {
		log.Print(err)

		os.Exit(1)
	}

	log.Printf("Upload directory: %#v", uploadDirectory)

	storage := &storagefs.Storage{
		UploadDirectory: "files",
	}
	storage.SetURLPrefix(urlPrefix)

	s := &pomf.Pomf{
		Storage: storage,
		Logger:  log,
	}

	err = s.SelfCheck()
	if err != nil {
		log.Print(err)

		os.Exit(1)
	}

	err = http.ListenAndServe(addr, s)
	if err != nil {
		log.Print(err)
	}
}
