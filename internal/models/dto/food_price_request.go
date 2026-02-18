package dto

type CheckoutItem struct {
	Id  string `json:"id"`
	Qty int    `json:"qty"`
}

type FoodCheckoutRequest struct {
	BookId string `json:"book_id"`
	Items []CheckoutItem `json:"items"`
}
