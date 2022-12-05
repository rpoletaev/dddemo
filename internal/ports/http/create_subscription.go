package http

import (
	"subscription/gen/http/restapi/operations"
	"subscription/internal/domain"

	"github.com/go-openapi/runtime/middleware"
)

func (api *Api) CreateSubscription(params operations.CreateSubscriptionParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()
	cmd, err := api.uowFactory.SubscribeCmd(ctx, uint64(*params.Req.UserID))
	if err != nil {
		return operations.NewCreateSubscriptionInternalServerError()
	}

	if err := cmd.Do(ctx); err != nil {
		switch err {
		case domain.ErrAlreadyExists:
			return operations.NewCreateSubscriptionConflict()
		default:
			return operations.NewCreateSubscriptionInternalServerError()
		}
	}
	return operations.NewCreateSubscriptionCreated()
}
