package ordersRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/IzePhanthakarn/kawaii-shop/modules/orders"
	"github.com/IzePhanthakarn/kawaii-shop/modules/orders/ordersPatterns"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
	FindOneOrder(orderId string) (*orders.Order, error)
	FindOrder(req *orders.OrderFilter) ([]*orders.Order, int)
	InsertOrder(req *orders.Order) (string, error)
	UpdateOrder(req *orders.Order) error
}

type ordersRepository struct {
	db *sqlx.DB
}

func OrdersRepository(db *sqlx.DB) IOrdersRepository {
	return &ordersRepository{
		db: db,
	}
}

func (r *ordersRepository) FindOneOrder(orderId string) (*orders.Order, error) {
	query := `
		SELECT
			to_jsonb("t")
		FROM (
			SELECT
				"o"."id",
				"o"."user_id",
				"o"."transfer_slip",
				(
					SELECT
						array_to_json(array_agg("pt"))
					FROM (
						SELECT
							"spo"."id",
							"spo"."qty",
							"spo"."product"
						FROM "products_orders" "spo"
						WHERE "spo"."order_id" = "o"."id"
					) AS "pt"
				) AS "products",
				"o"."address",
				"o"."contact",
				"o"."status",
				(
					SELECT
						SUM(COALESCE(("po"."product"->>'price')::FLOAT * ("po"."qty")::FLOAT))
					FROM "products_orders" "po"
					WHERE "po"."order_id" = "o"."id"
				) AS "total_paid",
				"o"."created_at",
				"o"."updated_at"
			FROM "orders" "o"
			WHERE "o"."id" = $1
		) AS "t";
	`

	orderDate := &orders.Order{
		Products: make([]*orders.ProductsOrder, 0),
	}

	raw := make([]byte, 0)
	if err := r.db.Get(&raw, query, orderId); err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if err := json.Unmarshal(raw, &orderDate); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return orderDate, nil
}

func (r *ordersRepository) FindOrder(req *orders.OrderFilter) ([]*orders.Order, int) {
	builder := ordersPatterns.FindOrderBuilder(r.db, req)
	engineer := ordersPatterns.FindOrderEngineer(builder)

	result := engineer.FindOrder()
	count := engineer.CountOrder()

	return result, count
}

func (r *ordersRepository) InsertOrder(req *orders.Order) (string, error) {
	builder := ordersPatterns.InsertOrderBuilder(r.db, req)
	orderId, err := ordersPatterns.InsertOrderEngineer(builder).InsertOrder()
	if err != nil {
		return "", err
	}

	return orderId, nil
}

func (r *ordersRepository) UpdateOrder(req *orders.Order) error {
	query := `
		UPDATE "orders" SET
	`

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	if req.Status != "" {
		values = append(
			values,
			req.Status,
		)

		queryWhereStack = append(
			queryWhereStack,
			fmt.Sprintf(`"status" = $%d?`, lastIndex),
		)

		lastIndex++
	}

	if req.TransferSlip != nil {
		values = append(
			values,
			req.TransferSlip,
		)

		queryWhereStack = append(
			queryWhereStack,
			fmt.Sprintf(`"transfer_slip" = $%d?`, lastIndex),
		)

		lastIndex++
	}

	values = append(values, req.Id)

	queryClose := fmt.Sprintf(` WHERE "id" = $%d`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1)
		}
	}
	query += queryClose

	fmt.Println(query)
	if _, err := r.db.ExecContext(context.Background(), query, values...); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}
