package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/PrivacyForge/nashenas/configs"
	"github.com/PrivacyForge/nashenas/database"
	"github.com/PrivacyForge/nashenas/routes"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	configs.LoadConfigs()

	app := fiber.New()
	
	app.Use(logger.New())

	if err := database.InitConnection(); err != nil {
		panic("database connection failed.")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	routes.DefineRoutes(app)

	app.Listen(":" + configs.ServerPort)
}
