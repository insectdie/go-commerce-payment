package products

import (
	"codebase-service/helper"
	model "codebase-service/models"
	"codebase-service/usecases/products"
	"codebase-service/util/middleware"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Svc products.ProductSvc
	v   *validator.Validate
}

func NewHandler(Svc products.ProductSvc, v *validator.Validate) *Handler {
	return &Handler{
		Svc: Svc,
		v:   v,
	}
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var req = new(model.GetProductReq)
	req.Id = r.PathValue("id")

	if err := h.v.Struct(req); err != nil {
		log.Printf("handler::GetProduct - failed to validate request, err: %v", err)
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	bRes, err := h.Svc.GetProduct(req)
	if err != nil {
		if err.Error() == "no product found" {
			helper.HandleResponse(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, bRes)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req = new(model.CreateProductReq)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	req.UserId = middleware.GetUserID(r.Context())

	if err := h.v.Struct(req); err != nil {
		log.Printf("handler::CreateProduct - failed to validate request, err: %v", err)
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	bRes, err := h.Svc.CreateProduct(req)
	if err != nil {
		if err.Error() == "user is not shop owner" {
			helper.HandleResponse(w, http.StatusForbidden, err.Error(), nil)
			return
		}

		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusCreated, helper.SUCCESS_MESSSAGE, bRes)
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var req = new(model.DeleteProductReq)
	req.Id = r.PathValue("id")
	req.UserId = middleware.GetUserID(r.Context())

	if err := h.v.Struct(req); err != nil {
		log.Printf("handler::DeleteProduct - failed to validate request, err: %v", err)
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.Svc.DeleteProduct(req)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, nil)
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	var req = new(model.GetProductsReq)
	var (
		page, _  = strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	)

	req.Page = page
	req.Limit = limit

	req.SetDefault()

	bRes, err := h.Svc.GetProducts(req)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.HandleResponse(w, http.StatusOK, helper.SUCCESS_MESSSAGE, bRes)
}
