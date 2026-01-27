-- Statement # 1
CREATE OR REPLACE FUNCTION generate_uid (first_name text, last_name text)
    RETURNS text
    AS $$
DECLARE
    first_letter char;
    last_letter char;
    digit1 int;
    digit2 int;
    new_uid text;
BEGIN
    -- Take first letters and make uppercase
    first_letter := upper(substr(first_name, 1, 1));
    last_letter := upper(substr(last_name, 1, 1));
    LOOP
        -- Generate two random digits
        digit1 := floor(random() * 10)::int;
        digit2 := floor(random() * 10)::int;
        -- Combine and store
        new_uid := first_letter || last_letter || digit1 || digit2;
        -- Check if it exists
        IF NOT EXISTS (
            SELECT
                1
            FROM
                employees
            WHERE
                emp_id = new_uid) THEN
        RETURN new_uid;
        -- unique, return it
    END IF;
    -- Otherwise, loop again and try new digits
END LOOP;
END;
$$
LANGUAGE plpgsql;

-- Statement # 2
CREATE TABLE IF NOT EXISTS employees (
    id serial PRIMARY KEY,
    emp_id varchar(4) UNIQUE,
    first_name varchar(50) NOT NULL,
    last_name varchar(50) NOT NULL,
    employee_pass bytea NOT NULL
);

-- Statement # 3
CREATE OR REPLACE FUNCTION set_emp_id ()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.emp_id := generate_uid (NEW.first_name, NEW.last_name);
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- Statement # 4
CREATE TRIGGER emp_id_trigger
    BEFORE INSERT ON employees
    FOR EACH ROW
    EXECUTE FUNCTION set_emp_id ();

