package main


type DiscountStrategy func(amount float64) float64

func StandardPricing(amount float64) float64{
	return amount
}
func BlackFridayPricing(amount float64) float64{
	return 0.8 * amount
}


type Checkout struct{
	Total float64
}
func (c *Checkout) CalculateFinal(strategy DiscountStrategy) float64{
	return strategy(c.Total)
}
