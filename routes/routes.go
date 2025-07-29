package routes

import (
	"database/sql"
	"github.com/Rickykn/rental-car/handler"
	"github.com/Rickykn/rental-car/repository/car"
	"github.com/Rickykn/rental-car/repository/order"
	"github.com/Rickykn/rental-car/service"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/v1")

	userRepository := car.NewUserRepository(db)
	orderRepository := order.NewOrderRepository(db)
	carService := service.NewCarService(userRepository, orderRepository, db)
	carHandler := handler.NewCarHandler(carService)

	api.Post("/register_car", carHandler.RegisterCar)
	api.Get("/car", carHandler.ShowCars)
	api.Post("/book_car", carHandler.BookCar)
}
