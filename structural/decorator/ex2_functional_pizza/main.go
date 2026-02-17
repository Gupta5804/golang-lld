package main

type Pizza func() int

func VeggieMania() Pizza {
	return func() int {
		return 15
	}
}

func AddCheese(p Pizza) Pizza {
	return func() int{
		return p() + 10
	}
}

func AddTomato(p Pizza) Pizza {
	return func() int{
		return p() + 5
	}
}
