package stream

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/matoous/go-nanoid/v2"
)

const hlsBasePath = "./static/videos"

func ProcessHLSHandler(c *fiber.Ctx) error {
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["videos"]

		if len(files) == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		file := files[0]

		fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

		folder, id := generateFolderWithId()
		filename := strings.ReplaceAll(file.Filename, " ", "-")

		if err := c.SaveFile(file, fmt.Sprintf("./%s/%s", folder, filename)); err != nil {
			return err
		}

		// check the reolution of the video

		go hlsConversion(filename, folder, id)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success",
			"content": fiber.Map{
				"video_id": id,
				"file":     fmt.Sprintf("%s/%s", folder, filename),
			},
		})
	}

	return c.SendStatus(fiber.StatusInternalServerError)
}

func hlsConversion(filename string, folder string, id string) {
	args := []string{"-i", filename, "-hls_time", "5", "-hls_playlist_type", "vod", "-hls_segment_filename", id + "%d.ts", "index.m3u8"}

	cmd := exec.Command("ffmpeg", args...)

	cmd.Dir = fmt.Sprintf("./%s", folder)

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}

	m3u8Path := filepath.Join(hlsBasePath, id, "index.m3u8")

	if _, err := os.Stat(m3u8Path); os.IsNotExist(err) {
		log.Fatal(err, "path not found!")
	}

	file, err := os.Open(m3u8Path)
	if err != nil {
		log.Fatal(err, "Error reading M3U8 file")
	}
	defer file.Close()

	var modifiedM3U8 string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if filepath.Ext(line) == ".ts" {
			line = fmt.Sprintf("/api/v1/videos/%s/segment/%s", id, line)
		}
		modifiedM3U8 += line + "\n"
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatal(err, "Error scanning M3U8 file")
	}

	// Step 2: Overwrite the file with modified content
	err = os.WriteFile(m3u8Path, []byte(modifiedM3U8), 0644)
	if err != nil {
		log.Fatal(err, "Error writing to M3U8 file")
	}

	fmt.Println("File successfully updated!")
}

func generateFolderWithId() (string, string) {
	id, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 6)
	if err != nil {
		return "", ""
	}

	if _, err := os.Stat(fmt.Sprintf("./static/videos/%s", id)); os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("./static/videos/%s", id), 0755)
	} else {
		return generateFolderWithId()
	}

	return fmt.Sprintf("static/videos/%s", id), id
}

func ProcessFetchStream(c *fiber.Ctx) error {
	videoID := c.Params("videoid")
	m3u8Path := filepath.Join(hlsBasePath, videoID, "index.m3u8")
	content, err := os.ReadFile(m3u8Path)
	if err != nil {
		log.Fatal(err)
		return c.SendStatus(500)
	}

	c.Set("Content-Type", "application/vnd.apple.mpegurl")
	return c.SendString(string(content))
}

func FetchSegments(c *fiber.Ctx) error {
	videoID := c.Params("videoid")
	segment := c.Params("segment")
	fmt.Println(videoID, segment)
	segmentPath := filepath.Join(hlsBasePath, videoID, segment)

	if _, err := os.Stat(segmentPath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString("Segment not found")
	}

	c.Set("Content-Type", "video/MP2T")
	return c.SendFile(segmentPath)
}
