package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnzipFile(fileName string) []string {
	collection := []string{}

	dst := "output"
	fileSetupPath := filepath.Join(dst, fileName)
	archive, err := zip.OpenReader(fileSetupPath)
	//fmt.Println(fileSetupPath)
	if err != nil {
		panic(err)
	}

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)

		macos := strings.Contains(filePath, "__MACOSX")

		if macos {
			continue
		}

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
		}

		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		// url = UploadToS3(f.Name)

		// urls = append(urls, url)

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
		urls := GetImageDimension(f.Name)
		for _, url := range urls {
			collection = append(collection, url)
		}

		//fmt.Println(url)

	}
	return collection
}
