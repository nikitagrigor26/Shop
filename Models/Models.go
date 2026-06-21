package Models

type Product struct {
	Id          int
	Name        string
	Price       float64
	Description string
	Quantity    int
}
type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Order struct {
	Id      int
	user_id int
	Status  string
}
type Order_item struct {
	Id         int
	product_id int
	order_id   int
	quantity   int
}
