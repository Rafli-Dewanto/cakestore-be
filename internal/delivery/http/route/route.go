package route

import (
	http "cakestore/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type RouteConfig struct {
	App            *fiber.App
	CakeController *http.CakeController
}

func (c *RouteConfig) Setup() {
	c.SetupRoute()
}

func (c *RouteConfig) SetupRoute() {
	api := c.App.Group("/cakes")
	c.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	api.Get("/", c.CakeController.GetAllCakes)
	api.Get("/:id", c.CakeController.GetCakeByID)
	api.Post("/", c.CakeController.CreateCake)
	api.Put("/:id", c.CakeController.UpdateCake)
	api.Delete("/:id", c.CakeController.DeleteCake)
}
