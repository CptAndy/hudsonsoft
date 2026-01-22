CREATE TABLE products (
id SERIAL PRIMARY KEY,
 product_name VARCHAR(100) UNIQUE NOT NULL,
  sales_number CHAR(8) UNIQUE NOT NULL,
  price NUMERIC NOT NULL,
  stock_quantity INT NOT NULL DEFAULT 0,
  type_id INT NULL,
  FOREIGN KEY (`type_id`) REFERENCES `Product_type` (`prod_type_id`)

);