package midtranssvc

import (
	"bytes"
	"codebase-service/config"
	entity "codebase-service/integration/midtrans/entity"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/sony/gobreaker"
)

var _ MidtransContract = &midtrans{}

type MidtransContract interface {
	CreatePayment(ctx context.Context, req *entity.CreatePaymentRequest) (entity.CreatePaymentResponse, error)
}

type midtrans struct {
	cfg *config.Config
	cb  *gobreaker.CircuitBreaker
}

func NewMidtransContract(
	cfg *config.Config,
	cb *gobreaker.CircuitBreaker,
) *midtrans {
	return &midtrans{
		cfg: cfg,
		cb:  cb,
	}
}

func (m *midtrans) CreatePayment(ctx context.Context, req *entity.CreatePaymentRequest) (entity.CreatePaymentResponse, error) {
	var (
		response  entity.CreatePaymentResponse
		ChargeURL string = m.cfg.MidtransChargeURL
	)

	respAny, err := m.cb.Execute(func() (interface{}, error) {
		bytesReq, err := json.Marshal(req)
		if err != nil {
			log.Printf("midtrans: failed to marshal request, err: %v", err)
			return response, err
		}

		request, err := http.NewRequest(http.MethodPost, ChargeURL, bytes.NewBuffer(bytesReq))
		if err != nil {
			log.Println("midtrans: failed to create request")
			return response, err
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", req.BasicAuthHeader)

		// create http client
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Println("midtrans: failed to do request")
			return response, err
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("midtrans: failed to read response body")
			return response, err
		}

		if err := json.Unmarshal(respBody, &response); err != nil {
			log.Println("midtrans: failed to unmarshal response")
			return response, err
		}

		return response, nil
	})
	if err != nil {
		return response, err
	}

	return respAny.(entity.CreatePaymentResponse), nil
}
