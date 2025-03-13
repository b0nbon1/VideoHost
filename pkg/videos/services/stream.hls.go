package stream

import (
	// "bufio"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

		// go hlsConversion(filename, folder, id)
		go hlsConversionWithResolutions(filename, folder, id)

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

func hlsConversionWithResolutions(filename string, folder string, id string) {

	width, height, err := getVideoResolution(fmt.Sprintf("./%s/%s", folder, filename))
	if err != nil {
		log.Fatal("Error detecting video resolution:", err)
	}
	fmt.Printf("Detected resolution: %dx%d\n", width, height)

	cmd := generateFFmpegCommand(fmt.Sprintf("./%s/%s", folder, filename), folder, width, height)
	if err := cmd.Run(); err != nil {
		log.Fatal("Error running FFmpeg:", err)
	}

	fmt.Println("HLS video processing complete!")

}

func HlsConversion(filename string, folder string, id string) {
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


type VideoStream struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ProbeResult struct {
	Streams []VideoStream `json:"streams"`
}

func getVideoResolution(videoPath string) (int, int, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_streams", videoPath)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, 0, fmt.Errorf("error running ffprobe: %v", err)
	}

	var result ProbeResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return 0, 0, fmt.Errorf("error parsing ffprobe output: %v", err)
	}

	// Find the highest resolution video stream
	var maxWidth, maxHeight int
	for _, stream := range result.Streams {
		if stream.Width > maxWidth {
			maxWidth = stream.Width
			maxHeight = stream.Height
		}
	}

	if maxWidth == 0 || maxHeight == 0 {
		return 0, 0, fmt.Errorf("could not determine video resolution")
	}

	return maxWidth, maxHeight, nil
}

func generateFFmpegCommand(inputFile string, outputDir string, width int, height int) *exec.Cmd {
	resolutions := []struct {
		width   int
		height  int
		bitrate string
	}{
		{3840, 2160, "8000k"}, // 4K (2160p)
		{2560, 1440, "5000k"}, // 2K (1440p)
		{1920, 1080, "3000k"}, // Full HD (1080p)
		{1280, 720, "1800k"},  // HD (720p)
		{854, 480, "1200k"},   // SD (480p)
		{640, 360, "800k"},    // 360p (Lowest)
	}

	// FFmpeg base arguments
	args := []string{
		"-y", "-i", inputFile,
		"-preset", "slow",
		"-g", "48",
		"-sc_threshold", "0",
	}

	// Add resolution mapping dynamically
	var streamMap []string
	streamIndex := 0
	for _, res := range resolutions {
		if res.width <= width && res.height <= height {
			args = append(args,
				"-map", "0:0", "-map", "0:1",
				"-s:v:"+strconv.Itoa(streamIndex), fmt.Sprintf("%dx%d", res.width, res.height),
				"-b:v:"+strconv.Itoa(streamIndex), res.bitrate,
			)
			streamMap = append(streamMap, fmt.Sprintf("v:%d,a:%d,name:%dp", streamIndex, streamIndex, res.height))
			streamIndex++
		}
	}

	// Append audio copy and output settings
	args = append(args,
		"-c:a", "copy",
		"-var_stream_map", strings.Join(streamMap, " "),
		"-master_pl_name", "master.m3u8",
		"-f", "hls",
		"-hls_time", "10",
		"-hls_playlist_type", "vod",
		"-hls_list_size", "0",
		"-hls_segment_filename", fmt.Sprintf("%s/v%%v/segment%%d.ts", outputDir),
		fmt.Sprintf("%s/v%%v/index.m3u8", outputDir),
	)

	// Print the full FFmpeg command (for debugging)
	fmt.Println("Running FFmpeg with arguments:", args)

	return exec.Command("ffmpeg", args...)
}

