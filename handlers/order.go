package handlers

import (
	"context"

	"github.com/SomeSuperCoder/HumaExampleProject/repository"
)

type OrderHandler struct {
	Repo *repository.Queries
}

type GetAllOrdersResponse struct {
	Body *[]repository.Order
}

func (h *OrderHandler) GetAll(ctx context.Context, input *struct{}) (*GetAllOrdersResponse, error) {
	resp := new(GetAllOrdersResponse)
	res, err := h.Repo.FindAllOrders(ctx)
	resp.Body = &res
	return resp, err
}

type CreateOrderRequest struct {
	Body struct {
		Details string `json:"name" example:"Bought an iPhone"`
	}
}
type CreateOrderResponse struct {
	Body *repository.Order
}

func (h *OrderHandler) Post(ctx context.Context, input *CreateOrderRequest) (*CreateOrderResponse, error) {
	resp := new(CreateOrderResponse)
	res, err := h.Repo.InsertOrder(ctx, repository.InsertOrderParams{
		Details: input.Body.Details,
	})
	resp.Body = &res
	return resp, err
}
