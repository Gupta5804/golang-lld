package domain

type Item struct {
	SKU      string
	Price    int64
	Quantity int
}

type Cart struct{
	Items []Item
}
func (c *Cart) BaseTotal()int64{
	var total int64
	for _,item := range c.Items{
		total += (int64(item.Quantity) * item.Price)
	}
	return total
}
