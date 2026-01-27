DROP TRIGGER IF EXISTS set_sales_num_trigger ON products;
DROP FUNCTION IF EXISTS set_sales_num();
DROP FUNCTION IF EXISTS generate_snum(TEXT);
DROP TABLE IF EXISTS products;