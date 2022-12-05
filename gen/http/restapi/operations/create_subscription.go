// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// CreateSubscriptionHandlerFunc turns a function with the right signature into a create subscription handler
type CreateSubscriptionHandlerFunc func(CreateSubscriptionParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateSubscriptionHandlerFunc) Handle(params CreateSubscriptionParams) middleware.Responder {
	return fn(params)
}

// CreateSubscriptionHandler interface for that can handle valid create subscription params
type CreateSubscriptionHandler interface {
	Handle(CreateSubscriptionParams) middleware.Responder
}

// NewCreateSubscription creates a new http.Handler for the create subscription operation
func NewCreateSubscription(ctx *middleware.Context, handler CreateSubscriptionHandler) *CreateSubscription {
	return &CreateSubscription{Context: ctx, Handler: handler}
}

/* CreateSubscription swagger:route POST /subscriptions createSubscription

Create new subscription for user. A user can have only one subscription


*/
type CreateSubscription struct {
	Context *middleware.Context
	Handler CreateSubscriptionHandler
}

func (o *CreateSubscription) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewCreateSubscriptionParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
