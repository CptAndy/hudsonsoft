CREATE TABLE IF NOT EXISTS stock (
  id SERIAL PRIMARY KEY,
  product_name TEXT UNIQUE,
  product_id TEXT UNIQUE,
  stock int DEFAULT 120,
  price NUMERIC(15,4) DEFAULT 10.00,
  inStock bool DEFAULT true,
  onOrder bool DEFAULT false
);