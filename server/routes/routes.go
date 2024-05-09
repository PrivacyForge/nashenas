package routes

import (
	"github.com/PrivacyForge/nashenas/handlers"
	"github.com/PrivacyForge/nashenas/middlewares"
	"github.com/gofiber/fiber/v2"
)

func DefineRoutes(app *fiber.App) *fiber.App {
	app.Get("/", middlewares.BearerToken, handlers.HelloWorld)
	app.Get("/me", middlewares.BearerToken, handlers.GetMe)
	app.Post("/set-username", middlewares.BearerToken, handlers.SetUsername)
	app.Post("/set-key", middlewares.BearerToken, handlers.SetPublicKey)
	app.Get("/get-messages", middlewares.BearerToken, handlers.GetMessages)
	app.Post("/send-message", handlers.SendMessage)
	app.Get("/profile/:username", handlers.GetProfile)
	app.Get("/confirm/:otp", handlers.ConfirmOTP)

	return app
}
