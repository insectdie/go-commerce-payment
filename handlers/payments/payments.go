package payments

import (
	"codebase-service/helper"
	midtranssvc "codebase-service/integration/midtrans"
	entity "codebase-service/integration/midtrans/entity"

	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	midtrans midtranssvc.MidtransContract
	v        *validator.Validate
}

func NewHandler(midtrans midtranssvc.MidtransContract, v *validator.Validate) *Handler {
	return &Handler{
		midtrans: midtrans,
		v:        v,
	}
}

func (h *Handler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var (
		req = new(entity.CreatePaymentRequest)
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	basicHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("SB-Mid-server-jz-9ZTjDo8yA-5kZCU6rgDNr"+":"))

	log.Printf("handler::CreatePayment - request: %+v", basicHeader)

	resp, err := h.midtrans.CreatePayment(r.Context(), req)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	statusCodeMap := map[string]int{
		"05":  http.StatusInternalServerError,
		"201": http.StatusCreated,
		"200": http.StatusOK,
		"400": http.StatusBadRequest,
		"401": http.StatusUnauthorized,
		"406": http.StatusNotAcceptable,
		"500": http.StatusInternalServerError,
		"503": http.StatusServiceUnavailable,
		"900": http.StatusInternalServerError,
	}

	statusCode, ok := statusCodeMap[resp.StatusCode]
	if !ok {
		log.Printf("handler::CreatePayment - failed to map status code, status code: %s", resp.StatusCode)
		statusCode = http.StatusInternalServerError
	}

	helper.HandleResponse(w, statusCode, resp.StatusMessage, resp)
}
