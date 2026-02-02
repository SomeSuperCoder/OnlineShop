package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/internal/middleware"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductHandler struct {
	Repo *repository.Queries
	Pool *pgxpool.Pool
}
type GetAllProductsResponse struct {
	Body *[]repository.Product
}

func (h *ProductHandler) GetAll(ctx context.Context, input *Pagination) (*GetAllProductsResponse, error) {
	resp := new(GetAllProductsResponse)
	res, err := h.Repo.FindProductsPaged(ctx, repository.FindProductsPagedParams{
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	resp.Body = &res
	return resp, err
}

type SearchForProductsRequest struct {
	Query string `query:"query"`
}
type SearchForProductsResponse struct {
	Body struct {
		Products []repository.SearchForProductsRow `json:"products"`
	}
}

func (h *ProductHandler) SearchForProducts(ctx context.Context, input *SearchForProductsRequest) (*SearchForProductsResponse, error) {
	resp := new(SearchForProductsResponse)
	res, err := h.Repo.SearchForProducts(ctx, repository.SearchForProductsParams{
		ToTsquery: input.Query,
	})
	resp.Body.Products = res
	return resp, err
}

type CreateProductRequest struct {
	Body repository.InsertProductParams
}
type CreateProductResponse struct {
	Body repository.Product
}

func (h *ProductHandler) Post(ctx context.Context, input *CreateProductRequest) (*CreateProductResponse, error) {
	resp := new(CreateProductResponse)

	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) (repository.Product, error) {
		return q.InsertProduct(ctx, input.Body)
	})
	resp.Body = *res
	return resp, err
}

type GetProductByIDRequest struct {
	ID uuid.UUID `path:"id"`
}
type GetProductByIDResponse struct {
	Body *repository.Product
}

func (h *ProductHandler) GetByID(ctx context.Context, input *GetProductByIDRequest) (*GetProductByIDResponse, error) {
	resp := new(GetProductByIDResponse)
	res, err := h.Repo.GetProductByID(ctx, repository.GetProductByIDParams(*input))
	resp.Body = &res
	return resp, err
}

type DeleteProductRequest struct {
	ID uuid.UUID `path:"id"`
}
type DeleteProductResponse struct {
	Body *repository.Product
}

func (h *ProductHandler) Delete(ctx context.Context, input *DeleteProductRequest) (*DeleteProductResponse, error) {
	resp := new(DeleteProductResponse)
	res, err := h.Repo.DeleteProduct(ctx, repository.DeleteProductParams(*input))
	resp.Body = &res
	return resp, err
}

type UpdateProductRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Name    string `json:"name"`
		Details string `json:"details"`
		Price   int32  `json:"price"`
	}
}
type UpdateProductResponse struct {
	Body repository.Product
}

func (h *ProductHandler) Patch(ctx context.Context, input *UpdateProductRequest) (*UpdateProductResponse, error) {
	resp := new(UpdateProductResponse)

	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) (repository.Product, error) {
		return q.UpdateProduct(ctx, repository.UpdateProductParams{
			ID:      input.ID,
			Name:    input.Body.Name,
			Details: input.Body.Details,
			Price:   input.Body.Price,
		})
	})

	resp.Body = *res
	return resp, err
}
