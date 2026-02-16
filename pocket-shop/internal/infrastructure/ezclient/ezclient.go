package ezclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	swagger "ez-client-go"

	"github.com/antihax/optional"

	"pocket-shop/config"
	"pocket-shop/internal/port"
)

type EzClient struct {
	client           *swagger.APIClient
	apiKey           string
	authToken        string
	sku              string
	retryMaxAttempts int
	retryBackoffSec  int
}

func NewEZClient(cfg *config.Config) port.EZClient {
	apiCfg := swagger.NewConfiguration()
	apiCfg.BasePath = cfg.EZBaseURL
	apiCfg.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	c := swagger.NewAPIClient(apiCfg)
	return &EzClient{
		client:           c,
		apiKey:           cfg.EZAPIKey,
		authToken:        cfg.EZAuthToken,
		sku:              cfg.EZSKU,
		retryMaxAttempts: cfg.EZRetryMaxAttempts,
		retryBackoffSec:  cfg.EZRetryBackoffSec,
	}
}

func (c *EzClient) authCtx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, swagger.ContextAPIKey, swagger.APIKey{Key: c.apiKey})
	ctx = context.WithValue(ctx, swagger.ContextAccessToken, c.authToken)
	return ctx
}

func (c *EzClient) CreateInstantOrder(ctx context.Context, orderID string) (*port.Order, error) {
	req := swagger.CreateInstantOrderRequest{
		Sku:                             c.sku,
		ClientOrderNumber:               orderID,
		EnableClientOrderNumberDupCheck: true,
	}
	var o *port.Order
	fn := func() (*http.Response, error) {
		res, httpResp, err := c.client.OrdersApi.CreateInstantOrder(c.authCtx(ctx), req)
		if err != nil {
			return httpResp, err
		}
		o = ToModelFromOrderResponse(res.Data, orderID)
		return httpResp, nil
	}
	_, err := withRetry(ctx, c.retryMaxAttempts, c.retryBackoffSec, fn)
	if err != nil {
		msg := err.Error()
		if gerr, ok := errors.AsType[swagger.GenericSwaggerError](err); ok {
			if b := gerr.Body(); len(b) > 0 {
				msg = msg + "; response: " + string(b)
			}
		}
		return nil, fmt.Errorf("EZ CreateInstantOrder: %s: %w", msg, err)
	}
	return o, nil
}

func (c *EzClient) GetOrder(ctx context.Context, orderID string) (*port.Order, error) {
	var o *port.Order
	opts := &swagger.OrdersApiGetOrdersOpts{
		ClientOrderNumber: optional.NewString(orderID),
	}
	fn := func() (*http.Response, error) {
		res, httpResp, err := c.client.OrdersApi.GetOrders(c.authCtx(ctx), opts)
		if err != nil {
			return httpResp, err
		}
		o = ToModelFromGetOrdersResponse(res.Data, orderID)
		return httpResp, nil
	}
	_, err := withRetry(ctx, c.retryMaxAttempts, c.retryBackoffSec, fn)
	if err != nil {
		msg := err.Error()
		if gerr, ok := errors.AsType[swagger.GenericSwaggerError](err); ok {
			if b := gerr.Body(); len(b) > 0 {
				msg = msg + "; response: " + string(b)
			}
		}
		return nil, fmt.Errorf("EZ GetOrders: %s: %w", msg, err)
	}
	return o, nil
}

func (c *EzClient) GetFirstRedeemCode(ctx context.Context, transactionID string) (string, error) {
	var code string
	fn := func() (*http.Response, error) {
		res, httpResp, err := c.client.CodesApi.GetCodes(c.authCtx(ctx), transactionID, nil)
		if err != nil {
			return httpResp, err
		}
		for _, p := range res.Data {
			for _, entry := range p.Codes {
				if entry.RedeemCode != "" {
					code = entry.RedeemCode
					return httpResp, nil
				}
			}
		}
		return httpResp, nil
	}
	_, err := withRetry(ctx, c.retryMaxAttempts, c.retryBackoffSec, fn)
	if err != nil {
		return "", fmt.Errorf("EZ GetCodes: %w", err)
	}
	return code, nil
}
