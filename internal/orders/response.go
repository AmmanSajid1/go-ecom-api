package orders

type OrderItemResponse struct {
	ProductID  int `json:"product_id"`
	Quantity   int `json:"quantity"`
	PriceCents int `json:"price_cents"`
}

type OrderResponse struct {
	ID         int                 `json:"id"`
	CustomerID int                 `json:"customer_id"`
	CreatedAt  string              `json:"created_at"`
	Items      []OrderItemResponse `json:"items"`
}
