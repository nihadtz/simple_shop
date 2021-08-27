package models

type Product struct {
	ID          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Category    string  `db:"category" json:"category"`
	Price       float64 `db:"price" json:"price"`
}
