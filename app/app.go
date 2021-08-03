package app

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"ws/model"
	"ws/store"
)

type ImageCollectorApp interface {
	AddImageInformation(ctx context.Context, fh *multipart.FileHeader) error
}

type imageCollector struct {
	db store.DataStore
}

func NewImageCollector(d store.DataStore) ImageCollectorApp {
	return &imageCollector{
		db: d,
	}
}

// AddImageInformation Store the uploaded file locally and save metadata into DB
func (ic *imageCollector) AddImageInformation(ctx context.Context, fh *multipart.FileHeader) error {

	metadata := model.Metadata{}

	metadata.Name = fh.Filename
	metadata.Size = fh.Size

	fileType, err := saveImage(fh)
	if err != nil {
		return err
	}

	metadata.ContentType = fileType
	// save the metadata to db
	err = ic.db.StoreFileMetaData(ctx, metadata)
	if err != nil {
		return nil
	}

	return nil
}

func saveImage(fh *multipart.FileHeader) (string, error) {
	// read the content of uploaded file
	upFile, err := fh.Open()
	if err != nil {
		return "", err
	}

	defer upFile.Close()

	buff := make([]byte, 512)
	_, err = upFile.Read(buff)
	if err != nil {
		return "", err
	}

	// identify the content type for -JPEG or PNG
	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" {
		return "", err
	}

	// reset the pointer back to start of the file
	_, err = upFile.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	// store all the uploaded files under files folder
	err = os.MkdirAll("./files", os.ModePerm)
	if err != nil {
		return "", err
	}

	// create a destination file - .fileNameTime
	saveFile, err := os.Create(fmt.Sprintf("./files/%s%d", filepath.Ext(fh.Filename), time.Now().UnixNano()))
	if err != nil {
		return "", err
	}

	defer saveFile.Close()

	_, err = io.Copy(saveFile, upFile)
	if err != nil {
		return "", err
	}

	return fileType, nil
}
