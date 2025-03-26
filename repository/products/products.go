package products

import (
	"codebase-service/helper"
	model "codebase-service/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ ProductRepository = &store{}

type store struct {
	db    *sql.DB
	redis *redis.Client
}

func NewStore(db *sql.DB, redis *redis.Client) *store {
	return &store{
		db:    db,
		redis: redis,
	}
}

type ProductRepository interface {
	IsShopOwner(userId, shopId string) error
	CreateProduct(req *model.CreateProductReq) (*model.GetProductResp, error)
	GetProduct(req *model.GetProductReq) (*model.GetProductResp, error)
	GetProducts(req *model.GetProductsReq) (*model.GetProductsResp, error)
	DeleteProduct(req *model.DeleteProductReq) error
}

func (s *store) GetProduct(req *model.GetProductReq) (*model.GetProductResp, error) {
	resp, err := s.getProductInRedis(req)
	if err != nil {
		if err != redis.Nil {
			log.Printf("repo::GetProduct - failed to get product data from redis: %v", err)
			return nil, err
		}

		resp, err = s.getProductInDB(req)
		if err != nil {
			log.Printf("repo::GetProduct - failed to get product data from db: %v", err)
			return nil, err
		}

		err = s.setProductInRedis(req, resp)
		if err != nil {
			log.Printf("repo::GetProduct - failed to set product data in redis: %v", err)
			return nil, err
		}
	}

	return resp, nil
}

func (s *store) getProductInDB(req *model.GetProductReq) (*model.GetProductResp, error) {
	log.Printf("repo::getProductInDB - fetching product data from db")
	var (
		res  = new(model.GetProductResp)
		args = make([]interface{}, 0)
	)

	query := `
		SELECT
			p.id,
			p.shop_id,
			p.category_id,
			s.name AS shop_name,
			c.name AS category_name,
			p.name,
			p.price,
			p.stock,
			p.image_url
		FROM
			products p
		JOIN
			shops s ON p.shop_id = s.id
		JOIN
			product_categories c ON p.category_id = c.id
		WHERE
			p.id = ?
	`
	args = append(args, req.Id)

	query = helper.RebindQuery(query)

	row := s.db.QueryRow(query, args...)
	if err := row.Scan(
		&res.Id,
		&res.ShopId,
		&res.CategoryId,
		&res.ShopName,
		&res.CategoryName,
		&res.Name,
		&res.Price,
		&res.Stock,
		&res.ImageUrl,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("repo::GetProduct - no product found")
			return nil, fmt.Errorf("no product found")
		}
		log.Printf("repo::GetProduct - failed to fetch product data: %v", err)
		return nil, err
	}

	return res, nil
}

func (s *store) getProductInRedis(req *model.GetProductReq) (*model.GetProductResp, error) {
	log.Printf("repo::getProductInRedis - fetching product data from redis")
	var (
		res = new(model.GetProductResp)
		key = fmt.Sprintf("product:%s", req.Id)
	)
	log.Printf("repo::getProductInRedis - key: %s", key)

	data, err := s.redis.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("repo::getProductInRedis - failed to get product data from redis: %v", err)
		return nil, err
	}

	// unmarshal json
	if err := json.Unmarshal([]byte(data), res); err != nil {
		log.Printf("repo::getProductInRedis - failed to unmarshal product data: %v", err)
		return nil, err
	}

	return res, nil
}

func (s *store) setProductInRedis(req *model.GetProductReq, res *model.GetProductResp) error {
	log.Printf("repo::setProductInRedis - setting product data in redis")
	var (
		key        = fmt.Sprintf("product:%s", req.Id)
		expiration = time.Minute * 5
		ctx        = context.Background()
	)

	// to string json
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("repo::setProductInRedis - failed to marshal product data: %v", err)
		return err
	}

	if err := s.redis.Set(ctx, key, string(data), expiration).Err(); err != nil {
		log.Printf("repo::setProductInRedis - failed to set product data in redis: %v", err)
		return err
	}

	return nil
}

func (s *store) CreateProduct(req *model.CreateProductReq) (*model.GetProductResp, error) {
	var (
		res  = new(model.GetProductResp)
		args = make([]interface{}, 0)
	)

	query := `
		INSERT INTO
			products (shop_id, category_id, name, price, stock, image_url)
		VALUES
			(?, ?, ?, ?, ?, ?)
		RETURNING
			id, (SELECT name FROM shops WHERE id = ?) AS shop_name, (SELECT name FROM product_categories WHERE id = ?) AS category_name
	`
	args = append(
		args, req.ShopId, req.CategoryId, req.Name, req.Price, req.Stock, req.ImageUrl,
		req.ShopId, req.CategoryId,
	)

	query = helper.RebindQuery(query)

	row := s.db.QueryRow(query, args...)
	if err := row.Scan(
		&res.Id,
		&res.ShopName,
		&res.CategoryName,
	); err != nil {
		log.Printf("repo::CreateProduct - failed to scan product id: %v", err)
		return nil, err
	}

	res.ShopId = req.ShopId
	res.CategoryId = req.CategoryId
	res.Name = req.Name
	res.Price = req.Price
	res.Stock = req.Stock
	res.ImageUrl = req.ImageUrl

	return res, nil
}

