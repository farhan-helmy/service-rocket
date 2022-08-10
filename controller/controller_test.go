package controller

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUploadImage(t *testing.T) {
	filePath := "download.png"
	fieldName := "image"
	body := new(bytes.Buffer)

	mw := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}

	w, err := mw.CreateFormFile(fieldName, filePath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(w, file); err != nil {
		t.Fatal(err)
	}

	// close the writer before making the request
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/image/upload", body)

	req.Header.Add("Content-Type", mw.FormDataContentType())

	res := httptest.NewRecorder()
	fmt.Println(res)
}
