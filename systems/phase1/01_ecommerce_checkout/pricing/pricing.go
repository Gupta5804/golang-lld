package pricing

import (
	"github.com/Gupta5804/golang-lld/systems/phase1/01_ecommerce_checkout/domain"
)

type PricingRule interface {
	Apply(cart *domain.Cart, currentTotal int64) int64
}

type PricingEngine struct {
	rules []PricingRule
}

func NewEngine(rules []PricingRule) *PricingEngine {
	return &PricingEngine{
		rules: rules,
	}
}

func (p *PricingEngine) Calculate(cart *domain.Cart) (int64, error) {
	current := cart.BaseTotal()
	for _, rule := range p.rules {
		current = rule.Apply(cart, current)
		if current < 0 {
			current = 0
		}
	}
	return current, nil
}

type PercentageDiscount struct {
	Percent int
}

func (p *PercentageDiscount) Apply(cart *domain.Cart, currentTotal int64) int64 {
	deduction := ((currentTotal * int64(p.Percent)) + 50) / 100
	return currentTotal - deduction
}

type FixedDiscount struct {
	FixedAmount int64
}

func (f *FixedDiscount) Apply(cart *domain.Cart, currentTotal int64) int64 {
	return currentTotal - f.FixedAmount
}

type TaxRule struct {
	Percent int
}

func (t *TaxRule) Apply(cart *domain.Cart, currentTotal int64) int64 {
	addition := (currentTotal * int64(t.Percent) + 50) / 100
	return currentTotal + addition
}
