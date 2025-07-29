package order

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Rickykn/rental-car/model"
	"time"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, tx *sql.Tx, order *model.Order) error
	GetLastOrderNumberToday(ctx context.Context) (int, error)
	GetOrderByOrderCode(ctx context.Context, tx *sql.Tx, orderCode string) (*model.Order, error)
	UpdateCheckingData(ctx context.Context, tx *sql.Tx, orderCode string, checkinDate time.Time) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) IOrderRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) CreateOrder(ctx context.Context, tx *sql.Tx, order *model.Order) error {
	query := `INSERT INTO public.orders ( car_id, order_date, pickup_date, dropoff_date, pickup_location, dropoff_location,order_code) 
		VALUES( $1, $2, $3, $4, $5, $6,$7);
	`
	_, err := tx.ExecContext(ctx, query, order.CarID, order.OrderDate, order.PickupDate, order.DropoffDate, order.PickupLocation, order.DropoffLocation, order.OrderCode)
	return err
}

func (o *orderRepository) GetLastOrderNumberToday(ctx context.Context) (int, error) {
	query := `SELECT count(id) FROM orders WHERE order_date = current_date`

	var count int
	err := o.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (o *orderRepository) GetOrderByOrderCode(ctx context.Context, tx *sql.Tx, orderCode string) (*model.Order, error) {
	query := `SELECT id,car_id, order_code FROM orders WHERE order_code = $1`
	row := tx.QueryRowContext(ctx, query, orderCode)
	fmt.Println(row)

	var order model.Order
	err := row.Scan(&order.OrderID, &order.CarID, &order.OrderCode)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (o *orderRepository) UpdateCheckingData(ctx context.Context, tx *sql.Tx, orderCode string, checkinDate time.Time) error {
	query := `UPDATE orders SET checkin_date = $1 WHERE order_code = $2`
	_, err := tx.ExecContext(ctx, query, checkinDate, orderCode)
	return err
}
