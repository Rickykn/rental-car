package routes

import (
	"database/sql"
	"github.com/Rickykn/rental-car/handler"
	"github.com/Rickykn/rental-car/repository"
	"github.com/Rickykn/rental-car/service"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/v1")

	userRepository := repository.NewUserRepository(db)
	carService := service.NewCarService(userRepository)
	carHandler := handler.NewCarHandler(carService)

	api.Post("/register_car", carHandler.RegisterCar)
	api.Get("/car", carHandler.ShowCars)
}
