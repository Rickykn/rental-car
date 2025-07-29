package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Rickykn/rental-car/model"
)

type ICarRepository interface {
	CreateCar(ctx context.Context, car model.Car) error
	GetCars(ctx context.Context) ([]model.Car, error)
}

type carRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) ICarRepository {
	return &carRepository{db: db}
}

func (c carRepository) CreateCar(ctx context.Context, car model.Car) error {
	fmt.Println(car)
	query := `
		INSERT INTO public.cars (car_name, day_rate, month_rate, image) 
		VALUES( $1, $2, $3, $4)
		RETURNING id;
	`

	err := c.db.QueryRowContext(ctx, query,
		car.CarName,
		car.DayRate,
		car.MonthRate,
		car.Image,
	).Scan(&car.ID)

	return err
}

func (r *carRepository) GetCars(ctx context.Context) ([]model.Car, error) {
	query := `SELECT id, car_name, day_rate, month_rate, image FROM cars`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []model.Car

	for rows.Next() {
		var car model.Car
		err := rows.Scan(
			&car.ID,
			&car.CarName,
			&car.DayRate,
			&car.MonthRate,
			&car.Image,
		)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}
