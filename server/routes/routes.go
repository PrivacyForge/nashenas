package routes

import (
	"github.com/PrivacyForge/nashenas/handlers"
	"github.com/PrivacyForge/nashenas/middlewares"
	"github.com/gofiber/fiber/v2"
)

func DefineRoutes(app *fiber.App) *fiber.App {
	app.Get("/", middlewares.RequiredBearerToken, handlers.HelloWorld)
	app.Get("/me", middlewares.RequiredBearerToken, handlers.GetMe)
	app.Get("/get-messages", middlewares.RequiredBearerToken, handlers.GetMessages)
	app.Get("/profile/:username", handlers.GetProfile)
	app.Get("/confirm/:otp", handlers.ConfirmOTP)
	app.Post("/set-username", middlewares.RequiredBearerToken, handlers.SetUsername)
	app.Post("/set-key", middlewares.RequiredBearerToken, handlers.SetPublicKey)
	app.Post("/send-message", middlewares.OptionalBearerToken, handlers.SendMessage)

	return app
}
