-- Statement # 1
CREATE OR REPLACE FUNCTION generate_snum (product_name text)
    RETURNS text
    AS $$
DECLARE
    -- Declare variables
    -- Take the first two characters from text
    two_char char(2);
    digits int;
    new_snum text;
BEGIN
    -- 2 characters from word
    two_char := upper(substr(product_name,1,2));
    LOOP
        -- Generate 4 random digits
        digits := floor(random() * (9999 - 1000) + 1000)::int;
        -- Store that value
        new_snum := two_char || digits;
        -- Check if exists
        IF NOT EXISTS (
            SELECT
                1
            FROM
                products
            WHERE
                sales_num = new_snum) THEN
        RETURN new_snum;
        -- return
    END IF;
    -- LOOP AGAIN IF MATCH WAS FOUND
END LOOP;
END;
$$
LANGUAGE plpgsql;

-- Statement # 2
CREATE TABLE IF NOT EXISTS products (
    id serial PRIMARY KEY,
    sales_num varchar(6) UNIQUE NOT NULL,
    product_name text UNIQUE NOT NULL
);

-- Statement # 3
-- TRIGGER CREATION for FUNCTION INTERACTION
CREATE OR REPLACE FUNCTION set_sales_num ()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.sales_num := generate_snum (NEW.product_name);
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- Statement # 4
CREATE TRIGGER set_sales_num_trigger
    BEFORE INSERT ON products
    FOR EACH ROW
    EXECUTE FUNCTION set_sales_num ();

