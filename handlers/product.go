package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/internal/middleware"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/danielgtaylor/huma/v2"
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

func (h *ProductHandler) GetPaged(ctx context.Context, input *Pagination) (*GetAllProductsResponse, error) {
	resp := new(GetAllProductsResponse)
	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) ([]repository.Product, error) {
		return q.FindProductsPaged(ctx, repository.FindProductsPagedParams{
			Limit:  input.Limit,
			Offset: input.Offset,
		})
	})
	resp.Body = res
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
	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) ([]repository.SearchForProductsRow, error) {
		return q.SearchForProducts(ctx, repository.SearchForProductsParams{
			ToTsquery: input.Query,
		})
	})
	resp.Body.Products = *res
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
	Body repository.Product
}

func (h *ProductHandler) GetByID(ctx context.Context, input *GetProductByIDRequest) (*GetProductByIDResponse, error) {
	resp := new(GetProductByIDResponse)
	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) (repository.Product, error) {
		return h.Repo.GetProductByID(ctx, repository.GetProductByIDParams(*input))
	})
	resp.Body = *res
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
	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) (repository.Product, error) {
		if is, err := q.IsProductOwner(ctx, repository.IsProductOwnerParams{
			ID: input.ID,
		}); err != nil {
			return repository.Product{}, err
		} else if !is {
			return repository.Product{}, huma.Error401Unauthorized("Access denied")
		}
		return q.DeleteProduct(ctx, repository.DeleteProductParams{
			ID: input.ID,
		})
	})
	resp.Body = res
	return resp, err
}

type UpdateProductRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Name    *string `json:"name,omitempty"`
		Details *string `json:"details,omitempty"`
		Price   *int32  `json:"price,omitempty"`
	}
}
type UpdateProductResponse struct {
	Body repository.Product
}

func (h *ProductHandler) Patch(ctx context.Context, input *UpdateProductRequest) (*UpdateProductResponse, error) {
	resp := new(UpdateProductResponse)
	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) (repository.Product, error) {
		if is, err := q.IsProductOwner(ctx, repository.IsProductOwnerParams{
			ID: input.ID,
		}); err != nil {
			return repository.Product{}, err
		} else if !is {
			return repository.Product{}, huma.Error401Unauthorized("Access denied")
		}

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
