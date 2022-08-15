package controller

import (
	"fmt"
	"strings"

	"github.com/farhan-helmy/service-rocket/utils"
	"github.com/gofiber/fiber/v2"
)

type response struct {
	Message string `json:"message"`
}

type responseurl struct {
	Message string   `json:"message"`
	Url     []string `json:"url"`
}

type ResponseUrlMultiple struct {
	Message string `json:"message"`
	Url     string `json:"urls"`
}

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
	var url []string
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

		if !isImage(file_type) {
			return c.Status(400).JSON(&message)
		}

		err := c.SaveFile(image, fmt.Sprintf("./output/%s", image.Filename))

		//url = utils.UploadToS3(image.Filename)

		url = utils.GetImageDimension(image.Filename)

		if err != nil {
			return c.Status(400).JSON(&message)
		}

	}

	defer utils.ClearFolder()

	return c.Status(200).JSON(&responseurl{Message: "Upload Successful!", Url: url})
}

func UploadMultipleImage(c *fiber.Ctx) error {
	var file_type string
	var urls []string

	form, err := c.MultipartForm()

	if err != nil {
		return c.Status(400).JSON(err)
	}

	files := form.File["file"]

	for _, file := range files {
		file_type = file.Header["Content-Type"][0]

		if !isZip(file_type) {
			return c.Status(400).JSON(err)
		}

		err := c.SaveFile(file, fmt.Sprintf("./output/%s", file.Filename))

		if err != nil {
			return c.Status(400).JSON(err)
		}

		urls = utils.UnzipFile(file.Filename)
		fmt.Println(urls)

	}

	defer utils.ClearFolder()

	return c.Status(200).JSON(urls)
}
