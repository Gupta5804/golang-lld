package main

import (
	"context"
	"fmt"

	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/gateways"
	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/pipeline"
)

func main() {
	vh := &pipeline.ValidateHandler{}
	fh := &pipeline.FraudHandler{}
	stripeAdaptor := &gateways.Stripe{}
	registry := map[string]gateways.PaymentGateway{
		"US": stripeAdaptor,
	}
	router := pipeline.NewRouterHandler(registry)
	vh.SetNext(fh).SetNext(router)
	request := &domain.PaymentRequest{
		Region:            "US",
		Amount:            100,
		ClientIP:          "192.1.1.2",
		DeviceFingerprint: "Safedevice",
	}

	if err := vh.Execute(context.Background(), request); err != nil {
        fmt.Printf("Payment failed: %v\n", err)
        return // Exit early!
    }
    
    // Happy path stays on the left margin
    fmt.Println("Payment executed successfully!")
}
