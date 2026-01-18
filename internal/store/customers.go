package store

type Customer struct {
	ID               int64   `json:"id"`
	Customer_id      string  `json:"customer_id" validate:"required,min=10,max=10"`
	First_name       string  `json:"first_name" validate:"required,max=50"`
	Last_name        string  `json:"last_name" validate:"required,max=50"`
	Email            string  `json:"email" validate:"email"`
	City             string  `json:"city" validate:"required,min=2,max=2"`
	State            string  `json:"state" validate:"min=2,max=2"`
	Amount_spent     float64 `json:"amount_spent"`
	Product_owned    int64   `json:"product_owned"`
	Product_returned int64   `json:"product_returned"`
}


