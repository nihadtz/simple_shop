package models

import "github.com/nihadtz/simple_shop/services"

type DashboardProducts struct {
	ID             int     `db:"id" json:"id"`
	Name           string  `db:"name" json:"name"`
	Category       string  `db:"category" json:"category"`
	Price          float64 `db:"price" json:"price"`
	SoldQuantity   int     `db:"sold_quantity" json:"sold_quantity"`
	Income         float64 `db:"income" json:"income"`
	Paid           float64 `db:"paid" json:"paid"`
	NotPaid        float64 `db:"not_paid" json:"not_paid"`
	FailedPayments float64 `db:"failed_payments" json:"failed_payments"`
}

//Returns items with following data:
//SoldQuantity is total quantity sold in all Purchases
//Income is total value of all Purchases
//Paid is amount successfully paid by Customers
//NotPaid is difference between Income and Paid
//FailedPayments is amount failed during payment process
func ListAdminDashboard() ([]DashboardProducts, error) {
	db := services.Access.GetDB()
	var products []DashboardProducts

	query := `SELECT 
				p.id, p.name, p.category, p.price, SUM(IFNULL(ps.quantity,0)) sold_quantity, 
				SUM(IFNULL(ps.total,0)) income, SUM(IF(pp.status='ok', ps.total,0)) paid, 
				SUM(IF(ps.status='new', ps.total,0)) not_paid,
				SUM(IF(pp.status='failed', ps.total,0)) failed_payments
			FROM
				Product p
			LEFT JOIN 
				Purchase ps 
				ON ps.product_id = p.id
			LEFT JOIN 
				Payment pp
				ON pp.purchase_id = ps.id
			GROUP BY p.id
			ORDER BY sold_quantity desc`

	rows, err := db.Queryx(query)

	if err != nil {
		services.LogError("Error listing Products", err)
		return products, err
	}

	defer rows.Close()

	if rows == nil {
		return products, nil
	}

	for rows.Next() {
		var product DashboardProducts

		err = rows.StructScan(&product)

		if err != nil {
			services.LogError("Error scaning Product", err)
			return products, err
		}

		product.NotPaid = product.Income - product.Paid

		products = append(products, product)
	}

	return products, err
}
