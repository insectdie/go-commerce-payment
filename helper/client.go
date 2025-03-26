package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

type NetClientRequest struct {
	NetClient  *http.Client
	RequestUrl string
	QueryParam []QueryParams
}

type QueryParams struct {
	Param string
	Value string
}

type Response struct {
	Res        []byte
	Err        error
	StatusCode int
}

var DefaultNetClient = &http.Client{
	Timeout: time.Second * 10,
}

func NewNetClientRequest(url string, client *http.Client) *NetClientRequest {
	if client == nil {
		client = DefaultNetClient
	}
	return &NetClientRequest{
		NetClient:  client,
		RequestUrl: url,
	}
}

func (ncr *NetClientRequest) AddQueryParam(param, value string) {
	ncr.QueryParam = append(ncr.QueryParam, QueryParams{Param: param, Value: value})
}

func (ncr *NetClientRequest) buildUrl() (string, error) {
	urlObj, err := url.Parse(ncr.RequestUrl)
	if err != nil {
		return "", err
	}

	if len(ncr.QueryParam) > 0 {
		query := urlObj.Query()
		for _, param := range ncr.QueryParam {
			query.Add(param.Param, param.Value)
		}
		urlObj.RawQuery = query.Encode()
	}
	return urlObj.String(), nil
}

func (ncr *NetClientRequest) sendRequest(method string, load interface{}, channel chan Response) {
	go func() {
		marshalled, err := json.Marshal(load)
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		urlString, err := ncr.buildUrl()
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		req, err := http.NewRequest(method, urlString, bytes.NewBuffer(marshalled))
		if err != nil {
			channel <- Response{Err: err}
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := ncr.NetClient.Do(req)
		if err != nil {
			channel <- Response{Err: err}
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		channel <- Response{Res: respBody, StatusCode: resp.StatusCode}
	}()
}

func (ncr *NetClientRequest) Get(load interface{}, channel chan Response) {
	ncr.sendRequest(http.MethodGet, load, channel)
}

func (ncr *NetClientRequest) Post(load interface{}, channel chan Response) {
	ncr.sendRequest(http.MethodPost, load, channel)
}

func (ncr *NetClientRequest) Patch(load interface{}, channel chan Response) {
	ncr.sendRequest(http.MethodPatch, load, channel)
}
