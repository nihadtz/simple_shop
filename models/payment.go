package models

import "github.com/nihadtz/simple_shop/services"

type Payment struct {
	ID               int    `db:"id" json:"id"`
	Date             int64  `db:"date" json:"date"`
	PurchaseID       int    `db:"purchase_id" json:"purchase_id"`
	Gateway          string `db:"gateway" json:"gateway"`
	GatewayPaymentID string `db:"gateway_payment_id" json:"gateway_payment_id"`
	Status           string `db:"status" json:"status"`
	Error            string `db:"error" json:"error"`
}

var createPayment = `CREATE TABLE Payment (
    id int(11) NOT NULL AUTO_INCREMENT,
    date int(11) NOT NULL DEFAULT 0,
	purchase_id int(11) NOT NULL,
	gateway varchar(64) NOT NULL,
	gateway_payment_id varchar(64) NOT NULL,
	status ENUM('ok','failed') NOT NULL,
	error TEXT,
    PRIMARY KEY(id),
	FOREIGN KEY (purchase_id) REFERENCES Purchase(id)
    ) ENGINE=INNODB`

func (p *Payment) Create(purchase Purchase) (err error) {
	db := services.Access.GetDB()

	query := `INSERT INTO Payment
								(date, purchase_id, gateway, gateway_payment_id, status, error) 
							VALUES
								(:date, :purchase_id, :gateway, :gateway_payment_id, :status, :error)`

	result, err := db.NamedExec(query, p)

	if err == nil {
		paymentID, _ := result.LastInsertId()
		p.ID = int(paymentID)

		if p.Status == "ok" {
			purchase.Status = "paid"
			err = purchase.UpdateStatus()

			if err != nil {
				services.LogError("Error updating Purchase status", err)
				return err
			}
		}
	}
	return err
}

func (p *Purchase) ListPurchasePayments() error {
	db := services.Access.GetDB()

	query := `SELECT 
					id, date, purchase_id, gateway, gateway_payment_id, status, error
				FROM
					Payment
				WHERE
					purchase_id = :id`

	rows, err := db.NamedQuery(query, p)

	if err != nil {
		services.LogError("Error listing Payments for Purchase", err)
		return err
	}

	defer rows.Close()

	if rows == nil {
		return nil
	}

	for rows.Next() {
		var payment Payment

		err = rows.StructScan(&payment)

		if err != nil {
			services.LogError("Error scaning Payment for Purchase", err)
			return err
		}

		p.Payments = append(p.Payments, payment)
	}

	return err
}
