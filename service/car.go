package service

import (
	"context"
	"github.com/Rickykn/rental-car/model"
	"github.com/Rickykn/rental-car/repository"
)

type ICarService interface {
	CreateNewCar(ctx context.Context, car model.Car) error
	ShowAllCar(ctx context.Context) ([]model.Car, error)
}

type carService struct {
	repo repository.ICarRepository
}

func NewCarService(repo repository.ICarRepository) ICarService {
	return &carService{repo: repo}
}

func (c carService) CreateNewCar(ctx context.Context, car model.Car) error {
	newCar := model.Car{
		CarName:   car.CarName,
		DayRate:   car.DayRate,
		MonthRate: car.MonthRate,
		Image:     car.Image,
	}
	err := c.repo.CreateCar(ctx, newCar)
	if err != nil {
		return err
	}

	return nil
}

func (c carService) ShowAllCar(ctx context.Context) ([]model.Car, error) {
	cars, err := c.repo.GetCars(ctx)
	if err != nil {
		return nil, err
	}

	return cars, nil
}
