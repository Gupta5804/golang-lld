package main

type Pizza interface{
	GetPrice() int
}

type VeggieMania struct{
	
}

func (v *VeggieMania) GetPrice() int{
	return 15
}

type CheeseTopping struct{
	pizza Pizza
}
func (c *CheeseTopping) GetPrice() int{
	return c.pizza.GetPrice() + 10
}

type TomatoTopping struct{
	pizza Pizza
}
func (t *TomatoTopping) GetPrice() int{
	return t.pizza.GetPrice() + 5
}