package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	bucket_name = "service-rocket-buck-211"
	region      = "ap-southeast-1"
)

func removeWhiteSpace(fileName string) string {
	return strings.ReplaceAll(fileName, " ", "+")
}

func removePath(fileName string) string {
	return strings.ReplaceAll(fileName, "./output/", "")
}
func UploadToS3(fileName string) string {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		panic(err)
	}

	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(sess)

	fileSetupPath := filepath.Join(fileName)
	fmt.Println(fileSetupPath)

	file, err := os.Open(fileSetupPath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket_name),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		// Print the error and exit.
		panic(err)
	}
	filtered_file_name := removePath(fileName)
	url := fmt.Sprintf("https://service-rocket-buck-211.s3.ap-southeast-1.amazonaws.com/%s", removeWhiteSpace(filtered_file_name))

	//fmt.Println(url)

	return url
}
