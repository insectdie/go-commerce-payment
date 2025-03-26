package products

import (
	model "codebase-service/models"
	"codebase-service/repository/products"
)

var _ ProductSvc = &svc{}

type svc struct {
	store products.ProductRepository
}

func NewProductSvc(store products.ProductRepository) *svc {
	return &svc{
		store: store,
	}
}

type ProductSvc interface {
	GetProduct(req *model.GetProductReq) (*model.GetProductResp, error)
	GetProducts(req *model.GetProductsReq) (*model.GetProductsResp, error)
	CreateProduct(req *model.CreateProductReq) (*model.GetProductResp, error)
	DeleteProduct(req *model.DeleteProductReq) error
}

func (s *svc) GetProduct(req *model.GetProductReq) (*model.GetProductResp, error) {
	return s.store.GetProduct(req)
}

func (s *svc) CreateProduct(req *model.CreateProductReq) (*model.GetProductResp, error) {
	err := s.store.IsShopOwner(req.UserId, req.ShopId)
	if err != nil {
		return nil, err
	}

	return s.store.CreateProduct(req)
}

func (s *svc) GetProducts(req *model.GetProductsReq) (*model.GetProductsResp, error) {
	return s.store.GetProducts(req)
}

func (s *svc) DeleteProduct(req *model.DeleteProductReq) error {
	return s.store.DeleteProduct(req)
}
