package routes

import (
	"github.com/PrivacyForge/nashenas/handlers"
	"github.com/PrivacyForge/nashenas/middlewares"
	"github.com/gofiber/fiber/v2"
)

func DefineRoutes(app *fiber.App) {
	app.Get("/", middlewares.AuthMiddleware, handlers.HelloWorld)

	app.Get("/me", middlewares.AuthMiddleware, handlers.GetMe)
	app.Get("/get-messages", middlewares.AuthMiddleware, handlers.GetMessages)
	app.Get("/profile/:username", handlers.GetProfile)
	app.Get("/get-key/:id", middlewares.AuthMiddleware, handlers.GetPublicKey)
	app.Post("/set-username", middlewares.AuthMiddleware, handlers.SetUsername)
	app.Post("/send-message", middlewares.AuthMiddleware, handlers.SendMessage)
	app.Post("/set-key", middlewares.AuthMiddleware, handlers.SetPublicKey)
	app.Post("/reply-message", middlewares.AuthMiddleware, handlers.ReplyMessage)
}
