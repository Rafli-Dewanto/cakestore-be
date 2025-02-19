package controller

import (
	"cakestore/internal/entity"
	"cakestore/internal/model"
	"cakestore/internal/usecase"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CakeController struct {
	cakeUseCase usecase.CakeUseCase
	logger      *logrus.Logger
	validator   *validator.Validate
}

func NewCakeController(cakeUseCase usecase.CakeUseCase, logger *logrus.Logger) *CakeController {
	return &CakeController{
		cakeUseCase: cakeUseCase,
		logger:      logger,
		validator:   validator.New(),
	}
}

func (c *CakeController) GetAllCakes(ctx *fiber.Ctx) error {
	cakes, err := c.cakeUseCase.GetAllCakes()
	if err != nil {
		c.logger.Errorf("Failed to fetch cakes: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch cakes",
		})
	}

	cakeResponses := make([]*model.CakeModel, len(cakes))
	for i, cake := range cakes {
		cakeResponses[i] = model.CakeToResponse(&cake)
	}

	return ctx.JSON(fiber.Map{
		"data": cakeResponses,
	})
}

func (c *CakeController) GetCakeByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cake ID",
		})
	}

	cake, err := c.cakeUseCase.GetCakeByID(id)
	if err != nil {
		c.logger.Errorf("Failed to get cake: %v", err)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cake not found",
		})
	}

	return ctx.JSON(model.CakeToResponse(cake))
}

func (c *CakeController) CreateCake(ctx *fiber.Ctx) error {
	var request model.CreateUpdateCakeRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("Failed to parse body: ", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	if err := c.validatePayload(request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	cake := &entity.Cake{
		Title:       request.Title,
		Description: request.Description,
		Rating:      float64(request.Rating),
		Image:       request.ImageURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   sql.NullTime{},
	}

	if err := c.cakeUseCase.CreateCake(cake); err != nil {
		c.logger.Error("Failed to create cake: ", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create cake"})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Cake created successfully",
		"data":    model.CakeToResponse(cake),
	})
}

func (c *CakeController) UpdateCake(ctx *fiber.Ctx) error {
	var request model.CreateUpdateCakeRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	if err := c.validatePayload(request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	cakeID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid cake ID"})
	}

	cake := &entity.Cake{
		ID:          cakeID,
		Title:       request.Title,
		Description: request.Description,
		Rating:      float64(request.Rating),
		Image:       request.ImageURL,
		UpdatedAt:   time.Now(),
	}

	if err := c.cakeUseCase.UpdateCake(cake); err != nil {
		c.logger.Error("Failed to update cake: ", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update cake"})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Cake updated successfully",
		"data":    model.CakeToResponse(cake),
	})
}

func (c *CakeController) DeleteCake(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cake ID",
		})
	}

	err = c.cakeUseCase.DeleteCake(id)
	if err != nil {
		c.logger.Errorf("Failed to delete cake: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete cake",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Cake deleted successfully",
	})
}

func (c *CakeController) validatePayload(request model.CreateUpdateCakeRequest) error {
	if err := c.validator.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errMessages := make([]string, len(validationErrors))
		for i, e := range validationErrors {
			errMessages[i] = "Field '" + e.Field() + "' failed on '" + e.Tag() + "' rule"
		}
		return fiber.NewError(http.StatusBadRequest, "Validation failed: " + strings.Join(errMessages, ", "))
	}
	return nil
}
