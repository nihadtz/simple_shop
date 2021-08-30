CREATE TABLE User (
    id int(11) NOT NULL AUTO_INCREMENT,
    name varchar(128) NOT NULL DEFAULT '',
    email varchar(128) NOT NULL DEFAULT '',
    password varchar(128) NOT NULL DEFAULT '',
    type ENUM('customer','admin') NOT NULL DEFAULT 'customer',
    token varchar(255) NOT NULL DEFAULT '',
    alg varchar(10) NOT NULL DEFAULT '',
    issued int(11) NOT NULL DEFAULT 0,
    PRIMARY KEY(id),
    UNIQUE KEY(email)
    ) ENGINE=INNODB;

CREATE TABLE Product (
    id int(11) NOT NULL AUTO_INCREMENT,
    name varchar(128) NOT NULL DEFAULT '',
	description TEXT,
    category varchar(128) NOT NULL DEFAULT '',
    price DECIMAL(7,2) NOT NULL DEFAULT '0.00',
    PRIMARY KEY(id)
    ) ENGINE=INNODB;

CREATE TABLE Purchase (
    id int(11) NOT NULL AUTO_INCREMENT,
    date int(11) NOT NULL DEFAULT 0,
	product_id int(11) NOT NULL,
	user_id int(11) NOT NULL,
    quantity int(11) NOT NULL,
	status ENUM('new','paid') NOT NULL DEFAULT 'new',
    total DECIMAL(9,2) NOT NULL DEFAULT '0.00',
    PRIMARY KEY(id),
	FOREIGN KEY (product_id) REFERENCES Product(id),
	FOREIGN KEY (user_id) REFERENCES User(id)
    ) ENGINE=INNODB;

CREATE TABLE Payment (
    id int(11) NOT NULL AUTO_INCREMENT,
    date int(11) NOT NULL DEFAULT 0,
	purchase_id int(11) NOT NULL,
	gateway varchar(64) NOT NULL,
	gateway_payment_id varchar(64) NOT NULL,
	status ENUM('ok','failed') NOT NULL,
	error TEXT,
    PRIMARY KEY(id),
	FOREIGN KEY (purchase_id) REFERENCES Purchase(id)
    ) ENGINE=INNODB;