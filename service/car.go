package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Rickykn/rental-car/model"
	"github.com/Rickykn/rental-car/repository/car"
	"github.com/Rickykn/rental-car/repository/order"
	"time"
)

type ICarService interface {
	CreateNewCar(ctx context.Context, car model.Car) error
	ShowAllCar(ctx context.Context) ([]model.Car, error)
	BookCar(ctx context.Context, order model.BookReq) error
}

type carService struct {
	cr car.ICarRepository
	or order.IOrderRepository
	db *sql.DB
}

func NewCarService(repo car.ICarRepository, or order.IOrderRepository, db *sql.DB) ICarService {
	return &carService{cr: repo, or: or, db: db}
}

var (
	ErrCarNotAvailable = errors.New("car is not available for booking")
	ErrCarNotFound     = errors.New("car not found")
)

func (c carService) CreateNewCar(ctx context.Context, car model.Car) error {
	newCar := model.Car{
		CarName:   car.CarName,
		DayRate:   car.DayRate,
		MonthRate: car.MonthRate,
		Image:     car.Image,
	}
	err := c.cr.CreateCar(ctx, newCar)
	if err != nil {
		return err
	}

	return nil
}

func (c carService) ShowAllCar(ctx context.Context) ([]model.Car, error) {
	cars, err := c.cr.GetCars(ctx)
	if err != nil {
		return nil, err
	}

	return cars, nil
}

func (s *carService) BookCar(ctx context.Context, req model.BookReq) error {
	pickupDate, err := time.Parse("2006-01-02", req.PickupDate)
	if err != nil {
		return errors.New("invalid pickup_date format, expected YYYY-MM-DD")
	}
	dropoffDate, err := time.Parse("2006-01-02", req.DropoffDate)
	if err != nil {
		return errors.New("invalid dropoff_date format, expected YYYY-MM-DD")
	}
	if dropoffDate.Before(pickupDate) {
		return errors.New("dropoff_date must be after pickup_date")
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	car, err := s.cr.GetCarByID(ctx, tx, req.CarID)
	if err != nil {
		return ErrCarNotFound
	}

	if car.Status != "available" {
		return ErrCarNotAvailable
	}

	err = s.cr.UpdateCarStatus(ctx, tx, car.ID, "book")
	if err != nil {
		return err
	}

	newOrder := &model.Order{
		CarID:           req.CarID,
		OrderDate:       time.Now(),
		PickupDate:      pickupDate,
		DropoffDate:     dropoffDate,
		PickupLocation:  req.PickupLocation,
		DropoffLocation: req.DropoffLocation,
	}

	err = s.or.CreateOrder(ctx, tx, newOrder)
	if err != nil {
		return err
	}

	return tx.Commit()
}
