package controller

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
)

type response struct {
	Message string `json:"message"`
}

type responseurl struct {
	Message string `json:"message"`
	Url     string `json:"url"`
}

type ResponseUrlMultiple struct {
	Message string   `json:"message"`
	Url     []string `json:"urls"`
}

const (
	bucket_name = "service-rocket-buck-211"
	region      = "ap-southeast-1"
)

func isImage(fileType string) bool {
	fileType = strings.TrimPrefix(fileType, "image/")
	imageTypes := [...]string{"jpg", "png", "jpeg"}
	for _, imageType := range imageTypes {
		if imageType == fileType {
			return true
		}
	}
	return false
}

func isZip(fileType string) bool {
	fileType = strings.TrimPrefix(fileType, "application/")
	imageTypes := [...]string{"zip"}
	for _, imageType := range imageTypes {
		if imageType == fileType {
			return true
		}
	}
	return false
}

func UploadImage(c *fiber.Ctx) error {
	var file_type string
	var file_name_filtered string
	var url string
	message := response{"Unsupported file"}
	messageSuccess := response{"image uploaded"}
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(200).JSON(&messageSuccess)
	}
	images := form.File["image"]

	for _, image := range images {
		file_type = image.Header["Content-Type"][0]
		fmt.Println(image.Filename, image.Size, image.Header["Content-Type"][0])
		file_name_filtered = strings.ReplaceAll(image.Filename, " ", "+")
		if !isImage(file_type) {
			return c.Status(400).JSON(&message)
		}

		err := c.SaveFile(image, fmt.Sprintf("./%s", image.Filename))

		if err != nil {
			return c.Status(400).JSON(&message)
		}

		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)

		if err != nil {
			return c.Status(400).JSON(&message)
		}

		// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
		// for more information on configuring part size, and concurrency.
		//
		// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
		uploader := s3manager.NewUploader(sess)

		file, err := os.Open(image.Filename)
		if err != nil {
			return c.Status(400).JSON(&message)
		}

		defer file.Close()

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket_name),
			Key:    aws.String(image.Filename),
			Body:   file,
		})
		if err != nil {
			// Print the error and exit.
			return c.Status(400).JSON(&message)
		}
		url = fmt.Sprintf("https://service-rocket-buck-211.s3.ap-southeast-1.amazonaws.com/%s", file_name_filtered)

	}

	return c.Status(200).JSON(&responseurl{Message: "Upload Successful!", Url: url})
}

func UploadMultipleImage(c *fiber.Ctx) error {
	var file_name_filtered string
	var file_type string
	var url string
	var urls []string
	dst := "output"

	form, err := c.MultipartForm()

	if err != nil {
		return c.Status(400).JSON(err)
	}

	files := form.File["file"]

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	for _, file := range files {
		file_type = file.Header["Content-Type"][0]

		if !isZip(file_type) {
			return c.Status(400).JSON(err)
		}

		err := c.SaveFile(file, fmt.Sprintf("./output/%s", file.Filename))

		if err != nil {
			return c.Status(400).JSON(err)
		}

		archive, err := zip.OpenReader(file.Filename)

		if err != nil {
			panic(err)
		}

		for _, f := range archive.File {
			filePath := filepath.Join(dst, f.Name)
			fmt.Println("unzipping file ", filePath)

			if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
				fmt.Println("invalid file path")
				return c.Status(400).JSON(err)
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

			// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
			// for more information on configuring part size, and concurrency.
			//
			// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
			uploader := s3manager.NewUploader(sess)

			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucket_name),
				Key:    aws.String(f.Name),
				Body:   fileInArchive,
			})
			if err != nil {
				// Print the error and exit.
				return c.Status(400).JSON(err)
			}
			file_name_filtered = strings.ReplaceAll(f.Name, " ", "+")

			url = fmt.Sprintf("https://service-rocket-buck-211.s3.ap-southeast-1.amazonaws.com/%s", file_name_filtered)
			urls = append(urls, url)
			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				panic(err)
			}

			dstFile.Close()
			fileInArchive.Close()
		}
	}

	//Clear folder output after upload to S3

	dir, err := ioutil.ReadDir("output")
	if err != nil {
		return c.Status(400).JSON(err)
	}
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"output", d.Name()}...))
	}
	return c.Status(200).JSON(&ResponseUrlMultiple{Message: "Upload Successful!", Url: urls})
}
