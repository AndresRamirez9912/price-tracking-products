-- Delete if any table already exists
DROP TABLE IF EXISTS products_users;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS product_history;

CREATE TABLE users(
	id VARCHAR(50) PRIMARY KEY UNIQUE,
	email VARCHAR(50) NOT NULL,
	user_name VARCHAR(30) NOT NULL UNIQUE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products(
	id VARCHAR(20) PRIMARY KEY UNIQUE,
	name VARCHAR(150) NOT NULL,
	brand VARCHAR(20) NOT NULL,
	higher_price VARCHAR(20),
	lower_price VARCHAR(20),
	other_price VARCHAR(20),
	discount VARCHAR(20),
	image_url VARCHAR(150) NOT NULL,
	product_url VARCHAR(150) NOT NULL,
	store VARCHAR(20) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products_users (
	user_id VARCHAR(50) NOT NULL,
	product_id VARCHAR(20) NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE product_history (
	id SERIAL PRIMARY KEY UNIQUE,
	product_id VARCHAR(20) NOT NULL,
	product_name VARCHAR(150) NOT NULL,
	higher_price VARCHAR(20),
	lower_price VARCHAR(20),
	other_price VARCHAR(20),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (product_id) REFERENCES products(id)
);
