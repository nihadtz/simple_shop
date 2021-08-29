package models

import (
	"github.com/nihadtz/simple_shop/services"
)

type Product struct {
	ID          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Category    string  `db:"category" json:"category"`
	Price       float64 `db:"price" json:"price"`
}

var createProduct = `CREATE TABLE Product (
    id int(11) NOT NULL AUTO_INCREMENT,
    name varchar(128) NOT NULL DEFAULT '',
	description TEXT,
    category varchar(128) NOT NULL DEFAULT '',
    price DECIMAL(7,2) NOT NULL DEFAULT '0.00',
    PRIMARY KEY(id)
    ) ENGINE=INNODB`

func (p *Product) Create() (err error) {
	db := services.Access.GetDB()

	query := `INSERT INTO Product
						(name, description, category, price) 
					VALUES
						(:name, :description, :category, :price)`

	result, err := db.NamedExec(query, p)

	if err == nil {
		productID, _ := result.LastInsertId()
		p.ID = int(productID)
	}

	return err
}

func (p *Product) Update() error {
	db := services.Access.GetDB()

	query := `UPDATE Product SET
				name=?, description=?, category=?, price=?
			WHERE 
				id=?`

	_, err := db.Exec(query, p.Name, p.Description, p.Category, p.Price, p.ID)

	if err != nil {
		services.LogError("Error updating Product", err)
	}

	return err
}

func (p *Product) Get(id int) error {
	db := services.Access.GetDB()

	query := `SELECT 
				id, name, description, category, price
			FROM
				Product
			WHERE
				id = ?`

	err := db.Get(p, query, id)

	return err
}

func ListProducts() ([]Product, error) {
	db := services.Access.GetDB()
	var products []Product

	query := `SELECT 
				id, name, description, category, price
			FROM
				Product`

	rows, err := db.Queryx(query)

	defer rows.Close()

	if err != nil {
		services.LogError("Error listing Products", err)
		return products, err
	}

	if rows == nil {
		return products, nil
	}

	for rows.Next() {
		var product Product

		err = rows.StructScan(&product)

		if err != nil {
			services.LogError("Error scaning Product", err)
			return products, err
		}

		products = append(products, product)
	}

	return products, err
}

func PopularProducts() ([]Product, error) {
	db := services.Access.GetDB()
	var products []Product

	query := `SELECT 
				p.id, p.name, p.description, p.category, p.price
			FROM
				Product p
			INNER JOIN 
				Purchase ps 
				ON ps.product_id = p.id
			GROUP by p.id
			ORDER BY SUM(IFNULL(ps.quantity,0)) desc
			LIMIT 5`

	rows, err := db.Queryx(query)

	if err != nil {
		services.LogError("Error listing popular Products", err)
		return products, err
	}

	if rows == nil {
		return products, nil
	}

	defer rows.Close()

	for rows.Next() {
		var product Product

		err = rows.StructScan(&product)

		if err != nil {
			services.LogError("Error scaning Product", err)
			return products, err
		}

		products = append(products, product)
	}

	return products, err
}
