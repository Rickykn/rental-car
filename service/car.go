package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Rickykn/rental-car/model"
	"github.com/Rickykn/rental-car/repository/car"
	"github.com/Rickykn/rental-car/repository/order"
	"time"
)

type ICarService interface {
	CreateNewCar(ctx context.Context, car model.Car) error
	ShowAllCar(ctx context.Context) ([]model.Car, error)
	BookCar(ctx context.Context, order model.BookReq) error
	Checkin(ctx context.Context, checkingReq model.CheckinReq) error
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
	ErrOrderNotFound   = errors.New("order not found")
)

func (c carService) GenerateOrderNumber(ctx context.Context) (string, error) {
	fmt.Println("MASUK SINI")
	prefix := "CR"
	dateStr := time.Now().Format("020106")

	countOrder, err := c.or.GetLastOrderNumberToday(ctx)
	if err != nil {
		return "", err
	}
	sequence := countOrder + 1
	orderNumber := fmt.Sprintf("%s%s%d", prefix, dateStr, sequence)

	return orderNumber, nil
}

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

func (c *carService) BookCar(ctx context.Context, req model.BookReq) error {
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

	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	car, err := c.cr.GetCarByID(ctx, tx, req.CarID)
	if err != nil {
		return ErrCarNotFound
	}

	if car.Status != "available" {
		return ErrCarNotAvailable
	}

	err = c.cr.UpdateCarStatus(ctx, tx, car.ID, "book")
	if err != nil {
		return err
	}
	orderCode, err := c.GenerateOrderNumber(ctx)
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
		OrderCode:       orderCode,
	}

	err = c.or.CreateOrder(ctx, tx, newOrder)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (c *carService) Checkin(ctx context.Context, checkingReq model.CheckinReq) error {
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	order, err := c.or.GetOrderByOrderCode(ctx, tx, checkingReq.OrderCode)

	if err != nil {
		return ErrOrderNotFound
	}

	// Specific time
	checkinTime := time.Date(2025, 7, 30, 14, 30, 0, 0, time.UTC)
	err = c.or.UpdateCheckingData(ctx, tx, order.OrderCode, checkinTime)
	if err != nil {
		return err
	}

	err = c.cr.UpdateCarStatus(ctx, tx, order.CarID, "rent")
	if err != nil {
		return err
	}

	return tx.Commit()
}
