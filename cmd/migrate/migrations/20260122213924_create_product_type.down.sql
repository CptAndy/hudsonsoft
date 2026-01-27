DROP TRIGGER IF EXISTS type_id_trigger ON product_type;
DROP FUNCTION IF EXISTS set_type_id();
DROP FUNCTION IF EXISTS generate_tid(TEXT);
DROP TABLE IF EXISTS product_type;