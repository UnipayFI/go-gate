package options

import (
	"context"
	"strconv"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// AmendOptionsOrderService -- PUT /api/v4/options/orders/{order_id} (private)
//
// Modifies an existing options order in place, replacing its contract, price and
// size. All three fields are required by Gate.
type AmendOptionsOrderService struct {
	c       *OptionsClient
	orderID int64
	body    map[string]any
}

// NewAmendOptionsOrderService amends order orderID with a new contract, price
// and size (size is the signed number of contracts: positive buy, negative
// sell).
func (c *OptionsClient) NewAmendOptionsOrderService(orderID int64, contract string, price decimal.Decimal, size int64) *AmendOptionsOrderService {
	return &AmendOptionsOrderService{c: c, orderID: orderID, body: map[string]any{
		"contract": contract,
		"price":    price.String(),
		"size":     size,
	}}
}

func (s *AmendOptionsOrderService) Do(ctx context.Context) (*OptionsOrder, error) {
	req := request.Put(ctx, s.c, "/api/v4/options/orders/"+strconv.FormatInt(s.orderID, 10), s.body).WithSign()
	return request.Do[OptionsOrder](req)
}