func (s *store) IsShopOwner(userId, shopId string) error {
	var isShopOwner bool

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM
				shops
			WHERE
				user_id = ?
				AND id = ?
		)
	`
	query = helper.RebindQuery(query)

	row := s.db.QueryRow(query, userId, shopId)
	if err := row.Scan(&isShopOwner); err != nil {
		log.Printf("repo::IsShopOwner - failed to check if user is shop owner: %v", err)
		return err
	}

	if !isShopOwner {
		log.Printf("repo::IsShopOwner - user is not shop owner")
		return fmt.Errorf("user is not shop owner")
	}

	return nil
}

func (s *store) DeleteProduct(req *model.DeleteProductReq) error {
	query := `
		UPDATE products
		SET deleted_at = NOW()
		FROM shops
		WHERE
			products.id = ?
			AND products.shop_id = shops.id
			AND shops.user_id = ?
			AND products.deleted_at IS NULL
	`

	query = helper.RebindQuery(query)

	if _, err := s.db.Exec(query, req.Id, req.UserId); err != nil {
		log.Printf("repo::DeleteProduct - failed to delete product: %v", err)
		return err
	}

	return nil
}

func (s *store) GetProducts(req *model.GetProductsReq) (*model.GetProductsResp, error) {
	resp, err := s.getProductsInRedis(req)
	if err != nil {
		if err != redis.Nil {
			log.Printf("repo::GetProducts - failed to get products data from redis: %v", err)
			return nil, err
		}

		resp, err = s.getProductsInDB(req)
		if err != nil {
			log.Printf("repo::GetProducts - failed to get products data from db: %v", err)
			return nil, err
		}

		err = s.setProductsInRedis(req, resp)
		if err != nil {
			log.Printf("repo::GetProducts - failed to set products data in redis: %v", err)
			return nil, err
		}
	}

	return resp, nil
}

func (s *store) getProductsInDB(req *model.GetProductsReq) (*model.GetProductsResp, error) {
	log.Printf("repo::getProductsInDB - fetching products data from db")
	var (
		totalData int
		res       = new(model.GetProductsResp)
		args      = make([]interface{}, 0)
	)
	res.Items = make([]*model.ProductItem, 0)
	res.Meta = new(model.Meta)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			p.id,
			p.name,
			p.price,
			p.stock,
			p.image_url
		FROM
			products p
		WHERE
			p.deleted_at IS NULL
		LIMIT ? OFFSET ?
	`
	args = append(args, req.Limit, (req.Page-1)*req.Limit)

	query = helper.RebindQuery(query)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Printf("repo::getProductsInDB - failed to fetch products data: %v", err)
		return nil, err
	}

	for rows.Next() {
		var d model.ProductItem
		if err := rows.Scan(
			&totalData,
			&d.Id,
			&d.Name,
			&d.Price,
			&d.Stock,
			&d.ImageUrl,
		); err != nil {
			rows.Close()
			log.Printf("repo::getProductsInDB - failed to scan product data: %v", err)
			return nil, err
		}
		res.Items = append(res.Items, &d)
	}
	rows.Close()

	res.Meta.SetMeta(req.Page, req.Limit, totalData)

	return res, nil
}

func (s *store) getProductsInRedis(req *model.GetProductsReq) (*model.GetProductsResp, error) {
	log.Printf("repo::getProductsInRedis - fetching products data from redis")
	var (
		res = new(model.GetProductsResp)
		key = fmt.Sprintf("products:page:%d:limit:%d", req.Page, req.Limit)
	)
	log.Printf("repo::getProductsInRedis - key: %s", key)

	data, err := s.redis.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("repo::getProductsInRedis - failed to get products data from redis: %v", err)
		return nil, err
	}

	// unmarshal json
	if err := json.Unmarshal([]byte(data), res); err != nil {
		log.Printf("repo::getProductsInRedis - failed to unmarshal products data: %v", err)
		return nil, err
	}

	return res, nil
}

func (s *store) setProductsInRedis(req *model.GetProductsReq, res *model.GetProductsResp) error {
	log.Printf("repo::setProductsInRedis - setting products data in redis")
	var (
		key        = fmt.Sprintf("products:page:%d:limit:%d", req.Page, req.Limit)
		expiration = time.Minute * 5
		ctx        = context.Background()
	)

	// to string json
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("repo::setProductsInRedis - failed to marshal products data: %v", err)
		return err
	}

	if err := s.redis.Set(ctx, key, string(data), expiration).Err(); err != nil {
		log.Printf("repo::setProductsInRedis - failed to set products data in redis: %v", err)
		return err
	}

	return nil
}
