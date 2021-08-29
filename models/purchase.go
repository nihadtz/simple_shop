package models

import "github.com/nihadtz/simple_shop/services"

type Purchase struct {
	ID        int     `db:"id" json:"id"`
	Date      int64   `db:"date" json:"date"`
	ProductID int     `db:"product_id" json:"product_id"`
	UserID    int     `db:"user_id" json:"user_id"`
	Quantity  int     `db:"quantity" json:"quantity"`
	Status    string  `db:"status" json:"status"`
	Total     float64 `db:"total" json:"total"`
}

var createPurchase = `CREATE TABLE Purchase (
    id int(11) NOT NULL AUTO_INCREMENT,
    date int(11) NOT NULL DEFAULT 0,
	product_id int(11) NOT NULL,
	user_id int(11) NOT NULL,
    quantity int(11) NOT NULL,
	status ENUM('new','paid') NOT NULL DEFAULT 'new',
    total DECIMAL(9,2) NOT NULL DEFAULT '0.00',
    PRIMARY KEY(id),
	FOREIGN KEY (product_id) REFERENCES Product(id)
	FOREIGN KEY (user_id) REFERENCES User(id)
    ) ENGINE=INNODB`

type PurchaseFilter struct {
	UserID int `db:"user_id"`
}

func (p *Purchase) Create() (err error) {
	db := services.Access.GetDB()

	query := `INSERT INTO Purchase
							(date, product_id, user_id, quantity, total, status) 
						VALUES
							(:date, :product_id, :user_id, :quantity, :total, :status)`

	result, err := db.NamedExec(query, p)

	if err == nil {
		purchaseID, _ := result.LastInsertId()
		p.ID = int(purchaseID)
	}

	return err
}

func (p *Purchase) Get(id int) error {
	db := services.Access.GetDB()

	query := `SELECT 
					id, date, product_id, user_id, quantity, total, status
				FROM
				Purchase
				WHERE
					id = ?`

	err := db.Get(p, query, id)

	return err
}

func ListPurchases(filter PurchaseFilter) ([]Purchase, error) {
	db := services.Access.GetDB()
	var purchases []Purchase

	var where string

	if filter.UserID != 0 {
		where = " WHERE user_id = :user_id"
	}

	query := `SELECT 
					id, date, product_id, user_id, quantity, total, status
				FROM
					Purchase`

	if where != "" {
		query += where
	}

	rows, err := db.NamedQuery(query, filter)

	if err != nil {
		services.LogError("Error listing Purchases", err)
		return purchases, err
	}

	defer rows.Close()

	if rows == nil {
		return purchases, nil
	}

	for rows.Next() {
		var purchase Purchase

		err = rows.StructScan(&purchase)

		if err != nil {
			services.LogError("Error scaning Purchase", err)
			return purchases, err
		}

		purchases = append(purchases, purchase)
	}

	return purchases, err
}
