package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"worktime-service/internal/gateway/dto/response"
	"worktime-service/pkg/constants"

	"github.com/hashicorp/consul/api"
)

type TermGateway interface {
	GetTermByID(ctx context.Context, id string) (*response.TermResponse, error)
}

type termGatewayImpl struct {
	serviceName string
	consul      *api.Client
}

func NewUserGateway(serviceName string, consulClient *api.Client) TermGateway {
	return &termGatewayImpl{
		serviceName: serviceName,
		consul:      consulClient,
	}
}

func (g *termGatewayImpl) GetTermByID(ctx context.Context, termID string) (*response.TermResponse, error) {
	token, ok := ctx.Value(constants.TokenKey).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	resp, err := client.Call("GET", "/api/v1/gateway/terms/"+termID, nil)
	if err != nil {
		return nil, fmt.Errorf("call API term fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[response.TermResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}
