
CREATE EXTENSION IF NOT EXISTS citext;


CREATE OR REPLACE FUNCTION generate_uid(first_name TEXT, last_name TEXT)
RETURNS TEXT AS $$
DECLARE
    first_letter CHAR;
    last_letter CHAR;
    digits INT;
	new_uid TEXT;
BEGIN
    -- Take first letters and make uppercase
    first_letter := upper(substr(first_name, 1, 1));
    last_letter := upper(substr(last_name, 1, 1));
LOOP
    -- Generate  random digits
    digits := floor(random() * 90000000 + 10000000)::int;

	-- Combine and store
	new_uid := first_letter || last_letter || digits;

	-- Check if it exists
	IF NOT EXISTS (SELECT 1 FROM customers WHERE cust_id = new_uid) THEN
            RETURN new_uid; -- unique, return it
        END IF;

        -- Otherwise, loop again and try new digits
    END LOOP;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE customers(
    id SERIAL PRIMARY KEY,
    cust_id varchar(10) NOT NULL UNIQUE,
    first_name varchar NOT NULL,
    last_name varchar NOT NULL,
    Email    CITEXT,  
    city  varchar (2),
    state varchar(2),
    amount_spent     NUMERIC DEFAULT 0.00,
    product_owned    INT DEFAULT 0,
    product_returned INT DEFAULT 0
);

CREATE OR REPLACE FUNCTION set_cust_id()
RETURNS TRIGGER as $$
BEGIN
    NEW.cust_id := generate_uid(NEW.first_name, NEW.last_name);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cust_id_trigger
BEFORE INSERT ON customers
FOR EACH ROW
EXECUTE FUNCTION set_cust_id();