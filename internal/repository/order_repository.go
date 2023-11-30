package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/shopspring/decimal"
)

type (
	OrderRepository interface {
		FindOrderBySellerID(ctx context.Context, sellerId int64, params *dto.OrderSellerParams) ([]dto.OrderSellerModel, error)
		FindOrderByBuyerID(ctx context.Context, buyerId int64, params *dto.OrderParams) ([]dto.OrderBuyerModel, error)
		FirstOrderByOrderID(ctx context.Context, orderId int64) (*model.Order, error)
		FindOrderBySellerIDMetadata(ctx context.Context, sellerId int64, params *dto.OrderSellerParams) ([]model.Order, error)
		FindOrderByBuyerIDMetadata(ctx context.Context, buyerId int64, params *dto.OrderParams) ([]model.Order, error)
		CreateOrder(ctx context.Context, accountId int64, priceList, delivery []decimal.Decimal, promotionAmount []float64, promotionName []string, payload dto.CreateOrderRequestPayload, transaction []*model.Transaction) error
		CreateOrderProductVariant(ctx context.Context, orderId int64, payload dto.CartOrderModel, totalPrice decimal.Decimal) error
		UpdateOrderStatus(ctx context.Context, orderId int64, status constant.OrderStatusType, eat *time.Time) error
		UpdateCancelOrder(ctx context.Context, orderId, accountId int64, transaction *model.Transaction) error
		UpdateReceiveOrder(ctx context.Context, orderId, buyerId, sellerId int64, transaction *model.Transaction) error
	}
	orderRepository struct {
		db *sqlx.DB
		tr TransactionRepository
	}
)

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderId int64, status constant.OrderStatusType, eat *time.Time) error {
	query := `
	UPDATE orders 
	SET status = $1, updated_at = $2`
	if eat != nil {
		query += ", estimated_time_arrival = '" + eat.Format("2006-01-02 15:04:05") + "'"
	}
	query += "\nWHERE id = $3"

	qs := `
	UPDATE promotions
	SET quota = quota-1, updated_at = $1
	WHERE id = (
		SELECT
			p.id
		FROM
			orders o
			LEFT JOIN promotions p ON o.promotion_name = p."name"
		WHERE o.id = $2

	)
	`
	qDetail := `
		SELECT * FROM order_details od WHERE od.order_id = $1;
	`

	orderDetails := make([]model.OrderDetail, 0)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, status, time.Now(), orderId)
	if err != nil {
		return err
	}

	if status == constant.ProcessOrderStatus {
		err := tx.Select(&orderDetails, qDetail, orderId)
		if err != nil {
			return err
		}
		for _, orderDetail := range orderDetails {
			variants := strings.Split(orderDetail.VariantName, "-")
			variantLength := len(variants)
			if variantLength == 1 {
				if variants[0] == "" {
					if err := DecreaseStock(tx, orderDetail.ProductCode, orderDetail.Quantity, constant.ProductVariantDefault, constant.ProductVariantDefault); err != nil {
						return err
					}
					continue
				}
				if err := DecreaseStock(tx, orderDetail.ProductCode, orderDetail.Quantity, variants[0], constant.ProductVariantDefault); err != nil {
					return err
				}
				continue
			}

			if err := DecreaseStock(tx, orderDetail.ProductCode, orderDetail.Quantity, variants[0], variants[1]); err != nil {
				return err
			}
		}
		_, err = tx.Exec(qs, time.Now(), orderId)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) FindOrderBySellerID(ctx context.Context, sellerId int64, params *dto.OrderSellerParams) ([]dto.OrderSellerModel, error) {
	orders := make([]dto.OrderSellerModel, 0)

	query := `
		SELECT 
			o.id,
			o.status,
			od.product_name,
			od.thumbnail_url,
			od.variant_name,
			od.sub_total_price,
			od.quantity,
			aa.receiver_name,
			aa.receiver_phone_number,
			aa.detail AS address_detail,
			c."name" AS courier_name,
			o.estimated_time_arrival,
			o.promotion_amount
		FROM (
			SELECT * FROM orders o ORDER BY o.updated_at DESC LIMIT $1 OFFSET $2 
		) AS o
		LEFT JOIN order_details od ON od.order_id = o.id
		LEFT JOIN couriers c ON c.id = o.courier_id 
		LEFT JOIN account_addresses aa ON aa.account_id = o.buyer_id
		WHERE o.seller_id = $3 AND o.deleted_at IS NULL `
	if params.Status != "" {
		query += "AND status = '" + string(params.Status) + "'"
	}
	query += "ORDER BY o.updated_at DESC"
	limit := constant.OrderSellerDefaultItems
	offset := (params.Page - 1) * limit
	err := r.db.SelectContext(ctx, &orders, query, limit, offset, sellerId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) FindOrderBySellerIDMetadata(ctx context.Context, sellerId int64, params *dto.OrderSellerParams) ([]model.Order, error) {
	orders := make([]model.Order, 0)

	query := `
		SELECT * FROM orders o
		WHERE o.seller_id = $1 AND o.deleted_at IS NULL `
	if params.Status != "" {
		query += "AND status = '" + string(params.Status) + "'"
	}
	query += "\nORDER BY o.updated_at DESC"

	err := r.db.SelectContext(ctx, &orders, query, sellerId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) FirstOrderByOrderID(ctx context.Context, orderId int64) (*model.Order, error) {
	order := new(model.Order)
	query := `SELECT * FROM orders o WHERE o.id = $1 AND o.deleted_at IS NULL`
	err := r.db.GetContext(ctx, order, query, orderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrOrderIDNotFount
		}
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) FindOrderByBuyerID(ctx context.Context, buyerId int64, params *dto.OrderParams) ([]dto.OrderBuyerModel, error) {
	orders := make([]dto.OrderBuyerModel, 0)

	qs := `
	SELECT 
		s.name AS shop_name,
		o.id,
		o.status,
		od.product_code,
		od.product_name,
		od.thumbnail_url,
		od.variant_name,
		od.sub_total_price,
		od.quantity,
		aa.receiver_name,
		aa.receiver_phone_number,
		aa.detail AS address_detail,
		c."name" AS courier_name,
		o.delivery_cost,
		o.estimated_time_arrival,
		o.promotion_amount
	FROM (
		SELECT * FROM orders o ORDER BY o.updated_at DESC LIMIT $1 OFFSET $2 
	) AS o
	LEFT JOIN order_details od ON od.order_id = o.id
	LEFT JOIN couriers c ON c.id = o.courier_id 
	LEFT JOIN account_addresses aa ON aa.account_id = o.buyer_id
	LEFT JOIN shops s ON s.account_id = o.seller_id
	WHERE o.buyer_id = $3 AND o.deleted_at IS NULL `
	if params.Status != "" {
		qs += "AND status = '" + string(params.Status) + "'"
	}
	qs += "ORDER BY o.updated_at DESC"
	limit := constant.OrderBuyerDefaultItems
	offset := (params.Page - 1) * limit

	err := r.db.SelectContext(ctx, &orders, qs, limit, offset, buyerId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) FindOrderByBuyerIDMetadata(ctx context.Context, buyerId int64, params *dto.OrderParams) ([]model.Order, error) {
	orders := make([]model.Order, 0)

	query := `
		SELECT * FROM orders o
		WHERE o.buyer_id = $1 AND o.deleted_at IS NULL `
	if params.Status != "" {
		query += "AND status = '" + string(params.Status) + "'"
	}
	query += "\nORDER BY o.updated_at DESC"

	err := r.db.SelectContext(ctx, &orders, query, buyerId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) CreateOrderProductVariant(ctx context.Context, orderId int64, payload dto.CartOrderModel, totalPrice decimal.Decimal) error {
	qs2 := `
	INSERT INTO order_details (
		order_id,
		product_variant_id,
		sub_total_price,
		quantity
		) VALUES (
			$1,
			$2,
			$3,
			$4
		)
	`

	_, err := r.db.ExecContext(ctx, qs2, orderId, payload.ProductVariantID, totalPrice, payload.Qty)
	if err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) CreateOrder(ctx context.Context, accountId int64, priceList, delivery []decimal.Decimal, promotionAmount []float64, promotionName []string, payload dto.CreateOrderRequestPayload, transaction []*model.Transaction) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qs1 := `
	SELECT 
		c.seller_id,
		s.id AS shop_id,
		s.name AS shop_name,
		c.id AS cart_id,
		p.product_code AS product_code,
		p."name" AS product_name,
		p.id AS product_id,
		p.thumbnail_url AS image_url,
		c.product_variant_id,
		pv.price AS base_price,
		pv.discount,
		c.quantity,
		pv.stock AS remaining_quantity,
		vt."name" AS variant1_name,
		vt2."name" AS variant2_name,
		c.is_checked
	FROM carts c 
	LEFT JOIN accounts a ON a.id = c.seller_id 
	LEFT JOIN shops s ON a.id = s.account_id 
	LEFT JOIN product_variants pv ON pv.id = c.product_variant_id
	LEFT JOIN products p ON p.id = pv.product_id
	LEFT JOIN variant_types vt ON vt.id = pv.variant_type1_id
	LEFT JOIN variant_types vt2 ON vt2.id = pv.variant_type2_id 
	WHERE c.account_id = $1 AND c.is_checked AND s.id = $2
	ORDER BY c.id
	`

	var courierId int64
	queryCourier := `
	SELECT
		sc.courier_id 
	FROM
		shop_couriers sc 
	WHERE
		sc.id = $1`

	var orderId int64
	qs2 := `
	INSERT INTO orders (
		status,
		courier_id,
		seller_id,
		buyer_id,
		delivery_cost,
		transaction_id,
		promotion_name,
		promotion_amount
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		) RETURNING (id)
	`

	qs3 := `
	INSERT INTO order_details (
		order_id,
		product_code,
		product_name,
		thumbnail_url,
		variant_name,
		sub_total_price,
		quantity
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)
	`

	qs4 := `
		DELETE FROM carts
		WHERE id = $1
	`

	for idx, order := range payload.Orders {
		transactionID, err := r.tr.CreateTransaction(tx, transaction[idx])
		if err != nil {
			return err
		}

		if err := TransferUserTemp(tx, accountId, transaction[idx].Amount); err != nil {
			return err
		}
		cartOrders := make([]dto.CreateOrderModel, 0)
		err = tx.Select(&cartOrders, qs1, accountId, order.ShopId)
		if err != nil {
			return err
		}

		promoDec := decimal.NewFromFloat(promotionAmount[idx])

		err = tx.QueryRowx(queryCourier, order.CourierId).Scan(&courierId)
		if err != nil {
			return err
		}

		err = tx.QueryRowx(qs2, constant.NewOrderStatus, courierId, cartOrders[0].SellerID, accountId, delivery[idx], transactionID, promotionName[idx], promoDec).Scan(&orderId)
		if err != nil {
			return err
		}
		for _, cart := range cartOrders {
			if cart.Variant1Name == constant.ProductVariantDefault {
				cart.Variant1Name = ""
			}
			if cart.Variant2Name == constant.ProductVariantDefault {
				cart.Variant2Name = ""
			}
			variantName := cart.Variant1Name + cart.Variant2Name
			if cart.Variant1Name != "" && cart.Variant2Name != "" {
				variantName = cart.Variant1Name + "-" + cart.Variant2Name
			}
			_, err := tx.Exec(qs3, orderId, cart.ProductCode, cart.ProductName, cart.ImageUrl, variantName, priceList[0], cart.Qty)
			if err != nil {
				return err
			}
			priceList = priceList[1:]
			_, err = tx.Exec(qs4, cart.CartID)
			if err != nil {
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) UpdateCancelOrder(ctx context.Context, orderId, accountId int64, transaction *model.Transaction) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = r.tr.CreateTransaction(tx, transaction)
	if err != nil {
		return err
	}

	if err := RefundTempUser(tx, accountId, transaction.Amount); err != nil {
		return err
	}

	query := `
		UPDATE orders 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	_, err = tx.Exec(query, constant.CancelOrderStatus, time.Now(), orderId)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) UpdateReceiveOrder(ctx context.Context, orderId, buyerId, sellerId int64, transaction *model.Transaction) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = r.tr.CreateTransaction(tx, transaction)
	if err != nil {
		return err
	}

	if err := TransferTempSeller(tx, buyerId, sellerId, transaction.Amount); err != nil {
		return err
	}

	query := `
		UPDATE orders 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	_, err = tx.Exec(query, constant.ReceiveOrderStatus, time.Now(), orderId)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func NewOrderRepository(db *sqlx.DB, tr TransactionRepository) OrderRepository {
	return &orderRepository{
		db: db,
		tr: tr,
	}
}
