package httptransport

import (
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/EgorTarasov/shopx/internal/domain"
	api "github.com/EgorTarasov/shopx/internal/generated/api"
	"github.com/EgorTarasov/shopx/internal/usecase/shop"
)

func openapiEmail(s string) openapi_types.Email { return openapi_types.Email(s) }

func toCreateOrderInput(userID domain.UserID, req api.CreateOrderRequest) shop.CreateOrderInput {
	items := make([]shop.OrderItemInput, len(req.Items))
	for i, it := range req.Items {
		items[i] = shop.OrderItemInput{SKU: domain.SKU(it.Sku), Qty: it.Qty}
	}
	return shop.CreateOrderInput{UserID: userID, Items: items}
}

func fromGetOrderOutput(out shop.GetOrderOutput) api.Order {
	items := make([]api.OrderItem, len(out.Items))
	for i, it := range out.Items {
		items[i] = api.OrderItem{Sku: it.SKU.String(), Qty: it.Qty}
	}
	return api.Order{
		Id:     out.ID.String(),
		UserId: out.UserID.String(),
		Status: api.OrderStatus(out.Status),
		Items:  items,
	}
}
