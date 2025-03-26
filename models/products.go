package model

import (
	"github.com/google/uuid"
)

type Shop struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GetProductReq struct {
	Id string `json:"id" validate:"uuid"`
}

type GetProductResp struct {
	Id           string  `json:"id"`
	ShopId       string  `json:"shop_id"`
	CategoryId   string  `json:"category_id"`
	ShopName     string  `json:"shop_name"`
	CategoryName string  `json:"category_name"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Stock        int64   `json:"stock"`
	ImageUrl     string  `json:"image_url"`
}

type CreateProductReq struct {
	UserId     string  `json:"user_id" validate:"uuid"`
	ShopId     string  `json:"shop_id" validate:"uuid"`
	CategoryId string  `json:"category_id" validate:"uuid"`
	Name       string  `json:"name" validate:"required"`
	Price      float64 `json:"price" validate:"required"`
	Stock      int64   `json:"stock" validate:"required"`
	ImageUrl   string  `json:"image_url" validate:"required"`
}

type DeleteProductReq struct {
	UserId string `json:"user_id" validate:"uuid"`
	Id     string `json:"id" validate:"uuid"`
}

type GetProductsReq struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (g *GetProductsReq) SetDefault() {
	if g.Page < 1 {
		g.Page = 1
	}

	if g.Limit < 1 {
		g.Limit = 10
	}
}

type GetProductsResp struct {
	Items []*ProductItem `json:"items"`
	Meta  *Meta          `json:"meta"`
}

type ProductItem struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Stock    int64   `json:"stock"`
	ImageUrl string  `json:"image_url"`
}
