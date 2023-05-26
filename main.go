package main

import (
	"bytes"
	"image"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/resize-image", func(c *fiber.Ctx) error {
		// Receive image from POST request
		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(400).SendString("Bad Request")
		}

		// Open image file
		openFile, err := file.Open()
		if err != nil {
			return c.Status(500).SendString("Could not open image file")
		}

		// Decode the image
		img, _, err := image.Decode(openFile)
		if err != nil {
			return c.Status(500).SendString("Could not decode image file")
		}

		height := img.Bounds().Max.Y
		width := img.Bounds().Max.X

		// Get scale from form data
		scale, err := strconv.Atoi(c.FormValue("scale"))
		if err != nil {
			return c.Status(400).SendString("Invalid scale value")
		}

		// Calculate the scale ratio
		scaleRatio := float64(scale) / 100.0

		// Resize the image
		resizedImg := imaging.Resize(img, int(float64(width)*scaleRatio), int(float64(height)*scaleRatio), imaging.Lanczos)

		// Create a buffer to store the encoded resized image
		buf := new(bytes.Buffer)
		err = imaging.Encode(buf, resizedImg, imaging.JPEG)
		if err != nil {
			return c.Status(500).SendString("Could not encode resized image")
		}

		// Create a new ReadCloser for the fiber response
		readCloser := ioutil.NopCloser(bytes.NewReader(buf.Bytes()))

		// Set the headers and status code
		c.Response().Header.Set(fiber.HeaderContentType, "image/jpeg")
		c.Response().Header.Set(fiber.HeaderContentDisposition, "attachment; filename=resize.jpg")
		c.Response().SetStatusCode(http.StatusOK)

		// Send the image data as the fiber response
		return c.SendStream(readCloser, int(buf.Len()))
	})

	app.Listen(":3000")
}
