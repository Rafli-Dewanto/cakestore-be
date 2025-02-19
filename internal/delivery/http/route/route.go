package route

import (
	http "cakestore/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
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

	api.Get("/", c.CakeController.GetAllCakes)
	api.Get("/:id", c.CakeController.GetCakeByID)
	api.Post("/", c.CakeController.CreateCake)
	api.Put("/:id", c.CakeController.UpdateCake)
	api.Delete("/:id", c.CakeController.DeleteCake)
}
