package main

import (
	"github.com/gofiber/fiber/v2"
)

func DefineRoutes(app *fiber.App) *fiber.App {
	app.Get("/", BearerTokenMiddleware, HelloWorld)
	app.Get("/me", BearerTokenMiddleware, GetMe)
	app.Post("/set-username", BearerTokenMiddleware, SetUsername)
	app.Post("/set-key", BearerTokenMiddleware, SetPublicKey)
	app.Get("/get-messages", BearerTokenMiddleware, GetMessages)
	app.Post("/send-message", SendMessage)
	app.Get("/profile/:username", GetProfile)
	app.Get("/confirm/:otp", ConfirmOTP)

	return app
}
