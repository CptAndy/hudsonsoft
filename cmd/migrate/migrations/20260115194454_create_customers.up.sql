
CREATE OR REPLACE FUNCTION generate_uid(first_name TEXT, last_name TEXT)
RETURNS TEXT AS $$
DECLARE
    first_letter CHAR;
    last_letter CHAR;
    digit INT;
	new_uid TEXT;
BEGIN
    -- Take first letters and make uppercase
    first_letter := upper(substr(first_name, 1, 1));
    last_letter := upper(substr(last_name, 1, 1));
LOOP
    -- Generate two random digits
    digit := floor(random() * 10)::text;

	-- Combine and store
	new_uid := first_letter || last_letter || digit;

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
    cust_id varchar(10),
    first_name varchar,
    last_name varchar,
    Email    CITEXT,  
    city  varchar ,
    state varchar,
    amount_spent     NUMERIC DEFAULT 0.00,
    product_owned    INT,
    product_returned INT,
);

CREATE OR REPLACE FUNCTION set_cust_id()
RETURNS TRIGGER as $$
BEGIN
    NEW.cust_id := generate_uid(NEW.first_name, NEW.last_name);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cust_id_trigger
BEFORE INSERT ON employees
FOR EACH ROW
EXECUTE FUNCTION set_cust_id();