package controller

import (
	"cakestore/internal/entity"
	"cakestore/internal/model"
	"cakestore/internal/usecase"
	"cakestore/utils"
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
	resp := new(utils.Response)

	cakes, err := c.cakeUseCase.GetAllCakes()
	if err != nil {
		c.logger.Errorf("Failed to fetch cakes: %v", err)
		resp.Message = "Failed to fetch cakes"
		resp.Success = false
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	cakeResponses := make([]*model.CakeModel, len(cakes))
	for i, cake := range cakes {
		cakeResponses[i] = model.CakeToResponse(&cake)
	}

	resp.Data = cakeResponses
	resp.Success = true
	resp.Message = "Cakes fetched successfully"
	return ctx.JSON(resp)
}

func (c *CakeController) GetCakeByID(ctx *fiber.Ctx) error {
	resp := new(utils.Response)

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		resp.Message = "Invalid cake ID"
		resp.Success = false
		return ctx.Status(fiber.StatusBadRequest).JSON(resp)
	}

	cake, err := c.cakeUseCase.GetCakeByID(id)
	if err != nil {
		c.logger.Errorf("Failed to get cake: %v", err)
		resp.Message = "Cake not found"
		resp.Success = false
		return ctx.Status(fiber.StatusNotFound).JSON(resp)
	}

	resp.Data = model.CakeToResponse(cake)
	resp.Message = "Cake fetched successfully"
	resp.Success = true
	return ctx.JSON(resp)
}

func (c *CakeController) CreateCake(ctx *fiber.Ctx) error {
	resp := new(utils.Response)

	var request model.CreateUpdateCakeRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("Failed to parse body: ", err)
		resp.Message = "Invalid request body"
		resp.Success = false
		return ctx.Status(http.StatusBadRequest).JSON(resp)
	}

	if err := c.validatePayload(request); err != nil {
		resp.Message = err.Error()
		resp.Success = false
		return ctx.Status(http.StatusBadRequest).JSON(resp)
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
		resp.Message = "Failed to create cake"
		resp.Success = false
		return ctx.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Success = true
	resp.Message = "Cake created successfully"
	resp.Data = model.CakeToResponse(cake)
	return ctx.Status(http.StatusOK).JSON(resp)
}

func (c *CakeController) UpdateCake(ctx *fiber.Ctx) error {
	resp := new(utils.Response)

	var request model.CreateUpdateCakeRequest
	if err := ctx.BodyParser(&request); err != nil {
		resp.Message = "Invalid request body"
		resp.Success = false
		return ctx.Status(http.StatusBadRequest).JSON(resp)
	}

	if err := c.validatePayload(request); err != nil {
		resp.Message = err.Error()
		resp.Success = false
		return ctx.Status(http.StatusBadRequest).JSON(resp)
	}

	cakeID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		resp.Message = "Invalid cake ID"
		resp.Success = false
		return ctx.Status(http.StatusBadRequest).JSON(resp)
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
		resp.Message = "Failed to update cake"
		resp.Success = false
		return ctx.Status(http.StatusInternalServerError).JSON(resp)
	}
	
	resp.Success = true
	resp.Message = "Cake updated successfully"
	resp.Data = model.CakeToResponse(cake)
	return ctx.Status(http.StatusOK).JSON(resp)
}

func (c *CakeController) DeleteCake(ctx *fiber.Ctx) error {
	resp := new(utils.Response)

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		resp.Message = "Invalid cake ID"
		resp.Success = false
		return ctx.Status(fiber.StatusBadRequest).JSON(resp)
	}

	err = c.cakeUseCase.DeleteCake(id)
	if err != nil {
		c.logger.Errorf("Failed to delete cake: %v", err)
		resp.Message = "Failed to delete cake"
		resp.Success = false
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp.Success = true
	resp.Message = "Cake deleted successfully"
	return ctx.JSON(resp)
}

func (c *CakeController) validatePayload(request model.CreateUpdateCakeRequest) error {
	if err := c.validator.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errMessages := make([]string, len(validationErrors))
		for i, e := range validationErrors {
			errMessages[i] = "Field '" + e.Field() + "' failed on '" + e.Tag() + "' rule"
		}
		return fiber.NewError(http.StatusBadRequest, "Validation failed: "+strings.Join(errMessages, ", "))
	}
	return nil
}
