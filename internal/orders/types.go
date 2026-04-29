package orders

type OrderItemInput struct {
	ProductID  int `json:"product_id"`
	Quantity   int `json:"quantity"`
	PriceCents int `json:"price_cents"`
}

type PlaceOrderRequest struct {
	CustomerID int              `json:"customer_id"`
	Items      []OrderItemInput `json:"items"`
}
