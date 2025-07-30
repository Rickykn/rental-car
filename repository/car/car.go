package car

import (
	"context"
	"database/sql"
	"github.com/Rickykn/rental-car/model"
	"github.com/google/uuid"
)

type ICarRepository interface {
	CreateCar(ctx context.Context, car model.Car) error
	GetCars(ctx context.Context) ([]model.Car, error)
	GetCarByID(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*model.Car, error)
	UpdateCarStatus(ctx context.Context, tx *sql.Tx, id uuid.UUID, status string) error
}

type carRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) ICarRepository {
	return &carRepository{db: db}
}

func (r *carRepository) CreateCar(ctx context.Context, car model.Car) error {
	query := `
		INSERT INTO public.cars (car_name, day_rate, month_rate, image) 
		VALUES( $1, $2, $3, $4)
		RETURNING id;
	`

	err := r.db.QueryRowContext(ctx, query,
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

func (r *carRepository) GetCarByID(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*model.Car, error) {
	query := `SELECT id, car_name, day_rate, month_rate, image, status FROM cars WHERE id = $1`
	row := tx.QueryRowContext(ctx, query, id)

	var car model.Car
	err := row.Scan(&car.ID, &car.CarName, &car.DayRate, &car.MonthRate, &car.Image, &car.Status)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func (r *carRepository) UpdateCarStatus(ctx context.Context, tx *sql.Tx, id uuid.UUID, status string) error {
	query := `UPDATE cars SET status = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, status, id)
	return err
}
