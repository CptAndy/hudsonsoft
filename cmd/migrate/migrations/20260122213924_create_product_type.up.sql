
-- Statement # 1
CREATE OR REPLACE FUNCTION generate_tid (type_name text)
    RETURNS text
    AS $$
DECLARE
    -- 2 chars
    type_char char(2);
    digits int := 9;
    new_typeid text;
BEGIN
    -- 2 characters from the variation
    type_char := substr(UPPER(type_name), 1, 2);
    LOOP
        -- start at a certain number
        digits := digits + 1;
        -- store it
        new_typeid := type_char || digits;
        -- Check if it exists - if it does reloop the plan is to add 1 to the previous and try again
        IF NOT EXISTS (
            SELECT
                1
            FROM
                product_type
            WHERE
                type_id = new_typeid) THEN
        RETURN new_typeid;
    END IF;
END LOOP;
END;
$$
LANGUAGE plpgsql;

-- Statement # 2
CREATE TABLE IF NOT EXISTS product_type (
    id serial PRIMARY KEY,
    type_id varchar(4) UNIQUE NOT NULL,
    type_name text UNIQUE NOT NULL
);

-- Statement # 3
CREATE OR REPLACE FUNCTION set_type_id ()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.type_id := generate_tid (NEW.type_name);
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- Statement # 4
CREATE TRIGGER type_id_trigger
    BEFORE INSERT ON product_type
    FOR EACH ROW
    EXECUTE FUNCTION set_type_id ();
