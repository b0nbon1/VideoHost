package videos

import (
	"github.com/gofiber/fiber/v2"
)

func VideosRoutes(router fiber.Router) {
	// router.Get("/todos", getAllTodos)
	// router.Get("/todos/:id", getTodoById)
	router.Get("/videos", GetVideos)
	router.Post("/videos", CreateVideo)
	// router.Delete("/todos/:id", deleteTodoById)
	// router.Patch("/todos/:id", updateTodoById)
}