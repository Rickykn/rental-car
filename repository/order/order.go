package order

import (
	"context"
	"database/sql"
	"github.com/Rickykn/rental-car/model"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, tx *sql.Tx, order *model.Order) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) IOrderRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) CreateOrder(ctx context.Context, tx *sql.Tx, order *model.Order) error {
	query := `INSERT INTO public.orders ( car_id, order_date, pickup_date, dropoff_date, pickup_location, dropoff_location) 
		VALUES( $1, $2, $3, $4, $5, $6);
	`
	_, err := tx.ExecContext(ctx, query, order.CarID, order.OrderDate, order.PickupDate, order.DropoffDate, order.PickupLocation, order.DropoffLocation)
	return err
}
