package controller

import (
	"fmt"
	"os"
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
