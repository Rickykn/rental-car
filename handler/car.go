package handler

import (
	"context"
	"fmt"
	"github.com/Rickykn/rental-car/model"
	"github.com/Rickykn/rental-car/service"
	"github.com/gofiber/fiber/v2"
)

type CarHandler struct {
	svc service.ICarService
}

func NewCarHandler(svc service.ICarService) *CarHandler {
	return &CarHandler{
		svc: svc,
	}
}

func (ch CarHandler) RegisterCar(c *fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)
	if !ok || ctx == nil {
		ctx = context.Background() // fallback if not set
	}

	var req model.Car

	err := c.BodyParser(&req)
	if err != nil {
		return err
	}

	err = ch.svc.CreateNewCar(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	// Success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "car created successfully",
	})
}

func (ch CarHandler) ShowCars(c *fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)
	if !ok || ctx == nil {
		ctx = context.Background() // fallback if not set
	}

	car, err := ch.svc.ShowAllCar(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "success get all car",
		"data":    car,
	})
}

func (ch CarHandler) BookCar(c *fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)
	if !ok || ctx == nil {
		ctx = context.Background() // fallback if not set
	}
	var req model.BookReq
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	fmt.Println(req)
	err = ch.svc.BookCar(ctx, req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "success book car",
	})
}
