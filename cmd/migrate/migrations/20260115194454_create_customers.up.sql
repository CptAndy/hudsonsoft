-- Statement # 1
CREATE EXTENSION IF NOT EXISTS citext;

-- Statement # 2
CREATE OR REPLACE FUNCTION generate_cuid (first_name text, last_name text)
    RETURNS text
    AS $$
DECLARE
    first_letter char;
    last_letter char;
    digits int;
    new_cuid text;
BEGIN
    -- First letters from name uppercase
    first_letter := UPPER(substr(first_name, 1, 1));
    last_letter := UPPER(substr(last_name, 1, 1));
    LOOP
        -- Generate 8 digits
        digits := floor(random() * (99999999 - 10000000) + 10000000)::int;
        -- stores
        new_cuid := first_letter || last_letter || digits;
        IF NOT EXISTS (
            SELECT
                1
            FROM
                customers
            WHERE
                cust_id = new_cuid) THEN
        RETURN new_cuid;
    END IF;
END LOOP;
END;
$$
LANGUAGE plpgsql;

-- Statement # 3
CREATE TABLE IF NOT EXISTS customers (
    id serial PRIMARY KEY,
    cust_id varchar(10) UNIQUE NOT NULL,
    first_name varchar NOT NULL,
    last_name varchar NOT NULL,
    email CITEXT,
    city varchar(100),
    state varchar(2),
    amount_spent numeric DEFAULT 0.00,
    product_owned int DEFAULT 0,
    product_returned int DEFAULT 0
);

-- Statement # 4
CREATE OR REPLACE FUNCTION set_cust_id()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.cust_id := generate_cuid (NEW.first_name, NEW.last_name);
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- Statement # 5
CREATE TRIGGER cust_id_trigger
    BEFORE INSERT ON customers
    FOR EACH ROW
    EXECUTE FUNCTION set_cust_id();

