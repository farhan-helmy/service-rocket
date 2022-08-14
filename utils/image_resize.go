package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func isJpeg(fileName string) bool {

	fileExtension := filepath.Ext(fileName)
	//fmt.Printf("this is gfile extension %s", fileExtension)

	return fileExtension == ".jpeg"

	//return false
	//fmt.Println("this is file extension", fileExtension)
}

func GetImageDimension(fileName string) []string {
	urls := []string{}
	dst := "output"
	filePath := filepath.Join(dst, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filePath, err)
	}

	if image.Width >= 128 && image.Height >= 128 {
		urls = ResizeImage(fileName)
	} else {
		UploadToS3(filePath)
	}

	defer file.Close()

	return urls
}

func ResizeImage(fileName string) []string {
	var img image.Image
	var newFile1 string
	var newFile2 string
	urls := []string{}
	// var url1 string
	// var url2 string
	dst := "output"
	filePath := filepath.Join(dst, fileName)

	//fmt.Printf("this is file to resize %s", fileName)

	// open "test.jpg"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	is_jpeg := isJpeg(fileName)
	// if is_jpeg nak generate jpeg kat sini tapi png nak generate png,
	if is_jpeg {
		img, err = jpeg.Decode(file)

	} else {
		img, err = png.Decode(file)
	}

	if err != nil {
		log.Fatal(err)
	}
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Thumbnail(32, 40, img, resize.Lanczos3)

	m2 := resize.Thumbnail(64, 80, img, resize.Lanczos3)

	if is_jpeg {
		newFile1 = fmt.Sprintf("./output/resized_32x40_%s.jpeg", fileName)
		newFile2 = fmt.Sprintf("./output/resized_64x80_%s.jpeg", fileName)
	} else {
		newFile1 = fmt.Sprintf("./output/resized_32x40_%s.png", fileName)
		newFile2 = fmt.Sprintf("./output/resized_64x80_%s.png", fileName)
	}

	out, err := os.Create(newFile1)
	if err != nil {
		log.Fatal(err)
	}

	out2, err := os.Create(newFile2)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	// write new image to file2
	jpeg.Encode(out2, m2, nil)

	link1 := UploadToS3(newFile1)
	link2 := UploadToS3(newFile2)

	urls = append(urls, link1)
	urls = append(urls, link2)

	//urls = append(urls, url1, url2)

	defer file.Close()
	return urls
}
