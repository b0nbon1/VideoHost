package videos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	initdb "github.com/b0nbon1/VidFlux/database"
	"github.com/b0nbon1/VidFlux/pkg/videos/db"
	"github.com/b0nbon1/VidFlux/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type (
	Video struct {
		Name        string `validate:"required,min=3,max=60"`
		VideoUrl    string `validate:"required,url"`
		Description string `validate:"required"`
	}
	ReturnedVideo struct {
		ID          uuid.UUID      `json:"id"`
		Name        string         `json:"name"`
		VideoUrl    string         `json:"video_url"`
		Description string `json:"description"`
	}
	VideoResponse struct {
		data    ReturnedVideo
		message string
	}
)

func mapVideo(video db.Video) interface{} {

	return ReturnedVideo{
		ID:          video.ID,
		Name:        video.Name,
		VideoUrl:    video.VideoUrl,
		Description: video.Description.String,
	}
}

type NullString sql.NullString

func (x *NullString) MarshalJSON() ([]byte, error) {
    if !x.Valid {
        return []byte("null"), nil
    }
    return json.Marshal(x.String)
}


// GetVideos func get Videos.
// @Description Be able to get a Videos
// @Summary get a Videos
// @Tags Videos
// @Accept json
// @Produce json
// @Success 200 {object} VideoResponse
// @Router /api/v1/videos [post]
func GetVideos(c *fiber.Ctx) error {
	queries := db.New(initdb.DB)

	videos, err := queries.GetAllVideos(c.Context())
	if err != nil {
		fmt.Printf("Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"content":    util.MapValues(videos, mapVideo),
	})
}

// CreateVideos func create Videos.
// @Description Be able to create a Videos
// @Summary create a Videos
// @Tags Videos
// @Accept json
// @Produce json
// @Success 200 {object} VideoResponse
// @Router /api/v1/videos [post]
func CreateVideo(c *fiber.Ctx) error {
	myValidator := &util.XValidator{}

	video := new(Video)
	if err := c.BodyParser(video); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fmt.Println(video)

	if errs := myValidator.Validate(video); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": strings.Join(errMsgs, " and "),
		})
	}

	queries := db.New(initdb.DB)
	
	videoResult, err := queries.CreateVideo(c.Context(), db.CreateVideoParams{
		Name:        video.Name,
		VideoUrl:    video.VideoUrl,
		Description: sql.NullString{String: video.Description, Valid: true},
	})

	if err != nil {
		fmt.Printf("Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"content": mapVideo(videoResult),
	})
}
