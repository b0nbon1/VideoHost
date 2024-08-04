package stream

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/matoous/go-nanoid/v2"
)

func ProcessHLSHandler(c *fiber.Ctx) error {
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["videos"]

		if len(files) == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		file := files[0]

		fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

		folder := generateFolderWithId()
		filename :=  strings.ReplaceAll(file.Filename, " ", "-")

		if err := c.SaveFile(file, fmt.Sprintf("./%s/%s", folder, filename)); err != nil {
			return err
		}

		go hlsConversion(filename, folder)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success",
			"content": fiber.Map{
				"file": fmt.Sprintf("%s/%s/%s", c.BaseURL(), folder, filename),
			},
		})
	}

	return c.SendStatus(fiber.StatusInternalServerError)
}

func hlsConversion(filename string, folder string) {
	args := []string{"-i", filename, "-hls_time", "10", "-hls_playlist_type", "vod", "-hls_segment_filename", "segment_%d", "index.m3u8"}

	cmd := exec.Command("ffmpeg", args...)

	cmd.Dir = fmt.Sprintf("./%s", folder)

	out, err := cmd.Output();
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out);
	}
}

func generateFolderWithId() string {
	id, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)
	if err != nil {
		return ""
	}

	if _, err := os.Stat(fmt.Sprintf("./static/videos/%s", id)); os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("./static/videos/%s", id), 0755)
	} else {
		return generateFolderWithId()
	}

	return fmt.Sprintf("static/videos/%s", id)
}
