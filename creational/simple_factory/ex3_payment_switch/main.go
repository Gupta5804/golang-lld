package main
import (
	"fmt"
)

type PaymentGateway interface {
	Pay(amount float64) error
}

type StripeGateway struct {
	stripeID int
}
func (s *StripeGateway) Pay(amount float64) error{
	fmt.Printf("Sending amount %f, to stripeID : %d",amount, s.stripeID )	
	return nil
}

type BitcoinGateway struct {
	bitcoinWallet int
}
func (b *BitcoinGateway) Pay(amount float64) error{
	fmt.Printf("Sending amount %f, to bitcoinWallet: %d", amount, b.bitcoinWallet)
	return nil
}

type CheckoutService struct {
	gateway PaymentGateway
}

func (c *CheckoutService) ProcessPayment(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("Invalid Amount")
	}
		
	return c.gateway.Pay(amount)
}

func NewCheckoutService(p PaymentGateway) *CheckoutService{
	return &CheckoutService{ gateway:p}
}

func NewPaymentGateway(method string) (PaymentGateway,error){
	switch method{
	case "stripe":
		return &StripeGateway{stripeID:123}, nil
	case "bitcoin":
		return &BitcoinGateway{bitcoinWallet:123},nil
	default:
		return nil, fmt.Errorf("Invalid payment method: %s",method)
	}
}

