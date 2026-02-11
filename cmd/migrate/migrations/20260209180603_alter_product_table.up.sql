ALTER TABLE products
ADD COLUMN type_id varchar(4) REFERENCES product_types(type_id);