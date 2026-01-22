
CREATE EXTENSION IF NOT EXISTS citext;

CREATE OR REPLACE FUNCTION generate_cuid(first_name TEXT, last_name TEXT)
RETURNS TEXT AS $$
DECLARE
	first_letter CHAR;
	last_letter CHAR;
	digits INT;
	new_cuid TEXT;
BEGIN
    -- First letters from name uppercase
    first_letter := UPPER(substr(first_name,1,1));
    last_letter := UPPER(substr(last_name,1,1));
LOOP
    -- Generate 8 digits
    digits := floor(random() * (99999999 - 10000000) + 10000000)::int;

    -- stores
    new_cuid := first_letter || last_letter || digits;

    IF NOT EXISTS (SELECT 1 FROM customers WHERE cust_id = new_cuid) THEN
            RETURN new_cuid; -- return
        END IF;

    -- LOOP AGAIN IF THE UNIQUE ID IS FOUND
    END LOOP;
END;
$$ LANGUAGE plpgsql;



CREATE TABLE customers(
    id SERIAL PRIMARY KEY,
    cust_id varchar(10) UNIQUE NOT NULL,
    first_name varchar NOT NULL,
    last_name varchar NOT NULL,
    Email    CITEXT,  
    city  varchar (100),
    state varchar(2),
    amount_spent     NUMERIC DEFAULT 0.00,
    product_owned    INT DEFAULT 0,
    product_returned INT DEFAULT 0
);

CREATE OR REPLACE FUNCTION set_cust_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.cust_id := generate_cuid(NEW.first_name, NEW.last_name);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cust_id_trigger
BEFORE INSERT ON customers
FOR EACH ROW
EXECUTE FUNCTION set_cust_id();