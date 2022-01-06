package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/laluardian/fiber-book-api/controllers"
)

func BookRoute(r fiber.Router) {
	r.Get("/", controllers.FetchBooks)
	r.Post("", controllers.InsertBook)
	r.Put("/:id", controllers.UpdateBook)
	r.Delete("/:id", controllers.DeleteBook)
	r.Get("/:id", controllers.FetchBook)
}
