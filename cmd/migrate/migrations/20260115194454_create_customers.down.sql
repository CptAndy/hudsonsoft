
-- Drop trigger first
DROP TRIGGER IF EXISTS cust_id_trigger ON customers;

-- Drop trigger function
DROP FUNCTION IF EXISTS set_cust_id();

-- Drop UID generator function
DROP FUNCTION IF EXISTS generate_cuid(TEXT, TEXT);

-- Drop table that depends on citext
DROP TABLE IF EXISTS customers;

-- Now drop the extension
DROP EXTENSION IF EXISTS citext;